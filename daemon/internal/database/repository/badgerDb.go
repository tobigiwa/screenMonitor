package repository

import (
	"errors"
	"fmt"
	"time"

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

			fmt.Printf("New appName:%v, time so far is: %v:%v\n\n", data.AppName, appInfo.IsCategorySet, appInfo.IsIconSet)
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
			fmt.Printf("Existing appName:%v, time so far is: %v:%v\n\n", data.AppName, appInfo.ScreenStat[Key()].Active, appInfo.ScreenStat[Key()].Open)
		}

		switch stat := appInfo.ScreenStat[Key()]; data.Type {
		case Active:
			stat.Active += data.Duration
			stat.ActiveTimeData = append(stat.ActiveTimeData, data.Interval)
			appInfo.ScreenStat[Key()] = stat
		case Inactive:
			stat.Inactive += data.Duration
			appInfo.ScreenStat[Key()] = stat
		case Open:
			stat.Open += data.Duration
			appInfo.ScreenStat[Key()] = stat
		}

		byteData, err := appInfo.serialize()
		if err != nil {
			return err
		}
		return txn.Set(dbAppKey(data.AppName), byteData)
	})
}

func (bs *BadgerDBStore) DeleteKey(key string) error {
	return bs.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
}

func (bs *BadgerDBStore) GetWeeklyScreenStats(s ScreenType, dayWhat string) ([]KeyValuePair, error) {

	var (
		newKey bool
		result = make([]KeyValuePair, 0, 7)
	)
	err := bs.db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("daily"))

		if newKey = errors.Is(err, badger.ErrKeyNotFound); err != nil && !newKey {
			return err
		}

		dailyST := make(dailyScreentimeAnalytics)
		if newKey {
			data, err := dailyST.serialize()
			if err != nil {
				return err
			}

			entry := badger.NewEntry([]byte("daily"), data)
			err = txn.SetEntry(entry)
			if err != nil {
				return err
			}

			return nil
		}

		var valCopy []byte
		if valCopy, err = item.ValueCopy(nil); err != nil {
			return err
		}
		if err = dailyST.deserialize(valCopy); err != nil {
			return err
		}

		wellFormattedDate, err := ParseKey(date(dayWhat))
		if err != nil {
			return err
		}

		statsForThatWeek := availableStatForThatWeek(wellFormattedDate)

		for i := 0; i < len(statsForThatWeek); i++ {
			day := statsForThatWeek[i]

			statsForThatDay, ok := dailyST[day]
			if !ok {
				statsForThatDay, err = bs.getDayActivity(day)
				if err != nil {
					continue
				}

				maybePastOrFutureDay, _ := ParseKey(day)
				today, _ := ParseKey(date(time.Now().Format(timeFormat)))
				if today.After(maybePastOrFutureDay) {
					dailyST[day] = statsForThatDay
					data, err := dailyST.serialize()
					if err == nil {
						txn.Set([]byte("daily"), data)
					}
				}
			}
			res := KeyValuePair{}
			switch s {
			case Active:
				res.Key = string(day)
				res.Value = statsForThatDay.Stats.Active
				result = append(result, res)
			case Inactive:
				res.Key = string(day)
				res.Value = statsForThatDay.Stats.Active
				result = append(result, res)
			case Open:
				res.Key = string(day)
				res.Value = statsForThatDay.Stats.Active
				result = append(result, res)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	if newKey {
		return nil, errors.New("new key created")
	}
	return result, nil
}

func (bs *BadgerDBStore) getDayActivity(day date) (dailyActiveScreentime, error) {

	var res dailyActiveScreentime
	err := bs.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = true
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()

		prefix := dbAppPrefix()
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			err := it.Item().Value(func(v []byte) error {

				var app appInfo
				if err := app.deserialize(v); err != nil {
					return err
				}

				thatDayStat := app.ScreenStat[day]

				res.Stats.Active += thatDayStat.Active
				res.Stats.Inactive += thatDayStat.Inactive
				res.Stats.Open += thatDayStat.Open

				return nil
			})

			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return dailyActiveScreentime{}, err
	}

	return res, nil
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
