package webserver

import (
	"context"
	"fmt"
	"log"
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

type WeekStatDataCache struct {
	Day  string
	Data Message
}

var (
	lastRequest       = time.Now()
	weekStatCache     = make(map[string][]byte, 20)
	cacheLastSaturday string
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
		if jsonResponse, ok := weekStatCache["thisweek"]; ok && time.Since(lastRequest) <= 10*time.Minute {
			w.Write(jsonResponse)
			return
		}

		lastRequest = time.Now()
		today := time.Now().Format(timeFormat)
		msg = Message{
			Endpoint:          endpoint,
			StringDataRequest: today,
		}

	case "lastweek":
		lastSaturday := returnLastSaturday(time.Now())
		if jsonResponse, ok := weekStatCache[lastSaturday]; ok {
			w.Write(jsonResponse)
			return
		}

		msg = Message{
			Endpoint:          endpoint,
			StringDataRequest: lastSaturday,
		}

		cacheLastSaturday = lastSaturday

	case "backward-arrow", "forward-arrow":
		var (
			t            time.Time
			err          error
			lastSaturday string
			q            string
		)

		if q = r.URL.Query().Get("saturday"); q == "" {
			fmt.Println("empty q")
			log.Fatal(err)
		}

		if t, err = time.Parse(timeFormat, q); err != nil {
			log.Fatal(err)
		}

		if query == "backward-arrow" {
			lastSaturday = returnLastSaturday(t)
			msg = Message{
				Endpoint:          endpoint,
				StringDataRequest: lastSaturday,
			}
		}

		if query == "forward-arrow" {
			if futureDate(t) {
				w.Write(weekStatCache["thisweek"])
				return
			}
			lastSaturday = returnNextSaturday(t)
			msg = Message{
				Endpoint:          endpoint,
				StringDataRequest: lastSaturday,
			}
		}

		if jsonResponse, ok := weekStatCache[lastSaturday]; ok {
			w.Write(jsonResponse)
			return
		}

		cacheLastSaturday = lastSaturday

	case "month":
		var lastSaturday, q string
		if q = r.URL.Query().Get("month"); q == "" {
			log.Fatalf("empty query")
		}
		if lastSaturday = lastSaturdayOfTheMonth(q); lastSaturday == "" {
			log.Fatal("invalid input")
		}

		if jsonResponse, ok := weekStatCache[lastSaturday]; ok {
			w.Write(jsonResponse)
			return
		}

		msg = Message{
			Endpoint:          endpoint,
			StringDataRequest: lastSaturday,
		}

		cacheLastSaturday = lastSaturday
	}

	fmt.Println("would be consulting the deamonservice")
	jsonResponse, err := a.writeToFrontend(msg)
	if err != nil {
		fmt.Println("error occurred in writeToFrontend", err)
		return
	}

	// Cache
	if query == "thisweek" {
		weekStatCache[query] = jsonResponse
	} else if query == "backward-arrow" || query == "forward-arrow" || query == "lastweek" {
		weekStatCache[cacheLastSaturday] = jsonResponse
	}

	w.Write(jsonResponse)
}

func (a *App) writeToFrontend(msg Message) ([]byte, error) {
	bytes, err := msg.encode() // encode message in byte
	if err != nil {
		return nil, err
	}
	if _, err = a.daemonConn.Write(bytes); err != nil { // write to socket
		return nil, err
	}

	buf := make([]byte, 1024)
	if _, err = a.daemonConn.Read(buf); err != nil { // wait and read response from socket
		return nil, err
	}

	if err = msg.decode(buf); err != nil { // decode response to Message struct
		return nil, err
	}

	return msg.decodeToJson() // convert response to json
}

func returnLastSaturday(t time.Time) string {

	if t.Weekday() == time.Saturday {
		return t.AddDate(0, 0, -7).Format(timeFormat)
	}

	daysSinceSaturday := int(t.Weekday()+1) % 7
	return t.AddDate(0, 0, -daysSinceSaturday).Format(timeFormat)
}

func returnNextSaturday(t time.Time) string {
	return t.AddDate(0, 0, 7).Format(timeFormat)
}

func futureDate(t time.Time) bool {
	today := time.Now()
	nextWeekDay := t.AddDate(0, 0, 7)
	return nextWeekDay.After(today)
}

func lastSaturdayOfTheMonth(month string) string {
	t, err := time.Parse("January", month)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	NextMonth := time.Date(time.Now().Year(), t.Month()+1, 1, 0, 0, 0, 0, time.UTC)

	var s time.Time
	for {
		NextMonth = NextMonth.AddDate(0, 0, -1)
		if NextMonth.Weekday() == time.Saturday {
			s = NextMonth
			break
		}
	}
	return s.Format(timeFormat)
}
