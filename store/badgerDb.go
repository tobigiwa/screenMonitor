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

func (bs *BadgerDBStore) WriteUsuage(data ScreenTime) error {
	item, err := bs.get(data.AppName)
	if err != nil {
		// TODO
		if errors.Is(err, badger.ErrKeyNotFound) {
			fmt.Println("new app")

			var appInfo appInfo

			appInfo.AppName = data.AppName
			appInfo.ActiveScreenStats = make(dailyAppScreenTime)
			appInfo.PasiveScreenStats = make(dailyAppScreenTime)

			if data.Type == Active {
				appInfo.ActiveScreenStats[Key()] = data.Time
			} else {
				appInfo.PasiveScreenStats[Key()] = data.Time
			}

			ser, err := appInfo.serialize(appInfo)
			if err != nil {
				return err
			}

			if err := bs.set([]byte(data.AppName), ser); err != nil {
				return err
			}
			return nil
		}

		return err
	}

	var appInfo appInfo

	if err := appInfo.derialize(item); err != nil {
		return err
	}

	if data.AppName != appInfo.AppName {
		return ErrAppKeyMismatch
	}

	if data.Type == Active {
		appInfo.ActiveScreenStats[Key()] += data.Time
	} else {
		appInfo.PasiveScreenStats[Key()] += data.Time
	}

	return nil
}

// type KV struct {
// 	db *badger.DB
// }

// func (k *KV) Exists(key string) (bool, error) {
// 	var exists bool
// 	err := k.db.View(
// 		func(tx *badger.Txn) error {
// 			if val, err := tx.Get([]byte(key)); err != nil {
// 				return err
// 			} else if val != nil {
// 				exists = true
// 			}
// 			return nil
// 		})
// 	if errors.Is(err, badger.ErrKeyNotFound) {
// 		err = nil
// 	}
// 	return exists, err
// }
// func (k *KV) Get(key string) (string, error) {
// 	var value string
// 	return value, k.db.View(
// 		func(tx *badger.Txn) error {
// 			item, err := tx.Get([]byte(key))
// 			if err != nil {
// 				return fmt.Errorf("getting value: %w", err)
// 			}
// 			valCopy, err := item.ValueCopy(nil)
// 			if err != nil {
// 				return fmt.Errorf("copying value: %w", err)
// 			}
// 			value = string(valCopy)
// 			return nil
// 		})
// }

// func (k *KV) Set(key, value string) error {
// 	return k.db.Update(
// 		func(txn *badger.Txn) error {
// 			return txn.Set([]byte(key), []byte(value))
// 		})
// }

// func (k *KV) Delete(key string) error {
// 	return k.db.Update(
// 		func(txn *badger.Txn) error {
// 			return txn.Delete([]byte(key))
// 		})
// }
// func d() {
// 	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer db.Close()

// 	err = db.View(func(txn *badger.Txn) error {
// 		item, err := txn.Get([]byte("answer"))
// 		handle(err)

// 		var valNot, valCopy []byte
// 		// err = item.Value(func(val []byte) error {
// 		// 	// This func with val would only be called if item.Value encounters no error.

// 		// 	// Accessing val here is valid.
// 		// 	fmt.Printf("The answer is: %s\n", val)

// 		// 	// Copying or parsing val is valid.
// 		// 	valCopy = append([]byte{}, val...)

// 		// 	// Assigning val slice to another variable is NOT OK.
// 		// 	valNot = val // Do not do this.
// 		// 	return nil
// 		// })
// 		//   handle(err)

// 		// DO NOT access val here. It is the most common cause of bugs.
// 		fmt.Printf("NEVER do this. %s\n", valNot)

// 		// You must copy it to use it outside item.Value(...).
// 		fmt.Printf("The answer is: %s\n", valCopy)

// 		// Alternatively, you could also use item.ValueCopy().
// 		valCopy, err = item.ValueCopy(nil)
// 		handle(err)
// 		fmt.Printf("The answer is: %s\n", valCopy)

// 		return nil
// 	})

// }

// func handle(err error) {}
