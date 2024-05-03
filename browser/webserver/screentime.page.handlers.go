package webserver

import (
	"context"
	"log/slog"
	"net/http"
	views "views"
)

func (a *App) ScreenTimePageHandler(w http.ResponseWriter, r *http.Request) {
	if err := views.ScreenTimePage().Render(context.TODO(), w); err != nil {
		a.logger.Log(context.TODO(), slog.LevelError, err.Error())
		return
	}
}
func (a *App) sendScreenGraphData(w http.ResponseWriter, r *http.Request) {

}
