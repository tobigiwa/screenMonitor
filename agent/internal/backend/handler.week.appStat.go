package backend

import (
	"agent/internal/backend/chart"
	views "agent/internal/frontend/components"
	"context"
	"errors"
	"fmt"
	"net/http"
	"pkg/types"
	"strings"

	"github.com/a-h/templ"
)

func (a *App) AppStatHandler(w http.ResponseWriter, r *http.Request) {

	endpoint := strings.TrimPrefix(r.URL.Path, "/")
	queryRange := r.URL.Query().Get("range")
	appName := r.URL.Query().Get("appName")

	if queryRange == "" || appName == "" {
		a.clientError(w, http.StatusBadRequest, errors.New("query params cannot be empty"))
		return
	}

	var (
		msg types.Message
		err error
	)

	msg.Endpoint = endpoint
	msg.AppStatRequest.AppName = appName

	switch queryRange {
	case "week":
		var start string
		if start = r.URL.Query().Get("start"); start == "" {
			a.clientError(w, http.StatusBadRequest, errors.New("query params:start: cannot be empty"))
			return
		}
		msg.AppStatRequest.StatRange = queryRange
		msg.AppStatRequest.Start = types.Date(start)
	case "month":
		msg.AppStatRequest.StatRange = queryRange
		msg.AppStatRequest.Month = ""
		msg.AppStatRequest.Year = ""
	case "range":
		msg.AppStatRequest.StatRange = queryRange
		msg.AppStatRequest.Start = ""
		msg.AppStatRequest.End = ""
	default:
		a.clientError(w, http.StatusBadRequest, fmt.Errorf(`value of query param "range" is unexpected,value - %s`, queryRange))
		return
	}

	if msg, err = a.commWithDaemonService(msg); err != nil {
		a.serverError(w, fmt.Errorf("error occurred in commWithDaemonService:%w", err))
		return
	}

	if err = appStatResponse(msg.AppStatResponse).Render(context.TODO(), w); err != nil {
		a.serverError(w, fmt.Errorf("templ reander error:%w", err))
	}
}

func appStatResponse(w types.AppStatMessage) templ.Component {
	return views.AppStatTempl(
		chart.AppStatBarChart(
			chart.BarChartData{
				AppName:     w.AppInfo.AppName,
				YAxis:       w.Values,
				XAxis:       w.FormattedDay,
				Month:       w.Month,
				Year:        w.Year,
				TotalUptime: w.TotalRangeUptime,
			}),
		w.AppInfo)
}
