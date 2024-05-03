package webserver

import "net/http"

func (a *App) Routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /screentime", a.ScreenTimePageHandler)

	fs := http.FileServer(http.Dir("../frontend/assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets", fs))
	return mux
}
