package webserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"pkg/types"
)

func (a *App) AppStatHandler(w http.ResponseWriter, r *http.Request) {
	queryRange := r.URL.Query().Get("range")
	appName := r.URL.Query().Get("appName")
	if queryRange == "" || appName == "" {
		a.clientError(w, http.StatusBadRequest, errors.New("query params cannot be empty"))
	}

	var (
		msg types.Message
		err error
	)

	switch queryRange {
	case "week":
		msg.Endpoint = queryRange
		msg.StringDataRequest = ""

	case "month":
	case "range":
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
