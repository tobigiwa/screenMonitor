package database

import (
	"cmp"
	"errors"
	"fmt"
	"time"

	"slices"
	utils "utils"

	badger "github.com/dgraph-io/badger/v4"
)

func (bs *BadgerDBStore) GetWeek(anyDayInTheWeek utils.Date) (WeeklyStat, error) {

	date := utils.ToTimeType(anyDayInTheWeek)
	if utils.IsFutureWeek(date) {
		return ZeroValueWeeklyStat, ErrFutureWeek
	}

	saturdayOfThatWeek := utils.SaturdayOfTheWeek(date)

	byteData, err := bs.Get(dbWeekKey(saturdayOfThatWeek))
	errKeyNotFound := errors.Is(err, badger.ErrKeyNotFound)

	if err != nil && !errKeyNotFound {
		return ZeroValueWeeklyStat, err
	}

	if errKeyNotFound {
		return bs.getWeeklyAppStat(anyDayInTheWeek)
	}

	weekStat, err := utils.DecodeJSON[WeeklyStat](byteData)
	if err != nil {
		return ZeroValueWeeklyStat, err
	}
	return weekStat, nil
}

func (bs *BadgerDBStore) getWeeklyAppStat(anyDayInTheWeek utils.Date) (WeeklyStat, error) {

	var (
		result     WeeklyStat
		weekTotal  utils.Stats
		tmpStorage = make(map[string]utils.Stats, 20)
	)

	date := utils.ToTimeType(anyDayInTheWeek)
	allConcernedDays := utils.DaysInThatWeek(date)

	err := bs.db.View(func(txn *badger.Txn) error {

		for i := 0; i < len(allConcernedDays); i++ {
			day := allConcernedDays[i]

			dayStat, err := bs.GetDay(day)
			if err != nil {
				result.DayByDayTotal[i] = utils.GenericKeyValue[utils.Date, utils.Stats]{Key: day, Value: utils.Stats{}}
				continue
			}

			// DayByDayTotal [7]stats
			result.DayByDayTotal[i].Key = day
			result.DayByDayTotal[i].Value = dayStat.DayTotal

			// WeekTotal stats
			weekTotal.Active += dayStat.DayTotal.Active
			weekTotal.Inactive += dayStat.DayTotal.Inactive
			weekTotal.Open += dayStat.DayTotal.Open

			// EachApp []appStat
			for j := 0; j < len(dayStat.EachApp); j++ {
				eachAppName := dayStat.EachApp[j].AppName
				eachAppStat := dayStat.EachApp[j].Usage

				//get for that app
				thatAppStat := tmpStorage[eachAppName]
				//update it stat
				thatAppStat.Active += eachAppStat.Active
				thatAppStat.Inactive += eachAppStat.Inactive
				thatAppStat.Open += eachAppStat.Open
				//put it back
				tmpStorage[eachAppName] = thatAppStat

			}
		}
		return nil
	})
	if err != nil {
		return ZeroValueWeeklyStat, err
	}

	size := len(tmpStorage)
	eachAppSlice := make([]utils.AppStat, 0, size)
	for app, stat := range tmpStorage {
		eachAppSlice = append(eachAppSlice, utils.AppStat{AppName: app, Usage: stat})
	}

	slices.SortFunc(eachAppSlice, func(a, b utils.AppStat) int {
		return cmp.Compare(b.Usage.Active, a.Usage.Active)
	})

	result.WeekTotal = weekTotal
	result.EachApp = eachAppSlice

	if utils.IsPastWeek(date) {
		byteData, _ := utils.EncodeJSON(result)
		saturdayOfThatWeek := allConcernedDays[6]
		err := bs.setOrUpdateKeyValue(dbWeekKey(saturdayOfThatWeek), byteData)
		if err != nil {
			fmt.Println("ERROR WRITING NEW WEEK ENTRY", saturdayOfThatWeek, "ERROR IS:", err)
		} else {
			fmt.Println("WRITING NEW WEEK ENTRY", saturdayOfThatWeek)
		}
	}

	return result, nil
}

func (bs *BadgerDBStore) ReportWeeklyUsage(anyDayInTheWeek time.Time) (string, error) {

	PreviousWeekSaturday := utils.PreviousWeekSaturday(anyDayInTheWeek)

	var (
		theWeekStat, previousWeekStat WeeklyStat
		err                           error
	)

	if theWeekStat, err = bs.GetWeek(utils.ToDateType(anyDayInTheWeek)); err != nil {
		return "", err
	}
	if previousWeekStat, err = bs.GetWeek(utils.ToDateType(PreviousWeekSaturday)); err != nil {
		return "", err
	}

	theWeekDays, previousWeekDays := make([]float64, 7), make([]float64, 7)
	for i := 0; i < 7; i++ {
		theWeekDays = append(theWeekDays, theWeekStat.DayByDayTotal[i].Value.Active)
		previousWeekDays = append(previousWeekDays, previousWeekStat.DayByDayTotal[i].Value.Active)
	}

	theWeekDailyAverage := calculateAverage(theWeekDays)
	previousWeekDailyAverage := calculateAverage(previousWeekDays)

	if theWeekDailyAverage > previousWeekDailyAverage {
		return fmt.Sprintf("Daily Average: %s  ⬆️%.2f%% from previous week", utils.UsageTimeInHrsMin(theWeekDailyAverage), (theWeekDailyAverage-previousWeekDailyAverage)*100), nil
	}

	if theWeekDailyAverage < previousWeekDailyAverage {
		return fmt.Sprintf("Daily Average: %s  ⬇️%.2f%% from previous week", utils.UsageTimeInHrsMin(theWeekDailyAverage), (previousWeekDailyAverage-theWeekDailyAverage)*100), nil
	}

	return fmt.Sprintf("Daily Average: %s  same with previous week", utils.UsageTimeInHrsMin(theWeekDailyAverage)), nil
}

func calculateAverage(data []float64) float64 {
	var (
		nominator   float64
		denominator = len(data)
	)

	for i := 0; i < len(data); i++ {
		nominator += data[i]
	}
	return nominator / float64(denominator)
}
