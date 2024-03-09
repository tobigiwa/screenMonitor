package store

import (
	"errors"
	"fmt"

	badger "github.com/dgraph-io/badger/v4"
)

type BadgerDBStore struct {
	db *badger.DB
}

func NewBadgerDb(pathToDb string) (*BadgerDBStore, error) {
	opts := badger.DefaultOptions(pathToDb)

	opts.Logger = nil
	badgerInstance, err := badger.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("opening kv: %w", err)
	}

	return &BadgerDBStore{db: badgerInstance}, nil
}

func (bs *BadgerDBStore) Close() error {
	return bs.db.Close()
}

func (bs *BadgerDBStore) get(key string) ([]byte, error) {
	var (
		valCopy []byte
		err     error
	)
	err = bs.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		valCopy, err = item.ValueCopy(nil)
		if err != nil {
			return err
		}
		return nil
	})
	return valCopy, err
}

func (bs *BadgerDBStore) set(key, value []byte) error {
	return bs.db.Update(
		func(txn *badger.Txn) error {
			return txn.Set([]byte(key), []byte(value))
		})
}

func (bs *BadgerDBStore) delete(key string) error {
	return bs.db.Update(
		func(txn *badger.Txn) error {
			return txn.Delete([]byte(key))
		})
}

func (bs *BadgerDBStore) WriteUsage(data ScreenTime) error {
	return bs.db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(data.AppName))
		var newApp bool
		if newApp = errors.Is(err, badger.ErrKeyNotFound); err != nil && !newApp {
			return err
		}

		var appInfo appInfo

		if newApp {
			fmt.Printf("new app :%v\n\n", data.AppName)
			appInfo.AppName = data.AppName
			appInfo.ActiveScreenStats = make(dailyAppScreenTime)
			appInfo.PassiveScreenStats = make(dailyAppScreenTime)
		}

		if err == nil {
			valCopy, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}
			if err := appInfo.deserialize(valCopy); err != nil {
				return err
			}
			if data.AppName != appInfo.AppName {
				return ErrAppKeyMismatch
			}
			fmt.Printf("existing appName:%v, time so far is: %v:%v\n\n", data.AppName, appInfo.ActiveScreenStats[Key()], appInfo.PassiveScreenStats[Key()])
		}

		if data.Type == Active {
			appInfo.ActiveScreenStats[Key()] += data.Time
		} else {
			appInfo.PassiveScreenStats[Key()] += data.Time
		}

		ser, err := appInfo.serialize()
		if err != nil {
			return err
		}

		return txn.Set([]byte(data.AppName), ser)
	})
}
