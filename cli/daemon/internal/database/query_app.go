package database

import (
	"errors"
	helperFuncs "pkg/helper"
	"pkg/types"
	"slices"

	badger "github.com/dgraph-io/badger/v4"
)

func (bs *BadgerDBStore) GetAppIconAndCategory(appNames []string) ([]types.AppIconCategoryAndCmdLine, error) {
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

			a := types.AppIconCategoryAndCmdLine{AppName: app.AppName}
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
			if app.IsCmdLineSet {
				a.CmdLine = app.CmdLine
				a.IsCmdLineSet = true
			}
			result[i] = a
		}
		return nil
	})

	return result, nil
}

func (bs *BadgerDBStore) AppWeeklyStat(appName string, anyDayInTheWeek types.Date) (types.AppRangeStat, error) {
	date, _ := ParseKey(anyDayInTheWeek)
	days := daysInThatWeek(date)
	return bs.appRangeStat(appName, days[:])
}

func (bs *BadgerDBStore) AppMonthlyStat(appName, month, year string) (types.AppRangeStat, error) {
	dates, err := AllTheDaysInMonth(year, month)
	if err != nil {
		return types.AppRangeStat{}, err
	}
	return bs.appRangeStat(appName, dates)
}

func (bs *BadgerDBStore) AppDateRangeStat(appName string, start, end types.Date) (types.AppRangeStat, error) {
	startDate, _ := ParseKey(start)
	endDate, _ := ParseKey(end)

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

	a, _ := bs.GetAppIconAndCategory([]string{appName})

	result.AppInfo = a[0]
	result.AppInfo.AppName = appName
	result.DaysRange = arr
	result.TotalRange = stat
	return result, nil
}
