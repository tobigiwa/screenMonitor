package database

import (
	"errors"
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
			_ = bs.db.Update(func(txn *badger.Txn) error {
				return txn.SetEntry(&badger.Entry{Key: dbTaskKey()})
			})
		}
		return nil, err
	}

	return helperFuncs.DecodeJSON[[]types.Task](byteData)
}

func (bs *BadgerDBStore) GetAllTask() ([]types.Task, error) {
	return bs.GetTaskByAppName("all")
}

func (bs *BadgerDBStore) RemoveTask(id uuid.UUID) error {
	taskArry, err := bs.getAllTasks()
	if err != nil {
		return err
	}
	slices.DeleteFunc(taskArry, func(s types.Task) bool {
		return s.UUID == id
	})

	return nil
}

func (bs *BadgerDBStore) GetTaskByAppName(appName string) ([]types.Task, error) {
	if appName == "" {
		return nil, errors.New("appName is empty")
	}

	taskArry, err := bs.getAllTasks()
	if err != nil {
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
