package database

import (
	"fmt"
	"time"

	"redirector/logging"

	bolt "go.etcd.io/bbolt"
)

var (
	db *bolt.DB
)

// InitDB initializes the Bolt database at the given path. If migrate is true, it performs database migrations and setups.
func InitDB(path string, migrate bool) {
	var err error
	db, err = bolt.Open(path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		panic(err)
	}
	if migrate {
		MigrateDB()
		InitHealthCheckData()
	}
}

// MigrateDB initializes required database buckets for storing redirects, API keys, and health checks. It panics on failure.
func MigrateDB() {
	database := GetDB()
	buckets := []string{"redirects", "api_keys", "health_checks"}
	err := database.Update(func(tx *bolt.Tx) error {
		for _, bucket := range buckets {
			createErr := checkOrCreateBucket(tx, []byte(bucket))
			if createErr != nil {
				return createErr
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}

// InitHealthCheckData initializes the health check data in the database, creating the required key if it does not exist.
func InitHealthCheckData() {
	database := GetDB()
	err := database.Update(func(tx *bolt.Tx) error {
		checkOrCreateHealthCheckKey(tx)
		return nil
	})
	if err != nil {
		panic(err)
	}
}

// GetDB returns the initialized Bolt database instance. Panics if the database has not been initialized.
func GetDB() *bolt.DB {
	if db == nil {
		panic(fmt.Errorf("database not initialized"))
	}
	return db
}

// CloseDB gracefully closes the Bolt database if it is initialized and sets the database instance to nil.
// It logs an error if the database is not initialized and panics on any close operation error.
func CloseDB() {
	if db == nil {
		logging.GetLogger().Error("attempting to close database, but it's not initialized")
	}
	err := db.Close()
	if err != nil {
		panic(err)
	}
	db = nil
}

// checkOrCreateBucket ensures a bucket with the given name exists by creating it if it does not already exist.
// It takes a bolt transaction and bucket name as input, returning an error if bucket creation fails.
func checkOrCreateBucket(tx *bolt.Tx, name []byte) error {
	logging.GetLogger().Debugf("Checking/Creating %q bucket", name)
	_, err := tx.CreateBucketIfNotExists(name)
	if err != nil {
		return err
	}
	return nil
}

// checkOrCreateHealthCheckKey ensures the "health" key in the "health_checks" bucket is set to "ok". Panics on errors or missing bucket.
func checkOrCreateHealthCheckKey(tx *bolt.Tx) {
	b := tx.Bucket([]byte("health_checks"))
	if b == nil {
		panic(fmt.Errorf("health_checks bucket not found"))
	}
	if err := b.Put([]byte("health"), []byte("ok")); err != nil {
		panic(err)
	}
}
