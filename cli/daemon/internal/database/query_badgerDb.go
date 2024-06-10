package database

import (
	"bytes"
	"fmt"

	helperFuncs "pkg/helper"

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

func (bs *BadgerDBStore) updateKeyValue(key, byteData []byte) error {
	err := bs.db.Update(func(txn *badger.Txn) error {
		err := txn.Set(key, byteData)
		return err
	})
	return err
}

func (bs *BadgerDBStore) DeleteKey(key []byte) error {
	return bs.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(key)
	})
}

func (bs *BadgerDBStore) Get(key []byte) ([]byte, error) {

	var valCopy []byte

	err := bs.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}

		err = item.Value(func(val []byte) error {
			valCopy = append([]byte{}, val...)
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return valCopy, nil
}

func (bs *BadgerDBStore) DeleteBucket(dbPrefix string) error {

	var prefix []byte
	switch dbPrefix {
	case "day":
		prefix = dbDayPrefix
	case "week":
		prefix = dbWeekPrefix
	case "app":
		prefix = dbAppPrefix
	case "tasks":
		prefix = dbTaskKey()
	default:
		return fmt.Errorf("no bucket of such key prefix - %s in the db", dbPrefix)
	}

	err := bs.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false //key-only iterations, several other magnitudes faster the docs says
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			key := item.Key()
			if bytes.HasPrefix(key, prefix) {
				if err := bs.DeleteKey(key); err != nil {
					fmt.Println("error deleting key", string(key))
					continue
				}
				fmt.Println("successfully deleted key:", string(key))
			}
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (bs *BadgerDBStore) UpdateOpertionOnBuCKET(dbPrefix string, opsFunc func([]byte) ([]byte, error)) error {

	var prefix []byte
	switch dbPrefix {
	case "day":
		prefix = dbDayPrefix
	case "week":
		prefix = dbWeekPrefix
	case "app":
		prefix = dbAppPrefix
	default:
		return fmt.Errorf("no bucket of such key prefix - %s in the db", dbPrefix)
	}

	err := bs.db.Update(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = true
		opts.PrefetchSize = 100
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			err := it.Item().Value(func(val []byte) error {
				updatedByteArr, err := opsFunc(val)
				if err != nil {
					return err
				}
				return txn.Set(it.Item().Key(), updatedByteArr)
			})

			if err != nil {
				return err
			}
		}
		return nil // all Update successful
	})

	if err != nil {
		return err
	}
	return nil
}

func ExampleOf_opsFunc(v []byte) ([]byte, error) {
	var (
		app AppInfo
		err error
	)

	if app, err = helperFuncs.DecodeJSON[AppInfo](v); err != nil {
		return nil, err
	}
	fmt.Println(app.AppName, app.IsCategorySet, app.DesktopCategories)
	return helperFuncs.EncodeJSON(app)
}
