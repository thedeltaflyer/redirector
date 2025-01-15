package models

import (
	"redirector/helpers"

	bolt "go.etcd.io/bbolt"
)

// KV defines an interface for key-value storage operations with methods for retrieval, insertion, replacement, and deletion.
type KV interface {
	Get(key []byte) ([]byte, error)
	Put(key []byte, value []byte) error
	ExclusivePut(key []byte, value []byte) error
	Replace(key []byte, value []byte) ([]byte, error)
	Delete(key []byte) error
}

// KVWrapper provides a wrapper around a BoltDB instance and a specific bucket for key-value operations using the KV interface.
type KVWrapper struct {
	DB     *bolt.DB
	Bucket []byte
}

// Get retrieves the value associated with the provided key from the underlying BoltDB bucket.
// Returns the value along with any error encountered during the retrieval process.
func (kv *KVWrapper) Get(key []byte) ([]byte, error) {
	var value []byte
	err := kv.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(kv.Bucket)
		value = b.Get(key)
		return nil
	})
	return value, err
}

// Put inserts or updates the specified key-value pair in the BoltDB bucket. Returns an error if the operation fails.
func (kv *KVWrapper) Put(key []byte, value []byte) error {
	return kv.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(kv.Bucket)
		err := b.Put(key, value)
		return err
	})
}

// ExclusivePut attempts to insert a key-value pair into the bucket only if the key does not already exist.
// Returns an AlreadyExistsError if the key is already present in the bucket.
// Returns an error if the operation fails.
func (kv *KVWrapper) ExclusivePut(key []byte, value []byte) error {
	return kv.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(kv.Bucket)
		testVal := b.Get(key)
		if testVal != nil {
			return helpers.NewAlreadyExistsError(key)
		}
		err := b.Put(key, value)
		return err
	})
}

// Replace updates the value for a given key and returns the old value. Returns an error if the key does not exist.
func (kv *KVWrapper) Replace(key []byte, value []byte) ([]byte, error) {
	var oldVal []byte
	err := kv.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(kv.Bucket)
		oldVal = b.Get(key)
		if oldVal == nil {
			return helpers.NewDoesNotExistError(key)
		}
		err := b.Put(key, value)
		return err
	})
	return oldVal, err
}

// Delete removes the specified key from the BoltDB bucket. Returns a DoesNotExistError if the key does not exist.
func (kv *KVWrapper) Delete(key []byte) error {
	return kv.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(kv.Bucket)
		oldVal := b.Get(key)
		if oldVal == nil {
			return helpers.NewDoesNotExistError(key)
		}
		return b.Delete(key)
	})
}
