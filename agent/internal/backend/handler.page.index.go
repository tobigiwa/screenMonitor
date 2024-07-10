package backend

import (
	views "agent/internal/frontend/components"
	"context"
	"log/slog"
	"net/http"
)

func (a *App) IndexPageHandler(w http.ResponseWriter, r *http.Request) {
	if err := views.IndexPage().Render(context.TODO(), w); err != nil {
		a.logger.Log(context.TODO(), slog.LevelError, err.Error())
		return
	}
}
