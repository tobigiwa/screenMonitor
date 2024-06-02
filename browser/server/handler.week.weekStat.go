package webserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"pkg/types"
	"strings"
	"time"

	helperFuncs "pkg/helper"

	"github.com/a-h/templ"
)

var (
	lastRequest    = time.Now()
	weekStatCache  = make(map[string]templ.Component, 20)
	cachedSaturday string
)

const HeaderKey = "Saturday"

func (a *App) WeekStatHandler(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query().Get("week")
	endpoint := strings.TrimPrefix(r.URL.Path, "/")

	var (
		msg types.Message
		err error
	)

	switch query {
	case "thisweek":
		saturdayOfTheWeek := helperFuncs.SaturdayOfTheWeek(time.Now())
		if templComp, ok := weekStatCache["thisweek"]; ok && time.Since(lastRequest) <= 10*time.Minute {
			w.Header().Set(HeaderKey, saturdayOfTheWeek)
			if err = templComp.Render(context.TODO(), w); err != nil {
				a.serverError(w, err)
			}
			return
		}
		lastRequest = time.Now()
		msg = types.Message{
			Endpoint:        endpoint,
			WeekStatRequest: saturdayOfTheWeek,
		}
		cachedSaturday = saturdayOfTheWeek

	case "lastweek":
		lastSaturday := helperFuncs.ReturnLastWeekSaturday(time.Now())
		if templComp, ok := weekStatCache[lastSaturday]; ok {
			w.Header().Set(HeaderKey, lastSaturday)
			if err = templComp.Render(context.TODO(), w); err != nil {
				a.serverError(w, err)
			}
			return
		}
		msg = types.Message{
			Endpoint:        endpoint,
			WeekStatRequest: lastSaturday,
		}
		cachedSaturday = lastSaturday

	case "month":
		var firstSaturdayOfTheMonth, q string
		if q = r.URL.Query().Get("month"); q == "" {
			a.clientError(w, http.StatusBadRequest, errors.New("query param:month: cannot be empty"))
			return
		}
		if firstSaturdayOfTheMonth = helperFuncs.FirstSaturdayOfTheMonth(q); firstSaturdayOfTheMonth == "" {
			a.clientError(w, http.StatusBadRequest, errors.New("query param:month: invalid"))
			return
		}
		if templComp, ok := weekStatCache[firstSaturdayOfTheMonth]; ok {
			w.Header().Set("saturday", firstSaturdayOfTheMonth)
			if err = templComp.Render(context.TODO(), w); err != nil {
				a.serverError(w, err)
			}
			return
		}

		msg = types.Message{
			Endpoint:        endpoint,
			WeekStatRequest: firstSaturdayOfTheMonth,
		}
		cachedSaturday = firstSaturdayOfTheMonth

	case "backward-arrow", "forward-arrow":
		var displayedWeek, saturday string
		var t time.Time

		if displayedWeek = r.Header.Get(saturday); displayedWeek == "" {
			a.clientError(w, http.StatusBadRequest, errors.New("missing header saturday"))
			return
		}
		if t, err = time.Parse(types.TimeFormat, displayedWeek); err != nil {
			a.clientError(w, http.StatusBadRequest, errors.New("header value 'lastSaturday' invalide"))
			return
		}

		if query == "backward-arrow" {
			saturday = helperFuncs.ReturnLastWeekSaturday(t)
			msg = types.Message{
				Endpoint:        endpoint,
				WeekStatRequest: saturday,
			}
		}

		if query == "forward-arrow" {
			if helperFuncs.IsFutureDate(t) {
				weekStatCache["thisweek"].Render(context.TODO(), w)
				return
			}
			saturday = helperFuncs.ReturnNexWeektSaturday(t)
			msg = types.Message{
				Endpoint:        endpoint,
				WeekStatRequest: saturday,
			}
		}

		if templComp, ok := weekStatCache[saturday]; ok {
			w.Header().Set(HeaderKey, saturday)
			if err = templComp.Render(context.TODO(), w); err != nil {
				a.serverError(w, err)
			}
			return
		}
		cachedSaturday = saturday
	}

	fmt.Println("would be consulting the deamonservice")
	msg, err = a.writeAndReadWithDaemonService(msg)
	if err != nil {
		fmt.Println("error occurred in writeAndReadWithDaemonService", err)
		return
	}

	templComp := prepareHtTMLResponse(msg)

	// Cache
	// if query == "thisweek" {
	// 	weekStatCache[query] = templComp
	// } else {
	// 	weekStatCache[cachedSaturday] = templComp
	// }

	w.Header().Set(HeaderKey, cachedSaturday)
	if err = templComp.Render(context.TODO(), w); err != nil {
		w.Header().Del("lastSaturday")
		a.serverError(w, err)
	}
}

// func (a *App) renderFromCache(key string, w http.ResponseWriter) {
// 	if templComp, ok := weekStatCache[key]; ok {
// 		w.Header().Set(saturday, cachedSaturday)
// 		if err := templComp.Render(context.TODO(), w); err != nil {
// 			a.serverError(w, err)
// 		}
// 	}
// }
