package webserver

import (
	"context"
	"log/slog"
	"net/http"
	views "views/components"
)

func (a *App) IndexPageHandler(w http.ResponseWriter, r *http.Request) {
	if err := views.IndexPage().Render(context.TODO(), w); err != nil {
		a.logger.Log(context.TODO(), slog.LevelError, err.Error())
		return
	}
}
