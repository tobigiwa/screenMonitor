package database

import (
	"errors"
	"fmt"
	helperFuncs "pkg/helper"
	"pkg/types"

	"slices"

	badger "github.com/dgraph-io/badger/v4"
	"github.com/google/uuid"
)

func (bs *BadgerDBStore) getAllTasks() ([]types.Task, error) {
	byteData, err := bs.Get(dbTaskKey)

	if err != nil {
		if errors.Is(err, badger.ErrKeyNotFound) {
			if err = bs.db.Update(func(txn *badger.Txn) error {
				return txn.Set(dbTaskKey, []byte{})
			}); err != nil {
				return nil, err
			}
		}

		return nil, err
	}

	if len(byteData) == 0 {
		return []types.Task{}, nil
	}
	return helperFuncs.DecodeJSON[[]types.Task](byteData)
}

func (bs *BadgerDBStore) GetAllTask() ([]types.Task, error) {
	return bs.GetTaskByAppName("all")
}

func (bs *BadgerDBStore) GetTaskByAppName(appName string) ([]types.Task, error) {
	if appName == "" {
		return nil, errors.New("appName is empty")
	}

	taskArry, err := bs.getAllTasks()
	if err != nil {
		fmt.Println("error came from here:", err)
		return nil, err
	}

	if appName == "all" {
		return taskArry, nil
	}

	requestedTaskArr := make([]types.Task, 0, len(taskArry))
	for i := 0; i < len(taskArry); i++ {
		if taskArry[i].AppName == appName {
			requestedTaskArr = append(requestedTaskArr, taskArry[i])
		}
	}
	return slices.Clip(requestedTaskArr), nil
}

func (bs *BadgerDBStore) GetTaskByUUID(taskID uuid.UUID) (types.Task, error) {

	taskArry, err := bs.getAllTasks()
	if err != nil {
		fmt.Println("error came from here:", err)
		return types.Task{}, err
	}

	for _, task := range taskArry {
		if task.UUID == taskID {
			return task, nil
		}
	}

	return types.Task{}, errors.New("task does not exist")
}

func (bs *BadgerDBStore) AddTask(task types.Task) error {
	taskArry, err := bs.getAllTasks()
	if err != nil {
		return err
	}

	if bs.checkIfLimitAppExist(task, taskArry) {
		fmt.Println("this happened")
		return types.ErrLimitAppExist
	}

	taskArry = append(taskArry, task)

	byteData, err := helperFuncs.EncodeJSON(taskArry)
	if err != nil {
		return err
	}
	return bs.setOrUpdateKeyValue(dbTaskKey, byteData)
}

func (bs BadgerDBStore) checkIfLimitAppExist(task types.Task, tasks []types.Task) bool {

	for i := 0; i < len(tasks); i++ {
		if tasks[i].Job == types.DailyAppLimit {
			if tasks[i].AppName == task.AppName {
				return true
			}
		}
	}
	return false
}

func (bs *BadgerDBStore) RemoveTask(id uuid.UUID) error {
	taskArray, err := bs.getAllTasks()
	if err != nil {
		return err
	}
	newTaskArray := slices.DeleteFunc(taskArray, func(s types.Task) bool {
		return s.UUID == id
	})

	byteData, err := helperFuncs.EncodeJSON(newTaskArray)
	if err != nil {
		return err
	}
	return bs.setOrUpdateKeyValue(dbTaskKey, byteData)
}

func (bs *BadgerDBStore) UpdateAppLimitStatus(taskID uuid.UUID) error {
	taskArry, err := bs.getAllTasks()
	if err != nil {
		return err
	}

	for i := 0; i < len(taskArry); i++ {
		if task := taskArry[i]; task.UUID == taskID {
			task.AppLimit.Today = helperFuncs.Today()
			task.AppLimit.IsLimitReached = true
			taskArry[i] = task
			break
		}
	}

	byteData, err := helperFuncs.EncodeJSON(taskArry)
	if err != nil {
		return err
	}

	return bs.setOrUpdateKeyValue(dbTaskKey, byteData)
}
