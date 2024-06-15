package webserver

import (
	"context"
	"fmt"
	"net/http"
	"pkg/types"
	"slices"
	"strconv"
	"strings"
	"time"

	views "views/components"

	"github.com/google/uuid"
)

func (a *App) ReminderAndAlertPageHandler(w http.ResponseWriter, r *http.Request) {
	msg := types.Message{
		Endpoint: "allReminderTask",
	}
	msg, err := a.commWithDaemonService(msg)
	if err != nil {
		a.serverError(w, err)
		return
	}

	slices.SortFunc(msg.ReminderAndLimitResponse.AllTask, func(a, b types.Task) int {
		return a.TaskTime.StartTime.Compare(b.TaskTime.StartTime)
	})
	views.ReminderAndAlertPage(true, msg.ReminderAndLimitResponse.AllTask).Render(context.TODO(), w)
}

func (a *App) CreateReminderHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	var task types.Task

	for key, value := range r.Form {
		switch key {
		case "title":
			task.UI.Title = value[0]

		case "subtitle":
			task.UI.Subtitle = value[0]

		case "reminder":
			startTime, err := time.ParseInLocation("2006-01-02T15:04", value[0], time.Local)
			fmt.Println(startTime, value[0])
			if err != nil {
				a.clientError(w, http.StatusBadRequest, err)
				return
			}
			task.TaskTime.StartTime = startTime

		case "soundNotification":
			if _, err := strconv.ParseBool(value[0]); err != nil {
				a.clientError(w, http.StatusBadRequest, fmt.Errorf("error parsing formData:%w", err))
				return
			}
			task.TaskTime.AlertSound[2] = true

		case "firstNotification":
			alert, err := strconv.Atoi(value[0])
			if err != nil {
				a.clientError(w, http.StatusBadRequest, fmt.Errorf("error parsing formData:%w", err))
				return
			}
			task.TaskTime.AlertTimesInMinutes[0] = alert

		case "soundFirstNotification":
			if _, err := strconv.ParseBool(value[0]); err != nil {
				a.clientError(w, http.StatusBadRequest, fmt.Errorf("error parsing formData:%w", err))
				return
			}
			task.TaskTime.AlertSound[0] = true

		case "secondNotification":
			alert, err := strconv.Atoi(value[0])
			if err != nil {
				a.clientError(w, http.StatusBadRequest, fmt.Errorf("error parsing formData:%w", err))
				return
			}
			task.TaskTime.AlertTimesInMinutes[1] = alert

		case "soundSecondNotification":
			if _, err := strconv.ParseBool(value[0]); err != nil {
				a.clientError(w, http.StatusBadRequest, fmt.Errorf("error parsing formData:%w", err))
				return
			}
			task.TaskTime.AlertSound[1] = true

		case "app":
			if task.AppName = value[0]; value[0] == "no-app" {
				task.Job = types.ReminderWithNoAction
			} else {
				task.Job = types.ReminderWithAction
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
	res, err := a.commWithDaemonService(msg)
	if err != nil {
		a.serverError(w, err)
		return
	}
	if !res.ReminderAndLimitResponse.CreatedNewTask {
		a.serverError(w, fmt.Errorf("error creating reminder"))
		return
	}
	http.Redirect(w, r, "/reminder", http.StatusSeeOther)
}

func (a *App) CreateLimitHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	var (
		msg      types.Message
		task     types.Task
		hrs, min int
		err      error
	)

	for key, value := range r.Form {
		switch key {
		case "app":
			task.AppName = value[0]
		case "hrs":
			hrs, err = strconv.Atoi(value[0])
			if err != nil {
				a.clientError(w, http.StatusBadRequest, fmt.Errorf("error parsing formData:%w", err))
				return
			}
		case "min":
			min, err = strconv.Atoi(value[0])
			if err != nil {
				a.clientError(w, http.StatusBadRequest, fmt.Errorf("error parsing formData:%w", err))
				return
			}
		case "recurring":
			val, err := strconv.ParseBool(value[0])
			if err != nil {
				a.clientError(w, http.StatusBadRequest, fmt.Errorf("error parsing formData:%w", err))
				return
			}
			if val {
				task.TaskTime.EveryDay = true
				task.Job = types.Limit
			}
		case "exitApp":
			val, err := strconv.ParseBool(value[0])
			if err != nil {
				a.clientError(w, http.StatusBadRequest, fmt.Errorf("error parsing formData:%w", err))
				return
			}
			if val {
				task.TaskTime.ExitApp = true
			}

		}
	}

	hours, minutes := time.Duration(hrs)*time.Hour, time.Duration(min)*time.Minute

	task.TaskTime.Limit = hours.Hours() + minutes.Hours()

	task.UUID = uuid.New()

	fmt.Printf("\n%+v\n\n", task)

	msg = types.Message{
		Endpoint:                strings.TrimPrefix(r.URL.Path, "/"),
		ReminderAndLimitRequest: task,
	}

	res, err := a.commWithDaemonService(msg)
	if err != nil {
		a.serverError(w, err)
		return
	}
	if !res.ReminderAndLimitResponse.CreatedNewTask {
		a.serverError(w, fmt.Errorf("error creating reminder"))
		return
	}
	http.Redirect(w, r, "/reminder", http.StatusSeeOther)
}
