package webserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	helperFuncs "pkg/helper"
	"pkg/types"
	"strings"
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

	if msg, err = a.writeAndReadWithDaemonService(msg); err != nil {
		a.serverError(w, fmt.Errorf("error occurred in writeAndReadWithDaemonService:%w", err))
	}

	templComp := prepareHtTMLResponse(msg)

	if err = templComp.Render(context.TODO(), w); err != nil {
		a.serverError(w, err)
	}

}
