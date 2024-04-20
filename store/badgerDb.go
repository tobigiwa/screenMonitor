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

func (bs *BadgerDBStore) get(key string) ([]byte, error) {
	var (
		valCopy []byte
		err     error
	)
	err = bs.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		valCopy, err = item.ValueCopy(nil)
		if err != nil {
			return err
		}
		return nil
	})
	return valCopy, err
}

func (bs *BadgerDBStore) set(key, value []byte) error {
	return bs.db.Update(
		func(txn *badger.Txn) error {
			return txn.Set([]byte(key), []byte(value))
		})
}

func (bs *BadgerDBStore) delete(key string) error {
	return bs.db.Update(
		func(txn *badger.Txn) error {
			return txn.Delete([]byte(key))
		})
}

func (bs *BadgerDBStore) WriteUsage(data ScreenTime) error {
	return bs.db.Update(func(txn *badger.Txn) error {

		item, err := txn.Get([]byte(data.AppName))

		var (
			newApp  bool
			appInfo appInfo
			valCopy []byte
		)

		if newApp = errors.Is(err, badger.ErrKeyNotFound); err != nil && !newApp {
			return err
		}

		if newApp {
			fmt.Printf("new app :%v\n\n", data.AppName)
			appInfo.AppName = data.AppName
			appInfo.ScreenStat = make(dailyAppScreenTime)
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

// func (bs *BadgerDBStore) ReadAll() error {

// 	c := map[date]float64{}
// 	err := bs.db.View(func(txn *badger.Txn) error {
// 		opts := badger.DefaultIteratorOptions
// 		opts.PrefetchValues = true
// 		opts.PrefetchSize = 10
// 		it := txn.NewIterator(opts)
// 		defer it.Close()
// 		for it.Rewind(); it.Valid(); it.Next() {
// 			item := it.Item()
// 			// k := item.Key()
// 			err := item.Value(func(v []byte) error {
// 				// fmt.Printf("key=%s\n", string(k))
// 				var app appInfo
// 				if err := app.deserialize(v); err != nil {
// 					return err
// 				}
// 				// fmt.Printf("value=%+v\n\n", app)
// 				for key, value := range app.ScreenStat {
// 					c[key] += value
// 				}

// 				return nil
// 			})
// 			if err != nil {
// 				return err
// 			}
// 		}
// 		return nil
// 	})
// 	fmt.Printf("c=%+v\n", c)
// 	return err
// }

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
