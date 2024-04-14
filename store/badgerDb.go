package store

import (
	"errors"
	"fmt"
	"log"
	"math"
	"strings"

	badger "github.com/dgraph-io/badger/v4"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
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
		var newApp bool
		if newApp = errors.Is(err, badger.ErrKeyNotFound); err != nil && !newApp {
			return err
		}

		var appInfo appInfo

		if newApp {
			fmt.Printf("new app :%v\n\n", data.AppName)
			appInfo.AppName = data.AppName
			appInfo.ActiveScreenStats = make(dailyAppScreenTime)
			appInfo.PassiveScreenStats = make(dailyAppScreenTime)
		}

		if err == nil {
			valCopy, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}
			if err := appInfo.deserialize(valCopy); err != nil {
				return err
			}
			if data.AppName != appInfo.AppName {
				return ErrAppKeyMismatch
			}
			fmt.Printf("existing appName:%v, time so far is: %v:%v\n\n", data.AppName, appInfo.ActiveScreenStats[Key()], appInfo.PassiveScreenStats[Key()])
		}

		if data.Type == Active {
			appInfo.ActiveScreenStats[Key()] += data.Time
		} else {
			appInfo.PassiveScreenStats[Key()] += data.Time
		}

		ser, err := appInfo.serialize()
		if err != nil {
			return err
		}

		return txn.Set([]byte(data.AppName), ser)
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
			appInfo.ActiveScreenStats = make(dailyAppScreenTime)
			appInfo.PassiveScreenStats = make(dailyAppScreenTime)
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
			fmt.Printf("existing appName:%v, time so far is: %v:%v\n\n", d.AppName, appInfo.ActiveScreenStats[Key()], appInfo.PassiveScreenStats[Key()])
		}

		if d.Type == Active {
			appInfo.ActiveScreenStats[Key()] += d.Time
		} else {
			appInfo.PassiveScreenStats[Key()] += d.Time
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

func (bs *BadgerDBStore) ReadAll() error {

	c := map[string]float64{}
	err := bs.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = true
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			// k := item.Key()
			err := item.Value(func(v []byte) error {
				// fmt.Printf("key=%s\n", string(k))
				var app appInfo
				if err := app.deserialize(v); err != nil {
					return err
				}
				// fmt.Printf("value=%+v\n\n", app)
				for key, value := range app.ActiveScreenStats {
					c[key] += value
				}

				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	fmt.Printf("c=%+v\n", c)
	barChat(c)
	return err
}

// func barChat(src map[string]float64) {
// 	if err := ui.Init(); err != nil {
// 		log.Fatalf("failed to initialize termui: %v", err)
// 	}
// 	defer ui.Close()

// 	var (
// 		labels = make([]string, 0, len(src))
// 		data   = make([]float64, 0, len(src))
// 	)

// 	for key, value := range src {
// 		form := ParseKey(key)
// 		labels = append(labels, fmt.Sprintf("%s %d %s", ShortenDay(form.Weekday()), form.Day(), form.Month()))
// 		data = append(data, value)

// 	}

// 	bc := widgets.NewBarChart()
// 	bc.Data = data
// 	bc.Labels = labels
// 	bc.Title = "Bar Chart"
// 	bc.SetRect(5, 5, 100, 25)
// 	bc.BarWidth = 5
// 	bc.BarColors = []ui.Color{ui.ColorRed, ui.ColorGreen}
// 	bc.LabelStyles = []ui.Style{ui.NewStyle(ui.ColorBlue)}
// 	bc.NumStyles = []ui.Style{ui.NewStyle(ui.ColorYellow)}

// 	ui.Render(bc)

// 	uiEvents := ui.PollEvents()
// 	for {
// 		e := <-uiEvents
// 		switch e.ID {
// 		case "q", "<C-c>":
// 			return
// 		}
// 	}
// }

func barChat(src map[string]float64) {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	var (
		labels = make([]string, 0, len(src))
		data   = make([]float64, 0, len(src))
	)

	for key, value := range src {
		form := ParseKey(key)
		// Split labels into multiple lines
		label := fmt.Sprintf("%s %d %s", ShortenDay(form.Weekday()), form.Day(), form.Month())
		label = strings.ReplaceAll(label, " ", "\n")
		labels = append(labels, label)
		data = append(data, value)
	}

	bc := widgets.NewBarChart()
	bc.Border = true
	bc.BorderStyle = ui.NewStyle(ui.ColorCyan)
	bc.Data = data
	bc.Labels = labels
	bc.Title = "Bar Chart"
	bc.TitleStyle.Bg = ui.ColorBlue
	bc.SetRect(5, 5, 100, 25)
	bc.BarWidth = 7
	bc.BarGap = 2
	bc.BarColors = []ui.Color{ui.ColorBlue, ui.ColorYellow}
	bc.LabelStyles = []ui.Style{ui.NewStyle(ui.ColorCyan)}
	bc.NumStyles = []ui.Style{ui.NewStyle(ui.ColorBlack)}
	bc.NumFormatter = func(n float64) string {
		return fmt.Sprintf("%.1f", roundTo1DP(n))
	}

	// Create a Paragraph widget for y-axis labels
	p := widgets.NewParagraph()
	p.Title = "Y-axis"
	p.Text = "0\n1\n2\n3\n4\n5\n6\n7\n8\n9\n10"
	p.SetRect(0, 20, 5, 25)

	ui.Render(bc, p)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		}
	}
}

func roundTo1DP(n float64) float64 {
	return math.Round(n*10) / 10
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
