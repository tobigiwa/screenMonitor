package database

import (
	"fmt"
	helperFuncs "pkg/helper"
	"pkg/types"

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
		fmt.Printf("Existing appName:%v, time so far is: %v:%v:%v:%v\n\n", data.AppName, app.ScreenStat[Key()].Active, app.ScreenStat[Key()].Open, app.IsCmdLineSet, app.IsCategorySet)
		return updateAppStats(data, &app, txn)

	})
}

func updateAppStats(data types.ScreenTime, app *AppInfo, txn *badger.Txn) error {

	switch stat := app.ScreenStat[Key()]; data.Type {
	case types.Active:
		stat.Active += data.Duration
		stat.ActiveTimeData = append(stat.ActiveTimeData, data.Interval)
		app.ScreenStat[Key()] = stat
	case types.Inactive:
		stat.Inactive += data.Duration
		app.ScreenStat[Key()] = stat
	case types.Open:
		stat.Open += data.Duration
		app.ScreenStat[Key()] = stat
	}

	byteData, err := helperFuncs.EncodeJSON(app)
	if err != nil {
		return err
	}
	return txn.Set(dbAppKey(data.AppName), byteData)
}

func addAppInfoForNewApp(windowId xproto.Window, app *AppInfo) {
	if icon, err := GetWmIcon(windowId); err == nil {
		app.Icon = icon
		app.IsIconSet = true
	}

	if r, err := getDesktopCategoryAndCmd(app.AppName); err == nil {
		app.DesktopCategories = r.desktopCategories
		app.IsCategorySet = true
		app.CmdLine = r.cmdLine
		fmt.Println("fetched cmdLine for new app", app.CmdLine)
		app.IsCmdLineSet = true
	}
}

func updateAppInfoForOldApp(windowId xproto.Window, app *AppInfo) {
	if !app.IsIconSet {
		if icon, err := GetWmIcon(windowId); err == nil {
			app.Icon = icon
			app.IsIconSet = true
		}
	}

	if !app.IsCmdLineSet {
		if r, err := getDesktopCategoryAndCmd(app.AppName); err == nil {
			if len(r.desktopCategories) != 0 {
				app.DesktopCategories = r.desktopCategories
				app.IsCategorySet = true
			}
			if r.cmdLine != "" {
				app.CmdLine = r.cmdLine
				app.IsCmdLineSet = true
				fmt.Println("fetched cmdLine for old app", app.AppName, app.CmdLine, app.DesktopCategories)
			}
		}
	}
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
