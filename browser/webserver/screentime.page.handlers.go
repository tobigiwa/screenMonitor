package webserver

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"
	views "views"
)

type ScreenType string

const (
	Active     ScreenType = "active"
	Inactive   ScreenType = "inactive"
	Open       ScreenType = "open"
	timeFormat string     = "2006-01-02"
)

func (a *App) ScreenTimePageHandler(w http.ResponseWriter, r *http.Request) {
	if err := views.ScreenTimePage().Render(context.TODO(), w); err != nil {
		a.logger.Log(context.TODO(), slog.LevelError, err.Error())
		return
	}
}
func (a *App) WeekStat(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("week")
	endpoint := strings.TrimPrefix(r.URL.Path, "/")

	var msg Message

	switch query {
	case "thisweek":
		today := time.Now().Format(timeFormat)
		msg = Message{
			Endpoint:   endpoint,
			StringData: today,
		}
	case "lastweek":
	}

	bytes, err := msg.encode()
	if err != nil {
		a.logger.Log(context.TODO(), slog.LevelError, err.Error())
		return
	}
	if _, err = a.daemonConn.Write(bytes); err != nil {
		a.logger.Log(context.TODO(), slog.LevelError, err.Error())
		return
	}

	buf := make([]byte, 512)
	if _, err = a.daemonConn.Read(buf); err != nil {
		a.logger.Log(context.TODO(), slog.LevelError, err.Error())
		return
	}

	if err = msg.decode(buf); err != nil {
		a.logger.Log(context.TODO(), slog.LevelError, err.Error())
		return
	}

	data, err := msg.decodeToJson()
	if err != nil {
		a.logger.Log(context.TODO(), slog.LevelError, err.Error())
	}

	fmt.Printf("\n%+v\n\n", msg)
	w.Write(data)

}

func (a *App) CloseDaemonConnection() error {

	msg := Message{
		Endpoint: "closeConnection",
	}

	bytes, err := msg.encode()
	if err != nil {
		a.logger.Log(context.TODO(), slog.LevelError, err.Error())
		return err
	}
	if _, err = a.daemonConn.Write(bytes); err != nil {
		a.logger.Log(context.TODO(), slog.LevelError, err.Error())
		return err
	}
	return a.daemonConn.Close()
}
