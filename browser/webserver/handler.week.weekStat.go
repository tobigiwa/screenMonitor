package webserver

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"pkg/types"
	"strings"
	"time"

	"github.com/a-h/templ"
)

const (
	timeFormat string = "2006-01-02"
)

type WeekStatDataCache struct {
	Day  string
	Data types.Message
}

var (
	lastRequest       = time.Now()
	weekStatCache     = make(map[string]templ.Component, 20)
	cacheLastSaturday string
)
 
func (a *App) WeekStatHandler(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query().Get("week")
	endpoint := strings.TrimPrefix(r.URL.Path, "/")

	var (
		msg types.Message
		err error
	)

	switch query {
	case "thisweek":
		if templComp, ok := weekStatCache["thisweek"]; ok && time.Since(lastRequest) <= 10*time.Minute {
			if err = templComp.Render(context.TODO(), w); err != nil {
				fmt.Println("error writing templ response", err)
				return
			}
			return
		}

		lastRequest = time.Now()
		today := time.Now().Format(timeFormat)
		msg = types.Message{
			Endpoint:          endpoint,
			StringDataRequest: today,
		}

	case "lastweek":
		lastSaturday := returnLastSaturday(time.Now())
		if templComp, ok := weekStatCache[lastSaturday]; ok {
			templComp.Render(context.TODO(), w)
			return
		}

		msg = types.Message{
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
			msg = types.Message{
				Endpoint:          endpoint,
				StringDataRequest: lastSaturday,
			}
		}

		if query == "forward-arrow" {
			if futureDate(t) {
				weekStatCache["thisweek"].Render(context.TODO(), w)
				return
			}
			lastSaturday = returnNextSaturday(t)
			msg = types.Message{
				Endpoint:          endpoint,
				StringDataRequest: lastSaturday,
			}
		}

		if templComp, ok := weekStatCache[lastSaturday]; ok {
			templComp.Render(context.TODO(), w)
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

		if templComp, ok := weekStatCache[lastSaturday]; ok {
			templComp.Render(context.TODO(), w)
			return
		}

		msg = types.Message{
			Endpoint:          endpoint,
			StringDataRequest: lastSaturday,
		}

		cacheLastSaturday = lastSaturday
	}

	// fmt.Println("would be consulting the deamonservice")

	msg, err = a.writeAndReadWithDaemonService(msg)
	if err != nil {
		fmt.Println("error occurred in writeAndReadWithDaemonService", err)
		return
	}

	templComp := prepareHtTMLResponse(msg)
	err = templComp.Render(context.TODO(), w)
	if err != nil {
		fmt.Println("err with templ:", err)
	}

	// Cache
	if query == "thisweek" {
		weekStatCache[query] = templComp
	} else if query == "backward-arrow" || query == "forward-arrow" || query == "lastweek" {
		weekStatCache[cacheLastSaturday] = templComp
	}

	templComp.Render(context.TODO(), w)
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
