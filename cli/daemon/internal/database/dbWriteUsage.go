package database

import (
	"fmt"
	helperFuncs "pkg/helper"
	"pkg/types"
	"strings"
	"time"

	"github.com/BurntSushi/xgb/xproto"
	badger "github.com/dgraph-io/badger/v4"
	"github.com/pkg/errors"
)

func (bs *BadgerDBStore) WriteUsage(data types.ScreenTime) error {
	return bs.db.Update(func(txn *badger.Txn) error {

		var (
			app     AppInfo
			valCopy []byte
		)

		item, err := txn.Get(dbAppKey(data.AppName))

		if newApp := errors.Is(err, badger.ErrKeyNotFound); err != nil {
			if !newApp { // there is an error AND is NOT ErrkeyNotFound
				return err
			}

			app.AppName = data.AppName
			app.ScreenStat = make(dailyAppScreenTime)

			addAppInfoForNewApp(data.WindowID, &app)
			fmt.Printf("New appName:%v, time so far is: %v:%v\n\n", app.AppName, app.IsCategorySet, app.IsIconSet)
			return updateAppStats(data, &app, txn)
		}

		if valCopy, err = item.ValueCopy(nil); err != nil {
			return err
		}

		if app, err = helperFuncs.DecodeJSON[AppInfo](valCopy); err != nil {
			return err
		}

		updateAppInfoForOldApp(data.WindowID, &app)
		app.AppName = data.AppName // !!!needs removing...
		fmt.Printf("Existing appName:%v, time so far is: %v:%v, brought in %f\n\n", data.AppName, app.ScreenStat[today()].Active, app.ScreenStat[today()].Open, data.Duration)
		return updateAppStats(data, &app, txn)

	})
}

func updateAppStats(data types.ScreenTime, app *AppInfo, txn *badger.Txn) error {

	todayStat, ok := app.ScreenStat[today()]

	if !ok { // we live to see a new day!!! 😎😎😎
		now := time.Now()
		midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		newDay := now.Sub(midnight)
		yesterdayDuration := data.Duration - newDay.Hours()

		if yesterdayDuration > 0 {
			updateYesterday(data.Type, app, yesterdayDuration)
			data.Duration = newDay.Hours()
		}
	}

	switch data.Type {
	case types.Active:
		todayStat.Active += data.Duration
		todayStat.ActiveTimeData = append(todayStat.ActiveTimeData, data.Interval)
	case types.Inactive:
		todayStat.Inactive += data.Duration
	case types.Open:
		todayStat.Open += data.Duration
	}

	app.ScreenStat[today()] = todayStat

	byteData, err := helperFuncs.EncodeJSON(app)
	if err != nil {
		return err
	}
	return txn.Set(dbAppKey(data.AppName), byteData)
}

func updateYesterday(screenType types.ScreenType, app *AppInfo, yesterdayDuration float64) {

	yesterdayStat := app.ScreenStat[yesterday()]
	switch screenType {
	case types.Active:
		yesterdayStat.Active += yesterdayDuration
	case types.Inactive:
		yesterdayStat.Inactive += yesterdayDuration
	case types.Open:
		yesterdayStat.Open += yesterdayDuration
	}
	app.ScreenStat[yesterday()] = yesterdayStat
}

func addAppInfoForNewApp(windowId xproto.Window, app *AppInfo) {
	if icon, err := GetWmIcon(windowId); err == nil {
		app.Icon = icon
		app.IsIconSet = true
	}

	if r, err := getDesktopCategoryAndCmd(app.AppName); err == nil {
		app.DesktopCategories = r.desktopCategories
		if len(r.desktopCategories) != 0 {
			app.DesktopCategories = r.desktopCategories
			for _, c := range r.desktopCategories {
				if category, ok := types.CategoryMap[strings.ToLower(c)]; ok {
					app.Category = category
					app.IsCategorySet = true
					break
				}
			}
		}
		app.CmdLine = r.cmdLine
		app.IsCmdLineSet = true
		fmt.Println("fetched info for new app", app.AppName, app.CmdLine, app.DesktopCategories)
	}
}

func updateAppInfoForOldApp(windowId xproto.Window, app *AppInfo) {
	if !app.IsIconSet {
		if icon, err := GetWmIcon(windowId); err == nil {
			app.Icon = icon
			app.IsIconSet = true
		}
	}
	fmt.Println("this app got here", app.AppName, !app.IsCmdLineSet && !app.IsCategorySet, app.IsCategorySet, app.IsCmdLineSet)
	if !app.IsCmdLineSet && !app.IsCategorySet {
		if r, err := getDesktopCategoryAndCmd(app.AppName); err == nil {
			if r.cmdLine != "" {
				app.CmdLine = r.cmdLine
				app.IsCmdLineSet = true
			}

			if len(r.desktopCategories) != 0 {
				app.DesktopCategories = r.desktopCategories
				for _, c := range r.desktopCategories {
					fmt.Printf("currently in category selection for app %s with c '%s'\n", app.AppName, strings.ToLower(c))
					if category, ok := types.CategoryMap[strings.ToLower(c)]; ok {
						app.Category = category
						app.IsCategorySet = true
						break
					}
				}
			}
			fmt.Println("fetched info for old app", app.AppName, app.CmdLine, app.DesktopCategories, app.Category)
		}
	}
}

func ExampleOf_opsFunc(v []byte) ([]byte, error) {
	var (
		app AppInfo
		err error
	)

	if app, err = helperFuncs.DecodeJSON[AppInfo](v); err != nil {
		return nil, err
	}
	// a := app.AppName
	// app.AppIconCategoryAndCmdLine = types.NoAppIconCategoryAndCmdLine
	// app.AppName = a
	fmt.Println(app.AppName, app.IsCategorySet, app.DesktopCategories, "category-", app.Category, app.IsCmdLineSet, app.CmdLine, app.IsIconSet)

	return helperFuncs.EncodeJSON(app)
}
