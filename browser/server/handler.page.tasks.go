package webserver

import (
	"context"
	"fmt"
	"net/http"
	"pkg/types"
	"strconv"
	"strings"
	"time"

	views "views/components"

	"github.com/a-h/templ"
	"github.com/google/uuid"

	helperFuncs "pkg/helper"
)

func (a *App) tasksPage(w http.ResponseWriter, r *http.Request) {

	var err error
	msg := types.Message{
		Endpoint: strings.TrimPrefix(r.URL.Path, "/"),
	}

	msg, err = a.commWithDaemonService(msg)
	if err != nil {
		a.serverError(w, err)
		return
	}

	views.TasksPage("", msg.ReminderAndLimitResponse.AllApps).Render(context.TODO(), w)
}

func (a *App) ReminderTasksHandler(w http.ResponseWriter, r *http.Request) {

	var err error
	msg := types.Message{
		Endpoint: strings.TrimPrefix(r.URL.Path, "/"),
	}

	msg, err = a.commWithDaemonService(msg)
	if err != nil {
		a.serverError(w, err)
		return
	}

	reminderTasks := msg.ReminderAndLimitResponse.AllTask
	if len(reminderTasks) == 0 {
		views.RenderTasks(true, templ.NopComponent).Render(context.TODO(), w)
		return
	}

	c := views.ReminderTasks(reminderTasks)
	views.RenderTasks(false, c).Render(context.TODO(), w)
}

func (a *App) appLimitTasksHandler(w http.ResponseWriter, r *http.Request) {

	var err error
	msg := types.Message{
		Endpoint: strings.TrimPrefix(r.URL.Path, "/"),
	}

	msg, err = a.commWithDaemonService(msg)
	if err != nil {
		a.serverError(w, err)
		return
	}

	reminderTasks := msg.ReminderAndLimitResponse.AllTask
	if len(reminderTasks) == 0 {
		views.RenderTasks(true, templ.NopComponent).Render(context.TODO(), w)
		return
	}

	c := views.AppLimitTasks(reminderTasks)
	views.RenderTasks(false, c).Render(context.TODO(), w)
}

func (a *App) newReminderHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		a.clientError(w, http.StatusBadRequest, err)
		return
	}

	var task types.Task

	for key, value := range r.Form {
		switch key {
		case "title":
			task.UI.Title = value[0]

		case "subtitle":
			task.UI.Subtitle = value[0]

		case "reminder":
			startTime, err := time.ParseInLocation("2006-01-02T15:04", value[0], time.Local)
			if err != nil {
				a.clientError(w, http.StatusBadRequest, err)
				return
			}
			task.Reminder.StartTime = startTime

		case "soundNotification":
			if _, err := strconv.ParseBool(value[0]); err != nil {
				a.clientError(w, http.StatusBadRequest, fmt.Errorf("error parsing formData:%w", err))
				return
			}
			task.Reminder.AlertSound[2] = true

		case "firstNotification":
			alert, err := strconv.Atoi(value[0])
			if err != nil {
				a.clientError(w, http.StatusBadRequest, fmt.Errorf("error parsing formData:%w", err))
				return
			}
			task.Reminder.AlertTimesInMinutes[0] = alert

		case "soundFirstNotification":
			if _, err := strconv.ParseBool(value[0]); err != nil {
				a.clientError(w, http.StatusBadRequest, fmt.Errorf("error parsing formData:%w", err))
				return
			}
			task.Reminder.AlertSound[0] = true

		case "secondNotification":
			alert, err := strconv.Atoi(value[0])
			if err != nil {
				a.clientError(w, http.StatusBadRequest, fmt.Errorf("error parsing formData:%w", err))
				return
			}
			task.Reminder.AlertTimesInMinutes[1] = alert

		case "soundSecondNotification":
			if _, err := strconv.ParseBool(value[0]); err != nil {
				a.clientError(w, http.StatusBadRequest, fmt.Errorf("error parsing formData:%w", err))
				return
			}
			task.Reminder.AlertSound[1] = true

		case "app":
			if task.AppName = value[0]; value[0] == "no-app" {
				task.Job = types.ReminderWithNoAppLaunch
			} else {
				task.Job = types.ReminderWithAppLaunch
			}

		case "note":
			task.UI.Notes = value[0]
		}
	}

	task.UUID = uuid.New()
	msg := types.Message{
		Endpoint:                strings.TrimPrefix(r.URL.Path, "/"),
		ReminderAndLimitRequest: task,
	}

	if _, err = a.commWithDaemonService(msg); err != nil {
		a.serverError(w, err)
		return
	}

	views.TasksPage("", msg.ReminderAndLimitResponse.AllApps).Render(context.TODO(), w)
}

func (a *App) newAppLimitHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		a.clientError(w, http.StatusBadRequest, err)
		return
	}

	var (
		msg      types.Message
		task     types.Task
		hrs, min int
	)

	for key, value := range r.Form {
		switch key {
		case "app":
			task.AppName = value[0]

		case "hrs":
			if value[0] == "" {
				hrs = 0
			} else {
				hrs, err = strconv.Atoi(value[0])
				if err != nil {
					a.clientError(w, http.StatusBadRequest, fmt.Errorf("error parsing formData:%w", err))
					return
				}
			}

		case "min":
			if value[0] == "" {
				min = 0
			} else {
				min, err = strconv.Atoi(value[0])
				if err != nil {
					a.clientError(w, http.StatusBadRequest, fmt.Errorf("error parsing formData:%w", err))
					return
				}
			}

		case "recurring":
			isEveryDay, err := strconv.ParseBool(value[0])
			if err != nil {
				a.clientError(w, http.StatusBadRequest, fmt.Errorf("error parsing formData:%w", err))
				return
			}
			if isEveryDay {
				task.AppLimit.OneTime = false
				task.Job = types.DailyAppLimit
			}

		case "exitApp":
			exitApp, err := strconv.ParseBool(value[0])
			if err != nil {
				a.clientError(w, http.StatusBadRequest, fmt.Errorf("error parsing formData:%w", err))
				return
			}
			if exitApp {
				task.AppLimit.ExitApp = true
			}

		}
	}

	hours, minutes := time.Duration(hrs)*time.Hour, time.Duration(min)*time.Minute

	task.AppLimit.Limit = hours.Hours() + minutes.Hours()
	task.AppLimit.Today = helperFuncs.Today()
	task.UUID = uuid.New()

	msg = types.Message{
		Endpoint:                strings.TrimPrefix(r.URL.Path, "/"),
		ReminderAndLimitRequest: task,
	}

	if _, err = a.commWithDaemonService(msg); err != nil {
		if strings.Contains(err.Error(), types.ErrLimitAppExist.Error()) {
		}

		a.serverError(w, err)
		return
	}

	views.TasksPage("AppLimit", msg.ReminderAndLimitResponse.AllApps).Render(context.TODO(), w)

}

func (a *App) removeTask(w http.ResponseWriter, r *http.Request) {

	var err error

	path := strings.Split(r.URL.Path, "/")
	if !(len(path) > 2) {
		a.clientError(w, http.StatusBadRequest, fmt.Errorf("error parsing url"))
		return
	}

	endpoint, pathParam := path[1], r.PathValue("uuid")

	msg := types.Message{
		Endpoint: endpoint,
	}

	if msg.ReminderAndLimitRequest.UUID, err = uuid.Parse(pathParam); err != nil {
		a.clientError(w, http.StatusBadRequest, fmt.Errorf("bad task uuid format:%w", err))
		return
	}

	if msg, err = a.commWithDaemonService(msg); err != nil {
		a.serverError(w, err)
		return
	}

	if r.URL.Query().Get("taskType") == "AppLimit" {
		http.Redirect(w, r, "/appLimits", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/reminders", http.StatusSeeOther)
}
