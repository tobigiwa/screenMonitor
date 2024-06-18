package service

import (
	db "LiScreMon/cli/daemon/internal/database"
	"LiScreMon/cli/daemon/internal/jobs"
	"cmp"
	"fmt"
	helperFuncs "pkg/helper"
	"pkg/types"
	"slices"
	"strings"
	"time"
)

type Service struct {
	db          DatabaseInterface
	taskManager *jobs.TaskManager
}

func (s *Service) StopTaskManger() error {
	return s.taskManager.CloseChan()
}

func (s *Service) getWeekStat(msg types.Date) types.WeekStatMessage {
	var (
		weekStat    db.WeeklyStat
		appsInfo    []types.AppIconCategoryAndCmdLine
		allCategory []types.Category
		err         error
	)

	if weekStat, err = s.db.GetWeek(msg); err != nil {
		return types.WeekStatMessage{IsError: true, Error: fmt.Errorf("error weekStat: %w", err)}
	}

	var (
		keys             = [7]string{}
		formattedDay     = [7]string{}
		values           = [7]float64{}
		sizeOfApps       = len(weekStat.EachApp)
		appNameInTheWeek = make([]string, 0, sizeOfApps)
		appCard          = make([]types.ApplicationDetail, 0, sizeOfApps)
	)

	for i := 0; i < 7; i++ {
		keys[i] = string(weekStat.DayByDayTotal[i].Key)
		values[i] = weekStat.DayByDayTotal[i].Value.Active
		formattedDay[i] = helperFuncs.FormattedDay(types.Date(weekStat.DayByDayTotal[i].Key))
	}

	month, year := helperFuncs.MonthAndYear(types.Date(weekStat.DayByDayTotal[6].Key))

	for i := 0; i < sizeOfApps; i++ {
		appNameInTheWeek = append(appNameInTheWeek, weekStat.EachApp[i].AppName)
	}

	if appsInfo, err = s.db.GetAppIconCategoryAndCmdLine(appNameInTheWeek); err != nil {
		return types.WeekStatMessage{IsError: true, Error: fmt.Errorf("err with GetAppIconAndCategory:%w", err)}
	}

	for i := 0; i < sizeOfApps; i++ {
		appCard = append(appCard, types.ApplicationDetail{AppInfo: appsInfo[i], Usage: weekStat.EachApp[i].Usage.Active})
	}

	if allCategory, err = s.db.GetAllACategories(); err != nil {
		return types.WeekStatMessage{IsError: true, Error: fmt.Errorf("err with GetAllCategories:%w", err)}
	}

	return types.WeekStatMessage{
		Keys:            keys,
		FormattedDay:    formattedDay,
		Values:          values,
		TotalWeekUptime: weekStat.WeekTotal.Active,
		AllCategory:     allCategory,
		Month:           month,
		Year:            fmt.Sprint(year),
		AppDetail:       appCard,
	}
}

func (s *Service) getAppStat(msg types.AppStatRequest) types.AppStatMessage {
	var (
		appStat types.AppRangeStat
		err     error
	)

	switch msg.StatRange {
	case "week":
		appStat, err = s.db.AppWeeklyStat(msg.AppName, msg.Start)
	case "month":
		appStat, err = s.db.AppMonthlyStat(msg.AppName, msg.Month, msg.Year)
	case "range":
		appStat, err = s.db.AppDateRangeStat(msg.AppName, msg.Start, msg.End)
	}

	if err != nil {
		fmt.Println("error weekStat:", err)
		return types.AppStatMessage{IsError: true, Error: err}
	}

	var (
		formattedDay      = make([]string, 0, len(appStat.DaysRange))
		values            = []float64{}
		lastDayOfTheRange = len(appStat.DaysRange) - 1
	)

	for i := 0; i < len(appStat.DaysRange); i++ {
		formattedDay = append(formattedDay, helperFuncs.FormattedDay(types.Date(appStat.DaysRange[i].Key)))
		values = append(values, appStat.DaysRange[i].Value.Active)
	}
	month, year := helperFuncs.MonthAndYear(types.Date(appStat.DaysRange[lastDayOfTheRange].Key))

	return types.AppStatMessage{
		FormattedDay:     formattedDay,
		Values:           values,
		Month:            month,
		Year:             year,
		TotalRangeUptime: appStat.TotalRange.Active,
		AppInfo:          appStat.AppInfo,
	}
}

func (s *Service) getDayStat(msg types.Date) types.DayStatMessage {
	dayStat, err := s.db.GetDay(msg)
	if err != nil {
		return types.DayStatMessage{IsError: true, Error: err}
	}
	d, _ := helperFuncs.ParseKey(msg)
	date := fmt.Sprintf("%s. %s %s, %d", strings.TrimSuffix(d.Weekday().String(), "day"), helperFuncs.AddOrdinalSuffix(d.Day()), d.Month().String(), d.Year())

	return types.DayStatMessage{EachApp: dayStat.EachApp, DayTotal: dayStat.DayTotal, Date: date}
}

func (s *Service) setAppCategory(msg types.SetCategoryRequest) types.SetCategoryResponse {
	if err := s.db.SetAppCategory(msg.AppName, msg.Category); err != nil {
		return types.SetCategoryResponse{IsError: true, Error: err}
	}
	return types.SetCategoryResponse{IsCategorySet: true}
}

func (s *Service) tasks() types.ReminderMessage {

	allApps, err := s.db.GetAllApp()
	if err != nil {
		return types.ReminderMessage{IsError: true, Error: err}
	}
	return types.ReminderMessage{AllApps: allApps}
}

func (s *Service) reminderTasks() types.ReminderMessage {

	tasks, err := s.db.GetAllTask()
	if err != nil {
		return types.ReminderMessage{IsError: true, Error: err}
	}

	validTask := make([]types.Task, 0, len(tasks))
	for _, task := range tasks {

		if task.Job == types.Limit {
			continue
		}

		now, taskStartTime := time.Now(), task.TaskTime.StartTime

		if taskStartTime.Before(now) {
			if err := s.db.RemoveTask(task.UUID); err != nil {
				return types.ReminderMessage{IsError: true, Error: err}
			}
		}
		validTask = append(validTask, task)
	}

	slices.SortFunc(validTask, func(a, b types.Task) int {
		return a.TaskTime.StartTime.Compare(b.TaskTime.StartTime)
	})

	return types.ReminderMessage{AllTask: slices.Clip(validTask)}
}

func (s *Service) limitTasks() types.ReminderMessage {
	tasks, err := s.db.GetAllTask()
	if err != nil {
		return types.ReminderMessage{IsError: true, Error: err}
	}

	limitTask := make([]types.Task, 0, len(tasks))
	for _, task := range tasks {
		if task.Job == types.Limit {
			limitTask = append(limitTask, task)
		}
	}

	slices.SortFunc(limitTask, func(a, b types.Task) int {
		return cmp.Compare(a.TaskTime.Limit, b.TaskTime.Limit)
	})

	return types.ReminderMessage{AllTask: slices.Clip(limitTask)}
}

func (s *Service) addNewReminder(task types.Task) types.ReminderMessage {

	if task.Job == types.ReminderWithAction {
		appInfo, err := s.db.GetAppIconCategoryAndCmdLine([]string{task.AppName})
		if err != nil {
			return types.ReminderMessage{IsError: true, Error: err}
		}
		task.AppIconCategoryAndCmdLine = appInfo[0]
	}

	err := s.taskManager.SendTaskToTaskManager(task)
	if err != nil {
		return types.ReminderMessage{IsError: true, Error: err}
	}

	return types.ReminderMessage{
		CreatedNewTask: true,
	}
}

func (s *Service) addNewLimitApp(msg types.Task) types.ReminderMessage {

	err := s.taskManager.SendTaskToTaskManager(msg)
	if err != nil {
		return types.ReminderMessage{IsError: true, Error: err}
	}

	return types.ReminderMessage{CreatedNewTask: true}
}
