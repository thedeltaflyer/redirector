package models

import (
	"os"
	"testing"

	"redirector/helpers"

	bolt "go.etcd.io/bbolt"
)

func TestReplace(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	bucket := []byte("testBucket")
	setupBucket(t, db, bucket)

	kv := &KVWrapper{
		DB:     db,
		Bucket: bucket,
	}

	tests := []struct {
		name        string
		setup       func()
		key         []byte
		newValue    []byte
		expectedOld []byte
		expectedErr error
	}{
		{
			name: "valid replacement",
			setup: func() {
				err := kv.Put([]byte("key1"), []byte("value1"))
				if err != nil {
					t.Fatalf("setup failed: %v", err)
				}
			},
			key:         []byte("key1"),
			newValue:    []byte("newValue1"),
			expectedOld: []byte("value1"),
			expectedErr: nil,
		},
		{
			name: "non-existent key",
			setup: func() {
				// No pre-existing key
			},
			key:         []byte("missingKey"),
			newValue:    []byte("newValue2"),
			expectedOld: nil,
			expectedErr: helpers.NewDoesNotExistError([]byte("missingKey")),
		},
		{
			name: "empty key",
			setup: func() {
				err := kv.Put([]byte("key3"), []byte("value3"))
				if err != nil {
					t.Fatalf("setup failed: %v", err)
				}
			},
			key:         []byte(""),
			newValue:    []byte("newValue3"),
			expectedOld: nil,
			expectedErr: helpers.NewDoesNotExistError([]byte("")),
		},
		{
			name: "empty value",
			setup: func() {
				err := kv.Put([]byte("key4"), []byte("value4"))
				if err != nil {
					t.Fatalf("setup failed: %v", err)
				}
			},
			key:         []byte("key4"),
			newValue:    []byte(""),
			expectedOld: []byte("value4"),
			expectedErr: nil,
		},
		{
			name: "replacement in empty bucket",
			setup: func() {
				kv.Bucket = []byte("emptyBucket")
				setupBucket(t, db, kv.Bucket)
			},
			key:         []byte("key5"),
			newValue:    []byte("value5"),
			expectedOld: nil,
			expectedErr: helpers.NewDoesNotExistError([]byte("key5")),
		},
		{
			name: "replacement with same value",
			setup: func() {
				err := kv.Put([]byte("key6"), []byte("value6"))
				if err != nil {
					t.Fatalf("setup failed: %v", err)
				}
			},
			key:         []byte("key6"),
			newValue:    []byte("value6"),
			expectedOld: []byte("value6"),
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Ensure setup is performed before the test
			tt.setup()

			// Reset to original bucket after test if changed
			defer func() {
				kv.Bucket = bucket
			}()

			oldValue, err := kv.Replace(tt.key, tt.newValue)

			// Compare error
			if (err != nil && tt.expectedErr == nil) || (err == nil && tt.expectedErr != nil) ||
				(err != nil && tt.expectedErr != nil && err.Error() != tt.expectedErr.Error()) {
				t.Fatalf("expected error %v, got %v", tt.expectedErr, err)
			}

			// Compare old value
			if string(oldValue) != string(tt.expectedOld) {
				t.Errorf("expected old value %q, got %q", tt.expectedOld, oldValue)
			}

			// If no error, verify the new value is stored
			if err == nil {
				storedValue, err := kv.Get(tt.key)
				if err != nil {
					t.Fatalf("failed to get key: %v", err)
				}

				if string(storedValue) != string(tt.newValue) {
					t.Errorf("expected new value %q, got %q", tt.newValue, storedValue)
				}
			}
		})
	}
}

func TestExclusivePut(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	bucket := []byte("testBucket")
	setupBucket(t, db, bucket)

	kv := &KVWrapper{
		DB:     db,
		Bucket: bucket,
	}

	tests := []struct {
		name        string
		setup       func()
		key         []byte
		value       []byte
		expectedErr error
	}{
		{
			name: "key does not exist, valid input",
			setup: func() {
				// No pre-existing key
			},
			key:         []byte("key1"),
			value:       []byte("value1"),
			expectedErr: nil,
		},
		{
			name: "key exists already",
			setup: func() {
				err := kv.Put([]byte("key2"), []byte("value2"))
				if err != nil {
					t.Fatalf("failed to setup test: %v", err)
				}
			},
			key:         []byte("key2"),
			value:       []byte("newValue"),
			expectedErr: helpers.NewAlreadyExistsError([]byte("key2")), // Custom error check
		},
		{
			name: "empty key",
			setup: func() {
				// No specific setup needed
			},
			key:         []byte(""),
			value:       []byte("value3"),
			expectedErr: bolt.ErrKeyRequired,
		},
		{
			name: "empty value with a valid key",
			setup: func() {
				// No specific setup needed
			},
			key:         []byte("key3"),
			value:       []byte(""),
			expectedErr: nil,
		},
		{
			name: "empty bucket",
			setup: func() {
				kv.Bucket = []byte("emptyBucket")
				setupBucket(t, db, kv.Bucket)
			},
			key:         []byte("key4"),
			value:       []byte("value4"),
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Ensure setup is performed before the test
			tt.setup()

			// Reset to original bucket after test if changed
			defer func() {
				kv.Bucket = bucket
			}()

			err := kv.ExclusivePut(tt.key, tt.value)

			if (err != nil && tt.expectedErr == nil) || (err == nil && tt.expectedErr != nil) ||
				(err != nil && tt.expectedErr != nil && err.Error() != tt.expectedErr.Error()) {
				t.Fatalf("expected error %v, got %v", tt.expectedErr, err)
			}

			if err == nil {
				storedValue, err := kv.Get(tt.key)
				if err != nil {
					t.Fatalf("failed to get key: %v", err)
				}

				// Ensure the stored value matches the expected value
				if string(storedValue) != string(tt.value) {
					t.Errorf("expected value %q, got %q", tt.value, storedValue)
				}
			}
		})
	}
}

func TestPut(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	bucket := []byte("testBucket")
	setupBucket(t, db, bucket)

	kv := &KVWrapper{
		DB:     db,
		Bucket: bucket,
	}

	tests := []struct {
		name        string
		key         []byte
		value       []byte
		expectedErr error
	}{
		{
			name:        "valid key and value",
			key:         []byte("key1"),
			value:       []byte("value1"),
			expectedErr: nil,
		},
		{
			name:        "empty key",
			key:         []byte(""),
			value:       []byte("value2"),
			expectedErr: bolt.ErrKeyRequired,
		},
		{
			name:        "empty value",
			key:         []byte("key3"),
			value:       []byte(""),
			expectedErr: nil,
		},
		{
			name:        "empty key and value",
			key:         []byte(""),
			value:       []byte(""),
			expectedErr: bolt.ErrKeyRequired,
		},
		{
			name:        "large key and value",
			key:         []byte("key-longer-than-usual"),
			value:       []byte("value-longer-than-usual-value-longer-than-usual"),
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := kv.Put(tt.key, tt.value)
			if err != tt.expectedErr {
				t.Fatalf("expected error %v, got %v", tt.expectedErr, err)
			}

			if tt.expectedErr == nil {
				got, err := kv.Get(tt.key)
				if err != nil {
					t.Fatalf("failed to get key %q: %v", tt.key, err)
				}

				if string(got) != string(tt.value) {
					t.Errorf("expected value %q, got %q", tt.value, got)
				}
			}
		})
	}
}

func setupTestDB(t *testing.T) (*bolt.DB, func()) {
	t.Helper()

	db, err := bolt.Open("test.db", 0600, nil)
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	cleanup := func() {
		db.Close()
		os.Remove("test.db")
	}

	return db, cleanup
}

func setupBucket(t *testing.T, db *bolt.DB, bucket []byte) {
	t.Helper()
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucket)
		return err
	})
	if err != nil {
		t.Fatalf("failed to create bucket: %v", err)
	}
}

func TestGet(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	bucket := []byte("testBucket")
	setupBucket(t, db, bucket)

	kv := &KVWrapper{
		DB:     db,
		Bucket: bucket,
	}

	tests := []struct {
		name        string
		setup       func()
		key         []byte
		expected    []byte
		expectedErr error
	}{
		{
			name: "key exists",
			setup: func() {
				err := kv.Put([]byte("key1"), []byte("value1"))
				if err != nil {
					t.Fatalf("failed to setup test: %v", err)
				}
			},
			key:         []byte("key1"),
			expected:    []byte("value1"),
			expectedErr: nil,
		},
		{
			name: "key does not exist",
			setup: func() {
				err := kv.Put([]byte("key2"), []byte("value2"))
				if err != nil {
					t.Fatalf("failed to setup test: %v", err)
				}
			},
			key:         []byte("missingKey"),
			expected:    nil,
			expectedErr: nil,
		},
		{
			name: "empty key",
			setup: func() {
				err := kv.Put([]byte("key3"), []byte("value3"))
				if err != nil {
					t.Fatalf("failed to setup test: %v", err)
				}
			},
			key:         []byte(""),
			expected:    nil,
			expectedErr: nil,
		},
		{
			name: "empty bucket",
			setup: func() {
				kv.Bucket = []byte("emptyBucket")
				setupBucket(t, db, kv.Bucket)
			},
			key:         []byte("key4"),
			expected:    nil,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Ensure setup is performed before the test
			tt.setup()

			// Reset to original bucket after test if changed
			defer func() {
				kv.Bucket = bucket
			}()

			value, err := kv.Get(tt.key)

			// Compare error and expected value
			if err != tt.expectedErr {
				t.Fatalf("expected error %v, got %v", tt.expectedErr, err)
			}

			if string(value) != string(tt.expected) {
				t.Errorf("expected value %s, got %s", tt.expected, value)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	bucket := []byte("testBucket")
	setupBucket(t, db, bucket)

	kv := &KVWrapper{
		DB:     db,
		Bucket: bucket,
	}

	tests := []struct {
		name        string
		setup       func()
		key         []byte
		expectedErr error
	}{
		{
			name: "delete existing key",
			setup: func() {
				err := kv.Put([]byte("key1"), []byte("value1"))
				if err != nil {
					t.Fatalf("setup failed: %v", err)
				}
			},
			key:         []byte("key1"),
			expectedErr: nil,
		},
		{
			name: "delete non-existent key",
			setup: func() {
				// No setup required
			},
			key:         []byte("missingKey"),
			expectedErr: helpers.NewDoesNotExistError([]byte("missingKey")),
		},
		{
			name: "delete with empty key",
			setup: func() {
				err := kv.Put([]byte("key2"), []byte("value2"))
				if err != nil {
					t.Fatalf("setup failed: %v", err)
				}
			},
			key:         []byte(""),
			expectedErr: helpers.NewDoesNotExistError([]byte("")),
		},
		{
			name: "delete from empty bucket",
			setup: func() {
				kv.Bucket = []byte("emptyBucket")
				setupBucket(t, db, kv.Bucket)
			},
			key:         []byte("key3"),
			expectedErr: helpers.NewDoesNotExistError([]byte("key3")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Ensure setup is performed before the test
			tt.setup()

			// Reset to the original bucket after the test if changed
			defer func() {
				kv.Bucket = bucket
			}()

			err := kv.Delete(tt.key)

			// Compare error
			if (err != nil && tt.expectedErr == nil) || (err == nil && tt.expectedErr != nil) ||
				(err != nil && tt.expectedErr != nil && err.Error() != tt.expectedErr.Error()) {
				t.Fatalf("expected error %v, got %v", tt.expectedErr, err)
			}

			// Verify the key has been deleted
			if tt.expectedErr == nil {
				value, err := kv.Get(tt.key)
				if err != nil {
					t.Fatalf("unexpected error retrieving value: %v", err)
				}

				if value != nil {
					t.Errorf("expected key %q to be deleted, but got value %q", tt.key, value)
				}
			}
		})
	}
}
