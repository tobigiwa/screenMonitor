package webserver

import (
	"context"
	"net/http"
	"pkg/types"
	"time"

	views "views/components"

	"github.com/google/uuid"
)

func (a *App) ReminderAndAlertPageHandler(w http.ResponseWriter, r *http.Request) {

	tasks := []types.Task{
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
	s1 := tasks
	s2 := tasks
	s3 := tasks
	s4 := tasks
	combined := append(s1, s2...)
	combined = append(combined, s3...)
	combined = append(combined, s4...)
	views.ReminderAndAlertPage(combined).Render(context.TODO(), w)
}
