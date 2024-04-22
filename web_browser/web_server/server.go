package backend

import (
	"browser/views"
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/dgraph-io/badger"
)

type App struct {
	logger *slog.Logger
}

func NewApp(logger *slog.Logger) *App {
	return &App{
		logger: logger,
	}
}

func (a App) ScreenTimePageHandler(w http.ResponseWriter, r *http.Request) {
	if err := views.ScreenTimePage().Render(context.TODO(), w); err != nil {
		a.logger.Log(context.TODO(), slog.LevelError, err.Error())
		return
	}
}
func (a App) sendScreenGraphData(w http.ResponseWriter, r *http.Request) {

}

func (a *App) Routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /screentime", a.ScreenTimePageHandler)

	fs := http.FileServer(http.Dir("./web_browser/assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets", fs))
	return mux
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
// 			err := item.Value(func(v []byte) error {
// 				var app appInfo
// 				if err := app.deserialize(v); err != nil {
// 					return err
// 				}
// 				// fmt.Printf("value=%+v\n\n", app)
// 				for key, value := range app.ScreenStat {
// 					// c[key] += value
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

type date string
type BadgerDBStore struct {
	db *badger.DB
}
type appInfo struct {
	AppName    string
	Icon       []byte
	Category   string
	ScreenStat dailyAppScreenTime
}

func (ap *appInfo) serialize() ([]byte, error) {
	var res bytes.Buffer
	encoded := gob.NewEncoder(&res)
	if err := encoded.Encode(ap); err != nil {
		return nil, fmt.Errorf("%v:%w", err, ErrSerilization)
	}
	return res.Bytes(), nil
}

func (ap *appInfo) deserialize(data []byte) error {
	decoded := gob.NewDecoder(bytes.NewReader(data))
	if err := decoded.Decode(ap); err != nil {
		return fmt.Errorf("%v:%w", err, ErrDeserilization)
	}
	return nil
}

type dailyAppScreenTime map[date]stats

type TimeInterval struct {
	Start time.Time
	End   time.Time
}

type stats struct {
	Active         float64
	Open           float64
	Inactive       float64
	ActiveTimeData []TimeInterval
}

var (
	ErrAppKeyMismatch = fmt.Errorf("key error: app name mismatch")
	ErrDeserilization = fmt.Errorf("error deserializing data")
	ErrSerilization   = fmt.Errorf("error serializing data")
)
