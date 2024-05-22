package database

import (
	"cmp"
	"errors"
	"fmt"
	"slices"

	helperFuncs "pkg/helper"
	"pkg/types"

	badger "github.com/dgraph-io/badger/v4"
)

type BadgerDBStore struct {
	db *badger.DB
}

func NewBadgerDb(pathToDb string) (*BadgerDBStore, error) {
	opts := badger.DefaultOptions(pathToDb)

	opts.Logger = nil
	badgerInstance, err := badger.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("opening kv: %w", err)
	}

	return &BadgerDBStore{db: badgerInstance}, nil
}

func (bs *BadgerDBStore) Close() error {
	return bs.db.Close()
}

func (bs *BadgerDBStore) setNewEntryToDB(key, byteData []byte) error {
	err := bs.db.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry(key, byteData)
		err := txn.SetEntry(e)
		return err
	})

	return err
}

func (bs *BadgerDBStore) DeleteKey(key string) error {
	return bs.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
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

func (bs *BadgerDBStore) WriteUsage(data ScreenTime) error {
	return bs.db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get(dbAppKey(data.AppName))

		var (
			newApp  bool
			app     AppInfo
			valCopy []byte
		)

		if newApp = errors.Is(err, badger.ErrKeyNotFound); err != nil && !newApp {
			return err
		}

		if newApp {
			app.AppName = data.AppName
			app.ScreenStat = make(dailyAppScreenTime)
			if icon, err := GetWmIcon(data.WindowID); err == nil {
				app.Icon = icon
				app.IsIconSet = true
			}
			if categories, err := getDesktopCategory(data.AppName); err == nil {
				app.DesktopCategories = categories
				app.IsCategorySet = true
			}

			fmt.Printf("New appName:%v, time so far is: %v:%v\n\n", data.AppName, app.IsCategorySet, app.IsIconSet)
		}

		if err == nil && !newApp {
			if valCopy, err = item.ValueCopy(nil); err != nil {
				return err
			}

			if app, err = helperFuncs.Decode[AppInfo](valCopy); err != nil {
				return err
			}

			if data.AppName != app.AppName {
				return ErrAppKeyMismatch
			}

			if !app.IsIconSet {
				if icon, err := GetWmIcon(data.WindowID); err == nil {
					app.Icon = icon
					app.IsIconSet = true
				}
			}

			if !app.IsCategorySet {
				if categories, err := getDesktopCategory(data.AppName); err == nil {
					app.DesktopCategories = categories
					app.IsCategorySet = true
				}

			}
			fmt.Printf("Existing appName:%v, time so far is: %v:%v\n\n", data.AppName, app.ScreenStat[Key()].Active, app.ScreenStat[Key()].Open)
		}

		switch stat := app.ScreenStat[Key()]; data.Type {
		case Active:
			stat.Active += data.Duration
			stat.ActiveTimeData = append(stat.ActiveTimeData, data.Interval)
			app.ScreenStat[Key()] = stat
		case Inactive:
			stat.Inactive += data.Duration
			app.ScreenStat[Key()] = stat
		case Open:
			stat.Open += data.Duration
			app.ScreenStat[Key()] = stat
		}

		byteData, err := helperFuncs.Encode(app)
		if err != nil {
			return err
		}
		return txn.Set(dbAppKey(data.AppName), byteData)
	})
}

func (bs *BadgerDBStore) GetWeek(day string) (WeeklyStat, error) {

	anyDayInTheWeek := Date(day)
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

	weekStat, err := helperFuncs.Decode[WeeklyStat](byteData)
	if err != nil {
		return ZeroValueWeeklyStat, err
	}
	return weekStat, nil
}

func (bs *BadgerDBStore) getWeeklyAppStat(anyDayInTheWeek Date) (WeeklyStat, error) {

	var (
		result     WeeklyStat
		weekTotal  Stats
		tmpStorage = make(map[string]Stats, 20)
	)

	date, _ := ParseKey(anyDayInTheWeek)
	allConcernedDays := daysInThatWeek(date)

	err := bs.db.View(func(txn *badger.Txn) error {

		for i := 0; i < len(allConcernedDays); i++ {
			day := allConcernedDays[i]

			dayStat, err := bs.GetDay(day)
			if err != nil {
				result.DayByDayTotal[i] = GenericKeyValue[Date, Stats]{Key: day, Value: Stats{}}
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

	dayStat, err := helperFuncs.Decode[DailyStat](byteData)
	if err != nil {
		return ZeroValueDailyStat, err
	}
	return dayStat, nil
}

func (bs *BadgerDBStore) getDailyAppStat(day Date) (DailyStat, error) {
	var (
		result       DailyStat
		dayTotalData Stats
		arr          = make([]AppStat, 0, 20)
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
					app            AppInfo
					appStatArrData AppStat
					err            error
				)

				if app, err = helperFuncs.Decode[AppInfo](v); err != nil {
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

	slices.SortFunc(arr, func(a, b AppStat) int {
		return cmp.Compare(b.Usage.Active, a.Usage.Active)
	})

	result.DayTotal.Active = dayTotalData.Active
	result.DayTotal.Inactive = dayTotalData.Inactive
	result.DayTotal.Open = dayTotalData.Open
	result.EachApp = arr

	if day != Date(formattedToDay().Format(timeFormat)) {
		byteData, _ := helperFuncs.Encode(result)
		err := bs.setNewEntryToDB(dbDayKey(day), byteData)
		if err != nil {
			fmt.Println("ERROR WRITING NEW DAY ENTRY", day, "ERROR IS:", err)
		} else {
			fmt.Println("WRITING NEW DAY ENTRY", day)
		}
	}

	return result, nil
}

func (bs *BadgerDBStore) GetAppIconAndCategory(appNames []string) ([]types.AppIconAndCategory, error) {
	result := make([]types.AppIconAndCategory, len(appNames))
	bs.db.View(func(txn *badger.Txn) error {

		for i := 0; i < len(appNames); i++ {

			appName := appNames[i]
			item, err := txn.Get(dbAppKey(appName))
			if err != nil {
				result[i] = types.AppIconAndCategory{AppName: appName}
				continue
			}
			byteData, err := item.ValueCopy(nil)
			if err != nil {
				result[i] = types.AppIconAndCategory{AppName: appName}
				continue
			}
			app, err := helperFuncs.Decode[AppInfo](byteData)
			if err != nil {
				result[i] = types.AppIconAndCategory{AppName: appName}
				continue
			}

			a := types.AppIconAndCategory{AppName: app.AppName}
			if app.IsIconSet {
				a.Icon = app.Icon
				a.IsIconSet = true
			}
			if app.IsCategorySet {
				a.Category = string(app.Category)
				a.IsCategorySet = true
			} else {
				a.DesktopCategories = app.DesktopCategories
			}
			result[i] = a
		}
		return nil
	})

	return result, nil
}

// func (bs *BadgerDBStore) BatchWriteUsage(data []ScreenTime) error {
// 	wb := bs.db.NewWriteBatch()
// 	defer wb.Cancel()

// 	for _, d := range data {
// 		item, err := bs.db.NewTransaction(false).Get([]byte(d.AppName))
// 		var newApp bool
// 		if newApp = errors.Is(err, badger.ErrKeyNotFound); err != nil && !newApp {
// 			return err
// 		}

// 		var appInfo appInfo

// 		if newApp {
// 			fmt.Printf("new app :%v\n\n", d.AppName)
// 			appInfo.AppName = d.AppName
// 			appInfo.ScreenStat = make(dailyAppScreenTime)
// 		}

// 		if err == nil {
// 			valCopy, err := item.ValueCopy(nil)
// 			if err != nil {
// 				return err
// 			}
// 			if err := appInfo.deserialize(valCopy); err != nil {
// 				return err
// 			}
// 			if d.AppName != appInfo.AppName {
// 				return ErrAppKeyMismatch
// 			}
// 			fmt.Printf("existing appName:%v, time so far is: %v:%v\n\n", d.AppName, appInfo.ScreenStat[Key()].Active, appInfo.ScreenStat[Key()].Inactive)
// 		}

// 		if d.Type == Active {
// 			stat := appInfo.ScreenStat[Key()]
// 			stat.Active += d.Duration
// 			appInfo.ScreenStat[Key()] = stat
// 		} else {
// 			stat := appInfo.ScreenStat[Key()]
// 			stat.Inactive += d.Duration
// 			appInfo.ScreenStat[Key()] = stat
// 		}

// 		ser, err := appInfo.serialize()
// 		if err != nil {
// 			return err
// 		}

// 		if err := wb.Set([]byte(d.AppName), ser); err != nil {
// 			return err
// 		}
// 	}

// 	return wb.Flush()
// }
