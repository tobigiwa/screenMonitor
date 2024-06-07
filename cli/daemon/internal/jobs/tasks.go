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
	RemoveTask(id uuid.UUID) error
	AddTask(task types.Task) error
}

func (tm *TaskManager) StartTaskManger() error {

	tasks, err := tm.dbHandle.GetAllTask()
	if err != nil {
		return fmt.Errorf("taskManager cannot be started: %w", err)
	}

	go tm.disperseTask()

	for _, task := range tasks {
		if task.TaskTime.StartTime.Before(time.Now()) {
			if err := tm.dbHandle.RemoveTask(task.UUID); err != nil {
				fmt.Printf("err deleting old task: %+v : %v\n\n", task, err)
			}
			fmt.Printf("task is in thhe past deleted: %+v\n\n", task)
			continue
		}
		tm.channel <- task
	}

	return nil
}

func (tm *TaskManager) SendTaskToTaskManager(task types.Task) error {
	if err := tm.dbHandle.AddTask(task); err != nil {
		return fmt.Errorf("error adding task: %w", err)
	}
	fmt.Println("waiting for task to be added to task manager")
	tm.channel <- task
	return nil
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

func (tm *TaskManager) disperseTask() {

	tm.gocron.Start()
	for {
		task := <-tm.channel
		fmt.Printf("task received: %+v\n", task)

		if reflect.ValueOf(task).IsZero() {
			close(tm.channel)
			fmt.Println("closing and cleaning TaskManager")
			break
		}

		switch task.Job {
		case types.ReminderWithNoAction:
			fmt.Println("it got here")
			tm.createRemidersWithNoAction(task)

		case types.ReminderWithAction:
			tm.createRemidersWithAction(task)

		case types.Limit:
		}
		tm.gocron.Start()
	}

}

func (tm *TaskManager) createRemidersWithNoAction(task types.Task) {
	tm.reminders(task)
	fmt.Println("reminder fired")

	j, err := tm.gocron.NewJob(
		gocron.OneTimeJob(gocron.OneTimeJobStartDateTime(task.TaskTime.StartTime)),
		gocron.NewTask(reminderFunc, task.UI, true),
	)
	if err != nil {
		fmt.Println("error creating job", err)
	}
	fmt.Println(j.ID())
}

func (tm *TaskManager) createRemidersWithAction(task types.Task) {
	tm.reminders(task)

	tm.gocron.NewJob(
		gocron.OneTimeJob(gocron.OneTimeJobStartDateTime(task.TaskTime.StartTime)),
		gocron.NewTask(reminderFunc, task.UI, true),
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
		notifyBeForeReminder, withSound := task.TaskTime.AlertTimesInMinutes[i], task.TaskTime.AlertSound[i]
		t := task.TaskTime.StartTime.Add(-time.Duration(notifyBeForeReminder) * time.Minute)
		fmt.Print(t.Date())
		fmt.Print(t.Clock())
		fmt.Println()

		_, err := tm.gocron.NewJob(
			gocron.OneTimeJob(gocron.OneTimeJobStartDateTime(t)),
			gocron.NewTask(reminderFunc, task.UI, withSound))

		if err != nil {
			fmt.Println("gocron failed to add notififcation")
		}
	}
}

func reminderFunc(task types.UItextInfo, withSound bool) {
	title := "Reminder: task.UI.Title"
	if withSound {
		beeep.Alert(title, task.Subtitle, "")
		fmt.Println("OUR WORK WAS DONE")
		return
	}
	beeep.Notify(title, task.Subtitle, "")
	fmt.Println("OUR WORK WAS DONE")
}
