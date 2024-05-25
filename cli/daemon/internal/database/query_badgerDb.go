package database

import (
	"errors"
	"fmt"

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

func (bs *BadgerDBStore) WriteUsage(data types.ScreenTime) error {
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

		byteData, err := helperFuncs.Encode(app)
		if err != nil {
			return err
		}
		return txn.Set(dbAppKey(data.AppName), byteData)
	})
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
