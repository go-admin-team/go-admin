//go:build sqlite3

package database

import (
	"path/filepath"
	"sync"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type lockRow struct {
	ID int `gorm:"primaryKey"`
	N  int
}

// hammer runs writers and readers against one file at the same time and
// returns how many statements failed.
func hammer(t *testing.T, dsn string) (failed int, firstErr error) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&lockRow{}); err != nil {
		t.Fatal(err)
	}
	db.Create(&lockRow{ID: 1})

	var mu sync.Mutex
	note := func(err error) {
		if err == nil {
			return
		}
		mu.Lock()
		failed++
		if firstErr == nil {
			firstErr = err
		}
		mu.Unlock()
	}

	var wg sync.WaitGroup
	for w := 0; w < 8; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 40; i++ {
				note(db.Create(&lockRow{N: i}).Error)
				// read-then-write transaction: the shape a deferred
				// transaction cannot upgrade without SQLITE_BUSY
				note(db.Transaction(func(tx *gorm.DB) error {
					var r lockRow
					if err := tx.First(&r, 1).Error; err != nil {
						return err
					}
					return tx.Model(&lockRow{}).Where("id = 1").Update("n", r.N+1).Error
				}))
			}
		}()
	}
	for r := 0; r < 4; r++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 80; i++ {
				var c int64
				note(db.Model(&lockRow{}).Count(&c).Error)
			}
		}()
	}
	wg.Wait()
	return failed, firstErr
}

func TestSqliteDefaultsPreventLocking(t *testing.T) {
	file := filepath.Join(t.TempDir(), "lock.db")
	if failed, err := hammer(t, sqliteDSN(file)); failed != 0 {
		t.Fatalf("%d statements failed with the defaults applied, first: %v", failed, err)
	}
}

// Documents what the defaults are for: the bare path loses statements.
func TestSqliteBarePathLocks(t *testing.T) {
	file := filepath.Join(t.TempDir(), "bare.db")
	if failed, _ := hammer(t, file); failed == 0 {
		t.Skip("bare path did not lock on this machine; the comparison proves nothing here")
	}
}
