package webserver

import "net/http"

func (a *App) Routes() *http.ServeMux {
	mux := http.NewServeMux()

	tmp := http.FileServer(http.Dir("/tmp/"))
	mux.Handle("/tmp/", http.StripPrefix("/tmp", tmp))

	fs := http.FileServer(http.Dir("../frontend/assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets", fs))

	mux.HandleFunc("GET /", a.IndexPageHandler)
	mux.HandleFunc("GET /weekStat", a.WeekStatHandler)
	mux.HandleFunc("GET /appStat", a.AppStatHandler)

	return mux
}
