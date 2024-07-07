package tasks

import (
	monitoring "LiScreMon/cli/daemon/internal/monitoring/linux"
	"fmt"
	"os/exec"
	"pkg/types"
	"reflect"
	"time"

	helperFuncs "pkg/helper"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

type TaskManagerDbRequirement interface {
	GetTaskByAppName(appName string) ([]types.Task, error)
	GetAllTask() ([]types.Task, error)
	RemoveTask(id uuid.UUID) error
	AddTask(task types.Task) error
	GetAppTodayActiveStatSoFar(appName string) (float64, error)
}

type TaskManager struct {
	dbHandle TaskManagerDbRequirement
	gocron   gocron.Scheduler
	channel  chan types.Task
}

func newTaskManger(dbHandle TaskManagerDbRequirement) *TaskManager {
	var tm TaskManager
	tm.dbHandle = dbHandle
	tm.gocron, _ = gocron.NewScheduler()
	tm.channel = make(chan types.Task)
	return &tm
}

func (tm *TaskManager) CloseChan() error {
	if err := tm.gocron.Shutdown(); err != nil {
		fmt.Println("error shutting down gocron Scheduler:", err)
		return err
	}
	tm.channel <- types.Task{}
	return nil
}

func (tm *TaskManager) SendTaskToTaskManager(task types.Task) error {
	if err := tm.dbHandle.AddTask(task); err != nil {
		return fmt.Errorf("error adding task to db: %w", err)
	}

	if reflect.ValueOf(task).IsZero() {
		return types.ErrZeroValueTask
	}

	tm.channel <- task
	return nil
}

func StartTaskManger(dbHandle TaskManagerDbRequirement) (*TaskManager, error) {

	tm := newTaskManger(dbHandle)

	tasks, err := tm.dbHandle.GetAllTask()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", types.ErrTaskMangerNotStarted.Error(), err)
	}

	go tm.disperseTask()

	for _, task := range tasks {

		switch {
		case task.Job != types.DailyAppLimit: // i.e reminders
			if task.Reminder.StartTime.Before(time.Now()) { // this is a double-check, the reminder should have been removed when done
				if err := tm.dbHandle.RemoveTask(task.UUID); err != nil {
					return nil, fmt.Errorf("%s %+v :err %w", types.ErrDeletingTask.Error(), task, err)
				}
			}

		case task.Job == types.DailyAppLimit: // this is a double-check, the oneTime app limit should have been removed when done
			if task.AppLimit.OneTime && task.AppLimit.Today != helperFuncs.Today() { // whether it reached limit or not,for a new day, the limit does not matter again, hence why it is a dailyLimit.
				if err := tm.dbHandle.RemoveTask(task.UUID); err != nil {
					return nil, fmt.Errorf("%s: %+v :%w", types.ErrDeletingTask.Error(), task, err)
				}
			}

			if task.AppLimit.IsLimitReached && task.AppLimit.Today == helperFuncs.Today() {
				// limit has been reached...and limit was reached that very day
				continue
			}
		}

		tm.channel <- task
	}

	return tm, nil
}

func (tm *TaskManager) disperseTask() {

	tm.gocron.Start()

	for {
		task := <-tm.channel

		if reflect.ValueOf(task).IsZero() {
			close(tm.channel)
			fmt.Println("closing and cleaning TaskManager")
			break
		}

		fmt.Printf("task received   %+v\n\n", task)

		switch task.Job {
		case types.ReminderWithNoAppLaunch:
			tm.reminderWithNoAppLaunch(task)

		case types.ReminderWithAppLaunch:
			tm.reminderWithAppLaunch(task)

		case types.DailyAppLimit:
			timeSofar, _ := tm.dbHandle.GetAppTodayActiveStatSoFar(task.AppName)
			monitoring.AddNewLimit(task, timeSofar)
		}

	}

}

func (tm *TaskManager) reminderWithNoAppLaunch(task types.Task) {

	tm.preNotify(task)

	if _, err := tm.gocron.NewJob(
		gocron.OneTimeJob(gocron.OneTimeJobStartDateTime(task.Reminder.StartTime)),
		gocron.NewTask(reminderAlert, task.UI, true),
		gocron.WithTags(task.UUID.String()),
		gocron.WithEventListeners(
			gocron.AfterJobRuns(func(jobID uuid.UUID, jobName string) {
				tm.dbHandle.RemoveTask(task.UUID)

			}),
		),
	); err != nil {
		fmt.Println("error creating job", err)
	}

}

func (tm *TaskManager) reminderWithAppLaunch(task types.Task) {
	tm.preNotify(task)

	tm.gocron.NewJob(
		gocron.OneTimeJob(gocron.OneTimeJobStartDateTime(task.Reminder.StartTime)),
		gocron.NewTask(reminderAlert, task.UI, true),
		gocron.WithTags(task.UUID.String()),
		gocron.WithEventListeners(
			gocron.AfterJobRuns(
				func(jobID uuid.UUID, jobName string) {
					cmd := exec.Command("bash", "-c", task.CmdLine)
					err := cmd.Start()
					if err != nil {
						fmt.Println(err)
					}

					err = cmd.Wait()
					if err != nil {
						fmt.Printf("Command finished with error: %v", err)
					}
				}),
			gocron.AfterJobRuns(func(jobID uuid.UUID, jobName string) {
				tm.dbHandle.RemoveTask(task.UUID)
			}),
		),
	)
}

func (tm *TaskManager) preNotify(task types.Task) {

	for i := 0; i < 2; i++ {
		var t time.Time
		notifyBeForeReminder, withSound := task.Reminder.AlertTimesInMinutes[i], task.Reminder.AlertSound[i]

		if t = task.Reminder.StartTime.Add(-time.Duration(notifyBeForeReminder) * time.Minute); t.Before(time.Now()) {
			continue // reminder is the past, useless
		}

		if _, err := tm.gocron.NewJob(
			gocron.OneTimeJob(gocron.OneTimeJobStartDateTime(t)),
			gocron.NewTask(preNotifyAlert, task.UI.Title, notifyBeForeReminder, withSound),
			gocron.WithTags(task.UUID.String())); err != nil {
			fmt.Println("gocron failed to add notififcation", err)
		}
	}
}

func preNotifyAlert(taskTitle string, durationbeforeTask int, withSound bool) {
	title := fmt.Sprintf("%d minutes to your task", durationbeforeTask)

	if withSound {
		helperFuncs.NotifyWithBeep(title, taskTitle)
		return
	}

	helperFuncs.NotifyWithoutBeep(title, taskTitle)
}

func reminderAlert(task types.UItextInfo, withSound bool) {
	title := fmt.Sprintf("Task: %s", task.Title)

	if withSound {
		helperFuncs.NotifyWithBeep(title, task.Subtitle)
		return
	}

	helperFuncs.NotifyWithoutBeep(title, task.Subtitle)
}

func (tm *TaskManager) RemoveTask(taskUUID uuid.UUID) {
	tm.gocron.RemoveByTags(taskUUID.String())
}
