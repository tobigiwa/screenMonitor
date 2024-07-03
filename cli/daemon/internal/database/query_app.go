package database

import (
	"errors"
	helperFuncs "pkg/helper"
	"pkg/types"
	"slices"

	badger "github.com/dgraph-io/badger/v4"
)

func (bs *BadgerDBStore) GetAppIconCategoryAndCmdLine(appNames []string) ([]types.AppIconCategoryAndCmdLine, error) {
	result := make([]types.AppIconCategoryAndCmdLine, len(appNames))
	bs.db.View(func(txn *badger.Txn) error {

		for i := 0; i < len(appNames); i++ {

			appName := appNames[i]
			item, err := txn.Get(dbAppKey(appName))
			if err != nil {
				result[i] = types.AppIconCategoryAndCmdLine{AppName: appName}
				continue
			}
			byteData, err := item.ValueCopy(nil)
			if err != nil {
				result[i] = types.AppIconCategoryAndCmdLine{AppName: appName}
				continue
			}
			app, err := helperFuncs.DecodeJSON[AppInfo](byteData)
			if err != nil {
				result[i] = types.AppIconCategoryAndCmdLine{AppName: appName}
				continue
			}

			a := app.AppIconCategoryAndCmdLine
			result[i] = a
		}
		return nil
	})

	return result, nil
}

func (bs *BadgerDBStore) AppWeeklyStat(appName string, anyDayInTheWeek types.Date) (types.AppRangeStat, error) {
	date, _ := helperFuncs.ParseKey(anyDayInTheWeek)
	days := daysInThatWeek(date)
	return bs.appRangeStat(appName, days[:])
}

func (bs *BadgerDBStore) AppMonthlyStat(appName, month, year string) (types.AppRangeStat, error) {
	dates, err := helperFuncs.AllTheDaysInMonth(year, month)
	if err != nil {
		return types.AppRangeStat{}, err
	}
	return bs.appRangeStat(appName, dates)
}

func (bs *BadgerDBStore) AppDateRangeStat(appName string, start, end types.Date) (types.AppRangeStat, error) {
	startDate, _ := helperFuncs.ParseKey(start)
	endDate, _ := helperFuncs.ParseKey(end)

	if !endDate.After(startDate) {
		return types.AppRangeStat{}, errors.New("end date is not after start date")
	}

	dates := make([]types.Date, 0, 31)
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		dates = append(dates, types.Date(d.Format(types.TimeFormat)))
	}

	return bs.appRangeStat(appName, slices.Clip(dates))
}

func (bs *BadgerDBStore) appRangeStat(appName string, dateRange []types.Date) (types.AppRangeStat, error) {

	var (
		result types.AppRangeStat
		app    AppInfo
		err    error
	)

	if err = bs.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(dbAppKey(appName))
		if err != nil {
			return err
		}
		byteData, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		app, err = helperFuncs.DecodeJSON[AppInfo](byteData)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return types.AppRangeStat{}, err
	}

	var stat types.Stats
	arr := make([]types.GenericKeyValue[types.Date, types.Stats], len(dateRange))
	for i := 0; i < len(dateRange); i++ {
		dayStat := app.ScreenStat[dateRange[i]]
		arr[i] = types.GenericKeyValue[types.Date, types.Stats]{Key: dateRange[i], Value: dayStat}
		stat.Active += dayStat.Active
		stat.Inactive += dayStat.Inactive
		stat.Open += dayStat.Open
	}

	a, _ := bs.GetAppIconCategoryAndCmdLine([]string{appName})

	result.AppInfo = a[0]
	result.AppInfo.AppName = appName
	result.DaysRange = arr
	result.TotalRange = stat
	return result, nil
}

func (bs *BadgerDBStore) SetAppCategory(appName string, category types.Category) error {
	byteData, err := bs.Get(dbAppKey(appName))
	if err != nil {
		return err
	}
	appInfo, err := helperFuncs.DecodeJSON[AppInfo](byteData)
	if err != nil {
		return err
	}
	appInfo.Category = category
	appInfo.IsCategorySet = true

	byteData, err = helperFuncs.EncodeJSON(appInfo)
	if err != nil {
		return err
	}

	return bs.setOrUpdateKeyValue(dbAppKey(appName), byteData)
}
func (bs *BadgerDBStore) GetAllACategories() ([]types.Category, error) {
	byteData, err := bs.Get(dbCategoryKey)

	if err != nil {
		if !errors.Is(err, badger.ErrKeyNotFound) {
			return nil, err
		}

		byteData, err := helperFuncs.EncodeJSON(types.DefalutCategory)
		if err != nil {
			return nil, err
		}

		if err = bs.setOrUpdateKeyValue(dbCategoryKey, byteData); err != nil {
			return nil, err
		}

		return types.DefalutCategory, nil
	}

	return helperFuncs.DecodeJSON[[]types.Category](byteData)
}

func (bs *BadgerDBStore) GetAllApp() ([]types.AppIconCategoryAndCmdLine, error) {
	res := make([]types.AppIconCategoryAndCmdLine, 0, 30)

	err := bs.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = true
		opts.PrefetchSize = 100
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Seek(dbAppPrefix); it.ValidForPrefix(dbAppPrefix); it.Next() {
			err := it.Item().Value(func(val []byte) error {

				var (
					app AppInfo
					err error
				)

				if app, err = helperFuncs.DecodeJSON[AppInfo](val); err != nil {
					return err
				}

				if len(res) == cap(res) {
					res = slices.Grow(res, 10)
				}
				// if app.IsCmdLineSet {
				res = append(res, app.AppIconCategoryAndCmdLine)
				// }

				return nil
			})

			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return slices.Clip(res), nil
}

func (bs *BadgerDBStore) GetAppByName(appName string) (AppInfo, error) {
	byteData, err := bs.Get(dbAppKey(appName))
	if err != nil {
		return AppInfo{}, err
	}
	return helperFuncs.DecodeJSON[AppInfo](byteData)
}

func (bs *BadgerDBStore) GetAppTodayActiveStatSoFar(appName string) (float64, error) {
	appInfo, err := bs.GetAppByName(appName)
	if err != nil {
		return 0, err
	}
	return appInfo.ScreenStat[helperFuncs.Today()].Active, nil
}

func (bs *BadgerDBStore) GetAllAppWithCmdLine() ([]types.AppIconCategoryAndCmdLine, error) {
	r, err := bs.GetAllApp()
	if err != nil {
		return nil, err
	}

	r = slices.DeleteFunc(r, func(i types.AppIconCategoryAndCmdLine) bool {
		return !i.IsCmdLineSet
	})
	return r, nil
}
