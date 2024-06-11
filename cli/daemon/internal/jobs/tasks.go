package jobs

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"pkg/types"
	"reflect"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

var appLogo = ""

type TaskManagerDbRequirement interface {
	GetTaskByAppName(appName string) ([]types.Task, error)
	GetAllTask() ([]types.Task, error)
	RemoveTask(id uuid.UUID) error
	AddTask(task types.Task) error
}

type TaskManager struct {
	dbHandle TaskManagerDbRequirement
	gocron   gocron.Scheduler
	channel  chan types.Task
}

func NewTaskManger(dbHandle TaskManagerDbRequirement) *TaskManager {
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
	tm.channel <- task
	return nil
}

func (tm *TaskManager) StartTaskManger() error {

	tasks, err := tm.dbHandle.GetAllTask()
	if err != nil {
		return fmt.Errorf("taskManager cannot be started: %w", err)
	}

	go tm.disperseTask()

	for _, task := range tasks {
		now, taskStartTime := time.Now(), task.TaskTime.StartTime
		if taskStartTime.Before(now) {
			if err := tm.dbHandle.RemoveTask(task.UUID); err != nil {
				return fmt.Errorf("err deleting old task: %+v :err %v", task, err)
			}
		}

		tm.channel <- task
	}

	return nil
}

func (tm *TaskManager) disperseTask() {

	// config directory
	homeDir, _ := os.UserHomeDir()
	configDir := homeDir + "/liScreMon/"
	appLogo = configDir + "liscremon.jpeg"

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
		case types.ReminderWithNoAction:
			tm.createRemidersWithNoAction(task)

		case types.ReminderWithAction:
			tm.createRemidersWithAction(task)

		case types.Limit:
		}
	}

}

func (tm *TaskManager) createRemidersWithNoAction(task types.Task) {
	tm.reminders(task)

	if _, err := tm.gocron.NewJob(
		gocron.OneTimeJob(gocron.OneTimeJobStartDateTime(task.TaskTime.StartTime)),
		gocron.NewTask(taskFunc, task.UI, true),
	); err != nil {
		fmt.Println("error creating job", err)
	}
}

func (tm *TaskManager) createRemidersWithAction(task types.Task) {
	tm.reminders(task)

	tm.gocron.NewJob(
		gocron.OneTimeJob(gocron.OneTimeJobStartDateTime(task.TaskTime.StartTime)),
		gocron.NewTask(taskFunc, task.UI, true),
		gocron.WithEventListeners(
			gocron.AfterJobRuns(
				func(jobID uuid.UUID, jobName string) {
					cmd := exec.Command("bash", "-c", task.AppInfo.CmdLine)
					err := cmd.Start()
					if err != nil {
						log.Println(err)
					}

					err = cmd.Wait()
					if err != nil {
						log.Printf("Command finished with error: %v", err)
					}
				}),
		),
	)
}

func (tm *TaskManager) reminders(task types.Task) {

	for i := 0; i < 2; i++ {
		var t time.Time
		notifyBeForeReminder, withSound := task.TaskTime.AlertTimesInMinutes[i], task.TaskTime.AlertSound[i]

		if t = task.TaskTime.StartTime.Add(-time.Duration(notifyBeForeReminder) * time.Minute); t.Before(time.Now()) {
			continue // reminder is the past, useless
		}

		if _, err := tm.gocron.NewJob(
			gocron.OneTimeJob(gocron.OneTimeJobStartDateTime(t)),
			gocron.NewTask(taskReminderFunc, task.UI.Title, notifyBeForeReminder, withSound)); err != nil {
			fmt.Println("gocron failed to add notififcation", err)
		}
	}
}

func taskReminderFunc(taskTitle string, durationbeforeTask int, withSound bool) {
	title := fmt.Sprintf("%d Minutes to your task", durationbeforeTask)
	if withSound {

		beeep.Alert(title, taskTitle, appLogo)
		return
	}
	beeep.Notify(title, taskTitle, appLogo)
}

func taskFunc(task types.UItextInfo, withSound bool) {
	title := fmt.Sprintf("Reminder: %s", task.Title)
	if withSound {

		beeep.Alert(title, task.Subtitle, appLogo)
		return
	}
	beeep.Notify(title, task.Subtitle, appLogo)
}
