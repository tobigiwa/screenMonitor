package database

import (
	"cmp"
	"errors"
	"fmt"
	helperFuncs "pkg/helper"
	"pkg/types"
	"slices"

	badger "github.com/dgraph-io/badger/v4"
)

func (bs *BadgerDBStore) GetWeek(day string) (WeeklyStat, error) {

	anyDayInTheWeek := types.Date(day)
	date, _ := ParseKey(anyDayInTheWeek)
	if IsFutureWeek(date) {
		return ZeroValueWeeklyStat, ErrFutureWeek
	}

	saturdayOfThatWeek := SaturdayOfTheWeek(date)

	byteData, err := bs.Get(dbWeekKey(types.Date(saturdayOfThatWeek)))
	errKeyNotFound := errors.Is(err, badger.ErrKeyNotFound)

	if err != nil && !errKeyNotFound {
		return ZeroValueWeeklyStat, err
	}

	if errKeyNotFound {
		return bs.getWeeklyAppStat(anyDayInTheWeek)
	}

	weekStat, err := helperFuncs.Decode[WeeklyStat](byteData)
	if err != nil {
		return ZeroValueWeeklyStat, err
	}
	return weekStat, nil
}

func (bs *BadgerDBStore) getWeeklyAppStat(anyDayInTheWeek types.Date) (WeeklyStat, error) {

	var (
		result     WeeklyStat
		weekTotal  types.Stats
		tmpStorage = make(map[string]types.Stats, 20)
	)

	date, _ := ParseKey(anyDayInTheWeek)
	allConcernedDays := daysInThatWeek(date)

	err := bs.db.View(func(txn *badger.Txn) error {

		for i := 0; i < len(allConcernedDays); i++ {
			day := allConcernedDays[i]

			dayStat, err := bs.GetDay(day)
			if err != nil {
				result.DayByDayTotal[i] = types.GenericKeyValue[types.Date, types.Stats]{Key: day, Value: types.Stats{}}
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
	eachAppSlice := make([]AppStat, 0, size)
	for app, stat := range tmpStorage {
		eachAppSlice = append(eachAppSlice, AppStat{app, stat})
	}

	slices.SortFunc(eachAppSlice, func(a, b AppStat) int {
		return cmp.Compare(b.Usage.Active, a.Usage.Active)
	})

	result.WeekTotal = weekTotal
	result.EachApp = eachAppSlice

	if IsPastWeek(date) {
		byteData, _ := helperFuncs.Encode(result)
		saturdayOfThatWeek := allConcernedDays[6]
		err := bs.setNewEntryToDB(dbWeekKey(saturdayOfThatWeek), byteData)
		if err != nil {
			fmt.Println("ERROR WRITING NEW WEEK ENTRY", saturdayOfThatWeek, "ERROR IS:", err)
		} else {
			fmt.Println("WRITING NEW WEEK ENTRY", saturdayOfThatWeek)
		}
	}

	return result, nil
}
