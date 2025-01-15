package database

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	bolt "go.etcd.io/bbolt"
)

var testDBPath = filepath.Join(os.TempDir(), "test.db")

func setupTestDB(t *testing.T) *bolt.DB {
	t.Helper()
	db, err := bolt.Open(testDBPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}
	return db
}

func cleanupTestDB(t *testing.T) {
	t.Helper()
	if err := os.Remove(testDBPath); err != nil && !os.IsNotExist(err) {
		t.Fatalf("failed to clean up test database: %v", err)
	}
}

func TestInitDB(t *testing.T) {
	defer cleanupTestDB(t)

	t.Run("valid initialization without migration", func(t *testing.T) {
		InitDB(testDBPath, false)
		defer CloseDB()

		if db == nil {
			t.Fatal("database was not initialized")
		}
	})

	t.Run("valid initialization with migration", func(t *testing.T) {
		InitDB(testDBPath, true)
		defer CloseDB()

		if db == nil {
			t.Fatal("database was not initialized")
		}

		// Check for buckets created during migration
		buckets := []string{"redirects", "api_keys", "health_checks"}
		err := db.View(func(tx *bolt.Tx) error {
			for _, bucket := range buckets {
				if tx.Bucket([]byte(bucket)) == nil {
					t.Fatalf("bucket %q not found", bucket)
				}
			}
			return nil
		})
		if err != nil {
			t.Fatalf("unexpected view error: %v", err)
		}
	})
}

func TestMigrateDB(t *testing.T) {
	defer cleanupTestDB(t)

	db = setupTestDB(t)
	defer CloseDB()

	MigrateDB()

	// Verify the buckets are created
	buckets := []string{"redirects", "api_keys", "health_checks"}
	err := db.View(func(tx *bolt.Tx) error {
		for _, bucket := range buckets {
			if tx.Bucket([]byte(bucket)) == nil {
				t.Fatalf("bucket %q not found", bucket)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected view error: %v", err)
	}
}

func TestInitHealthCheckData(t *testing.T) {
	defer cleanupTestDB(t)

	db = setupTestDB(t)
	defer CloseDB()

	// Create health_checks bucket manually
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("health_checks"))
		return err
	})
	if err != nil {
		t.Fatalf("failed to create health_checks bucket: %v", err)
	}

	InitHealthCheckData()

	// Verify data in the health_checks bucket
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("health_checks"))
		if b == nil {
			t.Fatal("health_checks bucket not found")
		}
		if v := b.Get([]byte("health")); string(v) != "ok" {
			t.Fatalf("unexpected value for health key: got %s, want %s", v, "ok")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected view error: %v", err)
	}
}

func TestGetDB(t *testing.T) {
	defer cleanupTestDB(t)

	t.Run("uninitialized database", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic for uninitialized database")
			}
		}()
		_ = GetDB()
	})

	t.Run("initialized database", func(t *testing.T) {
		InitDB(testDBPath, false)
		defer CloseDB()

		dbInstance := GetDB()
		if dbInstance != db {
			t.Fatal("GetDB did not return the correct database instance")
		}
	})
}

func TestCloseDB(t *testing.T) {
	defer cleanupTestDB(t)

	t.Run("close uninitialized database", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic when closing uninitialized database")
			}
		}()
		CloseDB()
	})

	t.Run("close initialized database", func(t *testing.T) {
		InitDB(testDBPath, false)
		CloseDB()

		if db != nil {
			t.Fatal("database was not set to nil after closing")
		}
	})
}

func TestCheckOrCreateBucket(t *testing.T) {
	defer cleanupTestDB(t)

	db = setupTestDB(t)
	defer CloseDB()

	err := db.Update(func(tx *bolt.Tx) error {
		if err := checkOrCreateBucket(tx, []byte("test_bucket")); err != nil {
			t.Fatalf("unexpected error creating bucket: %v", err)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = db.View(func(tx *bolt.Tx) error {
		if tx.Bucket([]byte("test_bucket")) == nil {
			t.Fatal("test_bucket was not created")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCheckOrCreateHealthCheckKey(t *testing.T) {
	defer cleanupTestDB(t)

	db = setupTestDB(t)
	defer CloseDB()

	t.Run("health_checks bucket does not exist yet", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic for uninitialized database")
			}
		}()
		err := db.Update(func(tx *bolt.Tx) error {
			checkOrCreateHealthCheckKey(tx)
			return nil
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("health_checks"))
		return err
	})
	if err != nil {
		t.Fatalf("failed to create health_checks bucket: %v", err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		checkOrCreateHealthCheckKey(tx)
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	//err = db.View(func(tx *bolt.Tx) error {
	//	b := tx.Bucket([]byte("health_checks"))
	//	if b == nil {
	//		t.Fatal("health_checks bucket not found")
	//	}
	//	if v := b.Get([]byte("health")); string(v) != "ok" {
	//		t.Fatalf("unexpected value for health key: got %s, want %s", v, "ok")
	//	}
	//	return nil
	//})
	//if err != nil {
	//	t.Fatalf("unexpected error: %v", err)
	//}
}
