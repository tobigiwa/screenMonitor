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
	msg, err := a.writeAndReadWithDaemonService(msg)
	if err != nil {
		a.serverError(w, err)
		return
	}

	slices.SortFunc(msg.ReminderResponse.AllTask, func(a, b types.Task) int {
		return a.TaskTime.StartTime.Compare(b.TaskTime.StartTime)
	})
	views.ReminderAndAlertPage(msg.ReminderResponse.AllTask).Render(context.TODO(), w)
	// views.ReminderAndAlertPage(longTask()).Render(context.TODO(), w)
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
			if task.AppInfo.AppName = value[0]; value[0] == "no-app" {
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
		Endpoint:        strings.TrimPrefix(r.URL.Path, "/"),
		ReminderRequest: task,
	}
	res, err := a.writeAndReadWithDaemonService(msg)
	if err != nil {
		a.serverError(w, err)
		return
	}
	if !res.ReminderResponse.CreatedNewTask {
		a.serverError(w, fmt.Errorf("error creating reminder"))
		return
	}
	http.Redirect(w, r, "/reminder", http.StatusSeeOther)
}

var tasks = []types.Task{
	{
		UUID: uuid.New(),
		AppInfo: types.AppIconCategoryAndCmdLine{
			AppName: "djesfefef",
		},
		TaskTime: types.TaskTime{
			StartTime:           time.Now(),
			EndTime:             time.Now().Add(1 * time.Hour),
			AlertTimesInMinutes: [3]int{15, 30, 45},
			AlertSound:          [3]bool{true, true, false},
		},
		UI: types.UItextInfo{
			Title:    "Task Title 1",
			Subtitle: "Task Subtitle 1",
			Notes:    "Task Notes 1",
		},
		Job: "TaskType1", // Replace with actual TaskType
	},
	{
		UUID: uuid.New(),
		AppInfo: types.AppIconCategoryAndCmdLine{
			AppName: "djesfefef",
		},
		TaskTime: types.TaskTime{
			StartTime:           time.Now().Add(2 * time.Hour),
			EndTime:             time.Now().Add(3 * time.Hour),
			AlertTimesInMinutes: [3]int{20, 40, 60},
			AlertSound:          [3]bool{false, true, false},
		},
		UI: types.UItextInfo{
			Title:    "Task Title 2",
			Subtitle: "Task Subtitle 2",
			Notes:    "Task Notes 2",
		},
		Job: "TaskType2", // Replace with actual TaskType
	},
	{
		UUID: uuid.New(),
		AppInfo: types.AppIconCategoryAndCmdLine{
			AppName: "djesfefef",
		},
		TaskTime: types.TaskTime{
			StartTime:           time.Now().Add(4 * time.Hour),
			EndTime:             time.Now().Add(5 * time.Hour),
			AlertTimesInMinutes: [3]int{25, 50, 75},
			AlertSound:          [3]bool{true, false, true},
		},
		UI: types.UItextInfo{
			Title:    "Task Title 3",
			Subtitle: "Task Subtitle 3",
			Notes:    "Task Notes 3",
		},
		Job: "TaskType3", // Replace with actual TaskType
	},
	{
		UUID: uuid.New(),
		AppInfo: types.AppIconCategoryAndCmdLine{
			AppName: "djesfefef",
		},
		TaskTime: types.TaskTime{
			StartTime:           time.Now().Add(6 * time.Hour),
			EndTime:             time.Now().Add(7 * time.Hour),
			AlertTimesInMinutes: [3]int{30, 60, 90},
			AlertSound:          [3]bool{false, false, false},
		},
		UI: types.UItextInfo{
			Title:    "Task Title 4",
			Subtitle: "Task Subtitle 4",
			Notes:    "Task Notes 4",
		},
		Job: "TaskType4", // Replace with actual TaskType
	},
}

func longTask() []types.Task {
	s1, s2, s3, s4 := tasks, tasks, tasks, tasks
	combined := append(s1, s2...)
	combined = append(combined, s3...)
	combined = append(combined, s4...)
	return combined
}

// title := r.Form.Get("title")
// subtitle := r.Form.Get("subtitle")
// reminder := r.Form.Get("reminder")
// soundNotification := r.Form.Get("soundNotification")
// firstNotification := r.Form.Get("firstNotification")
// soundNotificationFirst := r.Form.Get("soundFirstNotification")
// secondNotification := r.Form.Get("secondNotification")
// soundNotificationSecond := r.Form.Get("soundSecondNotification")
// app := r.Form.Get("app")
// note := r.Form.Get("note")

// startTime, err := time.Parse("2006-01-02T15:04", reminder)
// if err != nil {
// 	a.clientError(w, http.StatusBadRequest, err)
// 	return
// }
// alert1, err1 := strconv.Atoi(firstNotification)
// alert2, err2 := strconv.Atoi(secondNotification)
// if err1 != nil || err2 != nil {
// 	a.clientError(w, http.StatusBadRequest, fmt.Errorf("%w:%w", err1, err2))
// 	return
// }

// newTask := types.Task{
// 	UUID: uuid.New(),
// 	AppInfo: types.AppIconCategoryAndCmdLine{
// 		AppName: app,
// 	},
// 	TaskTime: types.TaskTime{
// 		StartTime:           startTime,
// 		AlertTimesInMinutes: [3]int{0, alert1, alert2},
// 		AlertSound:          [3]bool{},
// 	},
// }
