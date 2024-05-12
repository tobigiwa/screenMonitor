package repository

import (
	"bytes"
	"cmp"
	"encoding/gob"
	"errors"
	"fmt"
	"slices"
	"time"

	badger "github.com/dgraph-io/badger/v4"
)

type appStat struct {
	AppName string
	Usage   stats
}

type DailyStat struct {
	EachApp  []appStat
	DayTotal stats
}

type GenericKeyValue[K, V any] struct {
	Key   K
	Value V
}

type WeeklyStat struct {
	EachApp       []appStat
	WeekTotal     stats
	DayByDayTotal [7]GenericKeyValue[Date, stats]
}

func Encode[T any](tyPe T) ([]byte, error) {
	var r bytes.Buffer
	encoded := gob.NewEncoder(&r)
	if err := encoded.Encode(tyPe); err != nil {
		return nil, fmt.Errorf("%v:%w", err, ErrSerilization)
	}
	return r.Bytes(), nil
}

func Decode[T any](data []byte) (T, error) {
	var t, result T
	decoded := gob.NewDecoder(bytes.NewReader(data))
	if err := decoded.Decode(&result); err != nil {
		return t, fmt.Errorf("%v:%w", err, ErrDeserilization)
	}
	return result, nil
}

func (bs *BadgerDBStore) GetWeek(anyDayInTheWeek Date) (WeeklyStat, error) {
	date, _ := ParseKey(anyDayInTheWeek)
	if IsFutureWeek(date) {
		return ZeroValueWeeklyStat, ErrFutureWeek
	}

	saturdayOfThatWeek := SaturdayOfTheWeek(date)

	byteData, err := bs.Get(dbWeekKey(Date(saturdayOfThatWeek)))
	errKeyNotFound := errors.Is(err, badger.ErrKeyNotFound)

	if err != nil && !errKeyNotFound {
		return ZeroValueWeeklyStat, err
	}

	if errKeyNotFound {
		return bs.getWeeklyAppStat(anyDayInTheWeek)
	}

	weekStat, err := Decode[WeeklyStat](byteData)
	if err != nil {
		return ZeroValueWeeklyStat, err
	}
	return weekStat, nil
}

func (bs *BadgerDBStore) getWeeklyAppStat(anyDayInTheWeek Date) (WeeklyStat, error) {

	var (
		result     WeeklyStat
		weekTotal  stats
		tmpStorage = make(map[string]stats, 20)
	)

	date, _ := ParseKey(anyDayInTheWeek)
	allConcernedDays := daysInThatWeek(date)

	err := bs.db.View(func(txn *badger.Txn) error {

		for i := 0; i < len(allConcernedDays); i++ {
			day := allConcernedDays[i]

			dayStat, err := bs.GetDay(day)
			if err != nil {
				result.DayByDayTotal[i] = GenericKeyValue[Date, stats]{}
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
	eachAppSlice := make([]appStat, 0, size)
	for app, stat := range tmpStorage {
		eachAppSlice = append(eachAppSlice, appStat{app, stat})
	}

	slices.SortFunc(eachAppSlice, func(a, b appStat) int {
		return cmp.Compare(b.Usage.Active, a.Usage.Active)
	})

	result.WeekTotal = weekTotal
	result.EachApp = eachAppSlice

	if IsPastWeek(date) {
		byteData, _ := Encode(result)
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

func (bs *BadgerDBStore) GetDay(date Date) (DailyStat, error) {

	if day, _ := ParseKey(date); day.After(formattedToDay()) {
		return ZeroValueDailyStat, ErrFutureDay
	}

	byteData, err := bs.Get(dbDayKey(date))
	errKeyNotFound := errors.Is(err, badger.ErrKeyNotFound)

	if err != nil && !errKeyNotFound {
		return ZeroValueDailyStat, err
	}

	if errKeyNotFound {
		return bs.getDailyAppStat(date)
	}

	dayStat, err := Decode[DailyStat](byteData)
	if err != nil {
		return ZeroValueDailyStat, err
	}
	return dayStat, nil
}

func (bs *BadgerDBStore) getDailyAppStat(day Date) (DailyStat, error) {
	var (
		result       DailyStat
		dayTotalData stats
		arr          = make([]appStat, 0, 20)
	)

	err := bs.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = true
		opts.PrefetchSize = 100
		it := txn.NewIterator(opts)
		defer it.Close()

		prefix := dbAppPrefix
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			err := it.Item().Value(func(v []byte) error {

				var (
					app            appInfo
					appStatArrData appStat
				)

				if err := app.deserialize(v); err != nil {
					return err
				}

				thatDayStat := app.ScreenStat[day]

				// EachApp []appstat
				appStatArrData.AppName = app.AppName
				appStatArrData.Usage.Active = thatDayStat.Active
				appStatArrData.Usage.Inactive = thatDayStat.Inactive
				appStatArrData.Usage.Open = thatDayStat.Open
				arr = append(arr, appStatArrData)

				// DayTotall stats
				dayTotalData.Active += thatDayStat.Active
				dayTotalData.Inactive += thatDayStat.Inactive
				dayTotalData.Open += thatDayStat.Open

				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return ZeroValueDailyStat, err
	}

	slices.SortFunc(arr, func(a, b appStat) int {
		return cmp.Compare(b.Usage.Active, a.Usage.Active)
	})

	result.DayTotal.Active = dayTotalData.Active
	result.DayTotal.Inactive = dayTotalData.Inactive
	result.DayTotal.Open = dayTotalData.Open
	result.EachApp = arr

	if day != Date(formattedToDay().Format(timeFormat)) {
		byteData, _ := Encode(result)
		err := bs.setNewEntryToDB(dbDayKey(day), byteData)
		if err != nil {
			fmt.Println("ERROR WRITING NEW DAY ENTRY", day, "ERROR IS:", err)
		} else {
			fmt.Println("WRITING NEW DAY ENTRY", day)
		}
	}

	return result, nil
}

func (bs *BadgerDBStore) Get(key []byte) ([]byte, error) {

	var valCopy []byte

	err := bs.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}

		err = item.Value(func(val []byte) error {
			valCopy = append([]byte{}, val...)
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return valCopy, nil
}

func (bs *BadgerDBStore) setNewEntryToDB(key, byteData []byte) error {
	err := bs.db.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry(key, byteData)
		err := txn.SetEntry(e)
		return err
	})

	return err
}

func formattedToDay() time.Time {
	t, _ := ParseKey(Date(time.Now().Format(timeFormat)))
	return t
}

func IsFutureWeek(t time.Time) bool {
	now := time.Now()

	// Adjust the week start day to Sunday
	now = now.AddDate(0, 0, int(time.Sunday-now.Weekday()))
	t = t.AddDate(0, 0, int(time.Sunday-t.Weekday()))

	_, currentWeek := now.ISOWeek()
	_, inputWeek := t.ISOWeek()

	return inputWeek > currentWeek
}

func IsPastWeek(t time.Time) bool {
	now := time.Now()

	// Adjust the week start day to Sunday
	now = now.AddDate(0, 0, int(time.Sunday-now.Weekday()))
	t = t.AddDate(0, 0, int(time.Sunday-t.Weekday()))

	_, currentWeek := now.ISOWeek()
	_, inputWeek := t.ISOWeek()

	return inputWeek < currentWeek
}

func SaturdayOfTheWeek(t time.Time) string {
	daysUntilSaturday := 6 - int(t.Weekday())
	return t.AddDate(0, 0, daysUntilSaturday).Format(timeFormat)
}
