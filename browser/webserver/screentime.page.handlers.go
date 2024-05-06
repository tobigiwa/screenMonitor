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
	fmt.Println(endpoint, query)

	switch query {
	case "thisweek":
		today := time.Now().Format(timeFormat)
		msg := Message{
			Endpoint: endpoint,
			Body:     today,
		}

		bytes, err := msg.encode()
		if err != nil {
			a.logger.Log(context.TODO(), slog.LevelError, err.Error())
			return
		}
		i, err := a.daemonConn.Write(bytes)
		if err != nil {
			a.logger.Log(context.TODO(), slog.LevelError, err.Error())
			return
		}
		fmt.Println("write successful", i)
		return

	case "lastweek":
	}
}

func (a *App) CloseDaemonConnection() error {
	return a.daemonConn.Close()
}
