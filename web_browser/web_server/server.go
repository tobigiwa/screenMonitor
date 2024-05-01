package backend

import (
	"browser/views"
	"context"

	"log/slog"
	"net/http"
)

type App struct {
	logger *slog.Logger
}

func NewApp(logger *slog.Logger) *App {
	return &App{
		logger: logger,
	}
}

func (a App) ScreenTimePageHandler(w http.ResponseWriter, r *http.Request) {
	if err := views.ScreenTimePage().Render(context.TODO(), w); err != nil {
		a.logger.Log(context.TODO(), slog.LevelError, err.Error())
		return
	}
}
func (a App) sendScreenGraphData(w http.ResponseWriter, r *http.Request) {
	// daemon.App.Db.GetWeeklyScreenStats()
}

func (a *App) Routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /screentime", a.ScreenTimePageHandler)

	fs := http.FileServer(http.Dir("./web_browser/assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets", fs))
	return mux
}
