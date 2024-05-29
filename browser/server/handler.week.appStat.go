package webserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"pkg/types"
	"strings"
)

func (a *App) AppStatHandler(w http.ResponseWriter, r *http.Request) {

	endpoint := strings.TrimPrefix(r.URL.Path, "/")
	queryRange := r.URL.Query().Get("range")
	appName := r.URL.Query().Get("appName")

	if queryRange == "" || appName == "" {
		a.clientError(w, http.StatusBadRequest, errors.New("query params cannot be empty"))
	}

	var (
		msg types.Message
		err error
	)

	msg.Endpoint = endpoint
	msg.AppStatRequest.AppName = appName

	switch queryRange {
	case "week":
		msg.AppStatRequest.StatRange = queryRange
		msg.AppStatRequest.Start = types.Date("2024-05-26")
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

	if msg, err = a.writeAndReadWithDaemonService(msg); err != nil {
		a.serverError(w, fmt.Errorf("error occurred in writeAndReadWithDaemonService:%w", err))
		return
	}

	templComp := prepareHtTMLResponse(msg)
	if err = templComp.Render(context.TODO(), w); err != nil {
		a.serverError(w, fmt.Errorf("templ reander error:%w", err))
	}

	templComp.Render(context.TODO(), w)
}
