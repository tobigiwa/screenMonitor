package webserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"pkg/chart"
	helperFuncs "pkg/helper"
	"pkg/types"
	"strings"

	views "views/components"

	"github.com/a-h/templ"
)

func (a *App) DayStatHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("day")
	endpoint := strings.TrimPrefix(r.URL.Path, "/")

	if query == "" {
		a.clientError(w, http.StatusBadRequest, errors.New("query param:day: cannot be empty"))
		return
	}

	if !helperFuncs.ValidDateType(query) {
		a.clientError(w, http.StatusBadRequest, errors.New("query param:day: inavalid string format"))
		return
	}

	var (
		msg types.Message
		err error
	)

	msg = types.Message{
		Endpoint:       endpoint,
		DayStatRequest: types.Date(query),
	}

	if msg, err = a.commWithDaemonService(msg); err != nil {
		a.serverError(w, fmt.Errorf("error occurred in commWithDaemonService:%w", err))
	}

	if err = dayStatResponse(msg.DayStatResponse).Render(context.TODO(), w); err != nil {
		a.serverError(w, err)
	}
}

func dayStatResponse(w types.DayStatMessage) templ.Component {
	return views.DayStatTempl(
		chart.DayStatPieChart(
			chart.PieChartData{
				PieData:  w.EachApp,
				DayTotal: w.DayTotal,
				Date:     w.Date,
			}),
		nil,
	)
}
