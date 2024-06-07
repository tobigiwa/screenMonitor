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
	byteData, err := bs.Get(dbTaskKey())
	if err != nil {
		if errors.Is(err, badger.ErrKeyNotFound) {
			if err = bs.db.Update(func(txn *badger.Txn) error {
				byteData, _ := helperFuncs.EncodeJSON([]types.Task{})
				return txn.Set(dbTaskKey(), byteData)
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
		if taskArry[i].AppInfo.AppName == appName {
			requestedTaskArr = append(requestedTaskArr, taskArry[i])
		}
	}
	return slices.Clip(requestedTaskArr), nil
}

func (bs *BadgerDBStore) AddTask(task types.Task) error {
	taskArry, err := bs.getAllTasks()
	if err != nil {
		return err
	}

	taskArry = append(taskArry, task)

	byteData, err := helperFuncs.EncodeJSON(taskArry)
	if err != nil {
		return err
	}
	return bs.updateKeyValue(dbTaskKey(), byteData)
}

func (bs *BadgerDBStore) RemoveTask(id uuid.UUID) error {
	taskArry, err := bs.getAllTasks()
	if err != nil {
		return err
	}
	allTasks := slices.DeleteFunc(taskArry, func(s types.Task) bool {
		return s.UUID == id
	})

	byteData, err := helperFuncs.EncodeJSON(allTasks)
	if err != nil {
		return err
	}
	return bs.updateKeyValue(dbTaskKey(), byteData)
}
