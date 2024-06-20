package webserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"pkg/chart"
	"pkg/types"
	"strings"
	"time"

	helperFuncs "pkg/helper"

	views "views/components"

	"github.com/a-h/templ"
)

var lastRequestSaturday string

func (a *App) WeekStatHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("week")
	endpoint := strings.TrimPrefix(r.URL.Path, "/")

	var (
		msg types.Message
		err error
	)

	switch query {
	case "thisweek":
		msg = types.Message{
			Endpoint:        endpoint,
			WeekStatRequest: helperFuncs.SaturdayOfTheWeek(time.Now()),
		}

	case "lastweek":
		msg = types.Message{
			Endpoint:        endpoint,
			WeekStatRequest: helperFuncs.ReturnLastWeekSaturday(time.Now()),
		}

	case "month":
		var firstSaturdayOfTheMonth types.Date
		var q string
		if q = r.URL.Query().Get("month"); q == "" {
			a.clientError(w, http.StatusBadRequest, errors.New("query param:month: cannot be empty"))
			return
		}
		if firstSaturdayOfTheMonth = helperFuncs.FirstSaturdayOfTheMonth(q); firstSaturdayOfTheMonth == "" {
			a.clientError(w, http.StatusBadRequest, errors.New("query param:month: invalid data"))
			return
		}
		msg = types.Message{
			Endpoint:        endpoint,
			WeekStatRequest: firstSaturdayOfTheMonth,
		}

	case "backward", "forward":
		var t time.Time
		if t, err = time.Parse(types.TimeFormat, lastRequestSaturday); err != nil {
			a.clientError(w, http.StatusBadRequest, errors.New("header value 'lastSaturday' invalide"))
			return
		}

		if query == "backward" {
			msg = types.Message{
				Endpoint:        endpoint,
				WeekStatRequest: helperFuncs.ReturnLastWeekSaturday(t),
			}
		}

		if query == "forward" {
			msg = types.Message{
				Endpoint:        endpoint,
				WeekStatRequest: helperFuncs.ReturnNexWeektSaturday(t),
			}

			if helperFuncs.IsFutureDate(t) {
				msg = types.Message{
					Endpoint:        endpoint,
					WeekStatRequest: helperFuncs.SaturdayOfTheWeek(time.Now()), // show current week
				}
			}
		}
	}

	if msg, err = a.commWithDaemonService(msg); err != nil {
		a.serverError(w, fmt.Errorf("error occurred in commWithDaemonService:%w", err))
		return
	}

	if err = weekStatResponse(msg.WeekStatResponse).Render(context.TODO(), w); err != nil {
		a.serverError(w, err)
	}
	lastRequestSaturday = msg.WeekStatResponse.Keys[6]
}

func weekStatResponse(w types.WeekStatMessage) templ.Component {
	return views.WeekStatTempl(
		chart.WeekStatBarChart(chart.BarChartData{
			XAxis:       w.FormattedDay[:],
			YAxis:       w.Values[:],
			Keys:        w.Keys[:],
			Month:       w.Month,
			Year:        w.Year,
			TotalUptime: w.TotalWeekUptime,
		}),
		w.TotalWeekUptime,
		w.AppDetail,
		w.AllCategory,
		w.Keys[6],
	)
}
