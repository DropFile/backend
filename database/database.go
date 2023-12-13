package database

import (
	"strings"
	"sync"

	badger "github.com/dgraph-io/badger/v4"
)

type KVStore struct {
	db *badger.DB
	mu sync.RWMutex
}

// to connect with the database
func NewKVStore(dbPath string) (*KVStore, error) {
	opts := badger.DefaultOptions(dbPath)
	db, error := badger.Open(opts)

	if error != nil {
		return nil, error
	}

	return &KVStore{db: db}, nil
}

// close the connection with the database
func (kv *KVStore) Close() {
	kv.db.Close()
}

// set a list of filenames in database
func (kv *KVStore) Set(key string, value []string) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	return kv.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), []byte(strings.Join(value, ",")))
	})
}

// get a list of filenames in database
func (kv *KVStore) Get(key string) ([]string, error) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	var result []string
	err := kv.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		value, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}
		result = strings.Split(string(value), ",")
		return nil
	})
	return result, err
}