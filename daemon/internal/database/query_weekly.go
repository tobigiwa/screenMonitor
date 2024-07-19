package database

import (
	"cmp"
	"errors"
	"fmt"

	"slices"
	utils "utils"

	badger "github.com/dgraph-io/badger/v4"
)

func (bs *BadgerDBStore) GetWeek(day utils.Date) (WeeklyStat, error) {

	anyDayInTheWeek := utils.Date(day)
	date, _ := utils.ParseKey(anyDayInTheWeek)
	if IsFutureWeek(date) {
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

	date, _ := utils.ParseKey(anyDayInTheWeek)
	allConcernedDays := daysInThatWeek(date)

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

	if IsPastWeek(date) {
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
