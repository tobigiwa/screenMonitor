package webserver

import (
	"context"
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
			Endpoint:          endpoint,
			StringDataRequest: today,
		}
	case "lastweek":
		today := time.Now()
		var lastSaturday time.Time

		if today.Weekday() == time.Saturday {
			lastSaturday = today.AddDate(0, 0, -7)
		} else {
			daysSinceSaturday := int(today.Weekday()+1) % 7
			lastSaturday = today.AddDate(0, 0, -daysSinceSaturday)
		}

		msg = Message{
			Endpoint:          endpoint,
			StringDataRequest: lastSaturday.Format(timeFormat),
		}
	}

	bytes, err := msg.encode() // encode message in byte
	if err != nil {
		a.logger.Log(context.TODO(), slog.LevelError, err.Error())
		return
	}
	if _, err = a.daemonConn.Write(bytes); err != nil { // write to socket
		a.logger.Log(context.TODO(), slog.LevelError, err.Error())
		return
	}

	buf := make([]byte, 512)
	if _, err = a.daemonConn.Read(buf); err != nil { // wait and read response from socket
		a.logger.Log(context.TODO(), slog.LevelError, err.Error())
		return
	}

	if err = msg.decode(buf); err != nil { // decode response to Message struct
		a.logger.Log(context.TODO(), slog.LevelError, err.Error())
		return
	}

	data, err := msg.decodeToJson() // convert response to json
	if err != nil {
		a.logger.Log(context.TODO(), slog.LevelError, err.Error())
	}

	w.Write(data) // write json response to http response

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

func (a *App) writeToFrontend(buf []byte) error {
	return nil
}
