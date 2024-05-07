package webserver

import "net/http"

func (a *App) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	
	fs := http.FileServer(http.Dir("../frontend/assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets", fs))
	
	mux.HandleFunc("GET /screentime", a.ScreenTimePageHandler)
	mux.HandleFunc("GET /weekStat", a.WeekStat)

	return mux
}
