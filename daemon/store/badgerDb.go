package store

import (
	"errors"
	"fmt"

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

func (bs *BadgerDBStore) WriteUsage(data ScreenTime) error {
	return bs.db.Update(func(txn *badger.Txn) error {

		item, err := txn.Get(dbAppKey(data.AppName))
		var (
			newApp  bool
			appInfo appInfo
			valCopy []byte
		)

		if newApp = errors.Is(err, badger.ErrKeyNotFound); err != nil && !newApp {
			return err
		}

		if newApp {
			appInfo.AppName = data.AppName
			appInfo.ScreenStat = make(dailyAppScreenTime)
			if icon, err := GetWmIcon(data.WindowID); err == nil {
				appInfo.Icon = icon
				appInfo.IsIconSet = true
			}
			if categories, err := getDesktopCategory(data.AppName); err == nil {
				appInfo.DesktopCategories = categories
				appInfo.IsCategorySet = true
			}
		}

		if err == nil && !newApp {
			if valCopy, err = item.ValueCopy(nil); err != nil {
				return err
			}

			if err = appInfo.deserialize(valCopy); err != nil {
				return err
			}

			if data.AppName != appInfo.AppName {
				return ErrAppKeyMismatch
			}

			if !appInfo.IsIconSet {
				if icon, err := GetWmIcon(data.WindowID); err == nil {
					appInfo.Icon = icon
					appInfo.IsIconSet = true
				}
			}

			if !appInfo.IsCategorySet {
				if categories, err := getDesktopCategory(data.AppName); err == nil {
					appInfo.DesktopCategories = categories
					appInfo.IsCategorySet = true
				}

			}
			fmt.Printf("existing appName:%v, time so far is: %v:%v\n\n", data.AppName, appInfo.ScreenStat[Key()].Active, appInfo.ScreenStat[Key()].Open)
		}

		if data.Type == Active {
			stat := appInfo.ScreenStat[Key()]
			stat.Active += data.Duration
			stat.ActiveTimeData = append(stat.ActiveTimeData, data.Interval)
			appInfo.ScreenStat[Key()] = stat
		} else {
			stat := appInfo.ScreenStat[Key()]
			stat.Inactive += data.Duration
			appInfo.ScreenStat[Key()] = stat
		}

		byteData, err := appInfo.serialize()
		if err != nil {
			return err
		}

		return txn.Set([]byte(data.AppName), byteData)
	})
}

func dbAppKey(appName string) []byte {
	return []byte(fmt.Sprintf("app:%v", appName))
}

func getAppPrefix() []byte {
	return []byte("app:")
}

func (bs *BadgerDBStore) ReadAll() error {

	keeper := make([]appInfo, 0, 30)
	err := bs.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = true
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()

		prefix := getAppPrefix()
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {

			err := it.Item().Value(func(v []byte) error {
				var app appInfo
				if err := app.deserialize(v); err != nil {
					return err
				}
				keeper = append(keeper, app)
				return nil
			})

			if err != nil {
				return err
			}
		}
		return nil
	})

	for i := 0; i < len(keeper); i++ {
		s := keeper[i]
		fmt.Println(s.AppName, s.ScreenStat[Key()].Active)
	}
	return err
}

func (bs *BadgerDBStore) BatchWriteUsage(data []ScreenTime) error {
	wb := bs.db.NewWriteBatch()
	defer wb.Cancel()

	for _, d := range data {
		item, err := bs.db.NewTransaction(false).Get([]byte(d.AppName))
		var newApp bool
		if newApp = errors.Is(err, badger.ErrKeyNotFound); err != nil && !newApp {
			return err
		}

		var appInfo appInfo

		if newApp {
			fmt.Printf("new app :%v\n\n", d.AppName)
			appInfo.AppName = d.AppName
			appInfo.ScreenStat = make(dailyAppScreenTime)
		}

		if err == nil {
			valCopy, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}
			if err := appInfo.deserialize(valCopy); err != nil {
				return err
			}
			if d.AppName != appInfo.AppName {
				return ErrAppKeyMismatch
			}
			fmt.Printf("existing appName:%v, time so far is: %v:%v\n\n", d.AppName, appInfo.ScreenStat[Key()].Active, appInfo.ScreenStat[Key()].Inactive)
		}

		if d.Type == Active {
			stat := appInfo.ScreenStat[Key()]
			stat.Active += d.Duration
			appInfo.ScreenStat[Key()] = stat
		} else {
			stat := appInfo.ScreenStat[Key()]
			stat.Inactive += d.Duration
			appInfo.ScreenStat[Key()] = stat
		}

		ser, err := appInfo.serialize()
		if err != nil {
			return err
		}

		if err := wb.Set([]byte(d.AppName), ser); err != nil {
			return err
		}
	}

	return wb.Flush()
}

// func (bs *BadgerDBStore) WriteUsage(data []ScreenTime) error {
//     // Group data by app name.
//     groupedData := make(map[string][]ScreenTime)
//     for _, d := range data {
//         groupedData[d.AppName] = append(groupedData[d.AppName], d)
//     }

//     wb := bs.db.NewWriteBatch()
//     defer wb.Cancel()

//     for appName, appData := range groupedData {
//         item, err := bs.db.NewTransaction(false).Get([]byte(appName))
//         var newApp bool
//         if newApp = errors.Is(err, badger.ErrKeyNotFound); err != nil && !newApp {
//             return err
//         }

//         var appInfo appInfo

//         if newApp {
//             fmt.Printf("new app :%v\n\n", appName)
//             appInfo.AppName = appName
//             appInfo.ActiveScreenStats = make(dailyAppScreenTime)
//             appInfo.PassiveScreenStats = make(dailyAppScreenTime)
//         }

//         if err == nil {
//             valCopy, err := item.ValueCopy(nil)
//             if err != nil {
//                 return err
//             }
//             if err := appInfo.deserialize(valCopy); err != nil {
//                 return err
//             }
//             if appName != appInfo.AppName {
//                 return ErrAppKeyMismatch
//             }
//             fmt.Printf("existing appName:%v, time so far is: %v:%v\n\n", appName, appInfo.ActiveScreenStats[Key()], appInfo.PassiveScreenStats[Key()])
//         }

//         // Sum screen time for this app.
//         for _, d := range appData {
//             if d.Type == Active {
//                 appInfo.ActiveScreenStats[Key()] += d.Time
//             } else {
//                 appInfo.PassiveScreenStats[Key()] += d.Time
//             }
//         }

//         ser, err := appInfo.serialize()
//         if err != nil {
//             return err
//         }

//         if err := wb.Set([]byte(appName), ser); err != nil {
//             return err
//         }
//     }

//     return wb.Flush()
// }
