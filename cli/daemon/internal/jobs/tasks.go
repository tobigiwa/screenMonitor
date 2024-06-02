package jobs

import (
	"fmt"
	"log"
	"os/exec"
	"pkg/types"
	"reflect"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

type TaskManagerDbRequirement interface {
	GetTaskByAppName(appName string) ([]types.Task, error)
	GetAllTask() ([]types.Task, error)
}

func StartTaskManger(dbHandle TaskManagerDbRequirement) (chan<- types.Task, error) {
	tm := NewTaskManger(dbHandle)

	tasks, err := tm.dbHandle.GetAllTask()
	if err != nil {
		return nil, fmt.Errorf("taskManager cannot be started: %w", err)
	}

	go tm.disperseTask()

	for _, task := range tasks {
		if task.TaskTime.StartTime.Before(time.Now()) {
			
			continue
		}
		tm.Chan <- task
	}

	return tm.Chan, nil
}

type TaskManager struct {
	dbHandle TaskManagerDbRequirement
	gocron   gocron.Scheduler
	Chan     chan types.Task
}

func NewTaskManger(dbHandle TaskManagerDbRequirement) *TaskManager {
	var tm TaskManager
	tm.dbHandle = dbHandle
	tm.gocron, _ = gocron.NewScheduler()
	tm.Chan = make(chan types.Task)
	return &tm
}

func (tm *TaskManager) Close() {
	if err := tm.gocron.Shutdown(); err != nil {
		fmt.Println("error shutting down gocron Scheduler:", err)
	}
	tm.Chan <- types.Task{}
}

func (tm *TaskManager) disperseTask() {

	tm.gocron.Start()
	for {
		task := <-tm.Chan

		if reflect.ValueOf(task).IsZero() {
			close(tm.Chan)
			break
		}

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

	tm.gocron.NewJob(
		gocron.OneTimeJob(gocron.OneTimeJobStartDateTime(task.TaskTime.StartTime)),
		gocron.NewTask(reminderFunc, task.UI, true),
	)

}

func (tm *TaskManager) createRemidersWithAction(task types.Task) {
	tm.reminders(task)

	tm.gocron.NewJob(
		gocron.OneTimeJob(gocron.OneTimeJobStartDateTime(task.TaskTime.StartTime)),
		gocron.NewTask(reminderFunc, task.UI, true),
		gocron.WithEventListeners(
			gocron.AfterJobRuns(
				func(jobID uuid.UUID, jobName string) {
					cmd := exec.Command("bash", "-c", task.AppCmdLine)
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
		notifyBeForeReminder, withSound := task.TaskTime.AlertTimesInMinutes[i], task.TaskTime.AlertSound[i]
		t := task.TaskTime.StartTime.Add(-time.Duration(notifyBeForeReminder) * time.Minute)

		if j, err := tm.gocron.NewJob(
			gocron.OneTimeJob(gocron.OneTimeJobStartDateTime(t)),
			gocron.NewTask(reminderFunc, task.UI, withSound)); err != nil {
			fmt.Println("gocron failed to add notififcation", j.ID())
		}
	}
}

func reminderFunc(task types.UItextInfo, withSound bool) {
	title := "Reminder: task.UI.Title"
	if withSound {
		beeep.Alert(title, task.Subtitle, "")
		return
	}
	beeep.Notify(title, task.Subtitle, "")
}
