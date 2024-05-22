package webserver

import "net/http"

func (a *App) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	
	tmp := http.FileServer(http.Dir("/tmp/"))
	mux.Handle("/tmp/", http.StripPrefix("/tmp", tmp))

	fs := http.FileServer(http.Dir("../frontend/assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets", fs))
	
	mux.HandleFunc("GET /screentime", a.ScreenTimePageHandler)
	mux.HandleFunc("GET /weekStat", a.WeekStat)

	return mux
}
