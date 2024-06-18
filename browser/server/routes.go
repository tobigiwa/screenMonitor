package webserver

import "net/http"

func (a *App) Routes() *http.ServeMux {
	mux := http.NewServeMux()

	tmp := http.FileServer(http.Dir("/tmp/"))
	mux.Handle("/tmp/", http.StripPrefix("/tmp", tmp))

	fs := http.FileServer(http.Dir("../frontend/assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets", fs))

	// screentimePage
	mux.HandleFunc("GET /screentime", a.IndexPageHandler)
	mux.HandleFunc("GET /weekStat", a.WeekStatHandler)
	mux.HandleFunc("GET /appStat", a.AppStatHandler)
	mux.HandleFunc("GET /dayStat", a.DayStatHandler)
	mux.HandleFunc("POST /setCategory", a.SetCategory)

	// ReminderAndLimitPage
	mux.HandleFunc("GET /tasks", a.ReminderAndAlertPageHandler)
	mux.HandleFunc("GET /reminderTasks", a.AllReminderTask)
	mux.HandleFunc("GET /limitTasks", a.AllLimitTask)
	mux.HandleFunc("POST /createReminder", a.CreateReminderHandler)
	mux.HandleFunc("POST /createLimit", a.CreateLimitHandler)

	return mux
}
