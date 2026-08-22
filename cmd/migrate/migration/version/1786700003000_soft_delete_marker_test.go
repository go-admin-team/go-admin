package version

import (
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// oldUser is the shape of the table before this migration: deleted_at is a
// nullable timestamp, and nothing constrains the username.
type oldUser struct {
	Id        int64 `gorm:"primaryKey;autoIncrement"`
	Username  string
	DeletedAt *time.Time
}

func (oldUser) TableName() string { return "sys_user" }

func openWithOldSchema(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	if err := db.AutoMigrate(&oldUser{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

func TestConvertsDeletedAtAndKeepsWhoWasDeleted(t *testing.T) {
	db := openWithOldSchema(t)

	gone := time.Now().Add(-time.Hour)
	rows := []oldUser{
		{Username: "alice"},
		{Username: "bob", DeletedAt: &gone},
	}
	for i := range rows {
		if err := db.Create(&rows[i]).Error; err != nil {
			t.Fatalf("seed: %v", err)
		}
	}

	if err := convertDeletedAt(db, "sys_user"); err != nil {
		t.Fatalf("convert: %v", err)
	}

	var marks []struct {
		Username  string
		DeletedAt int64
	}
	if err := db.Raw("SELECT username, deleted_at FROM sys_user ORDER BY id").Scan(&marks).Error; err != nil {
		t.Fatalf("read back: %v", err)
	}
	if len(marks) != 2 {
		t.Fatalf("got %d rows, want 2", len(marks))
	}
	if marks[0].DeletedAt != 0 {
		t.Errorf("a live row carries %d, want 0", marks[0].DeletedAt)
	}
	if want := gone.UnixMilli(); marks[1].DeletedAt != want {
		t.Errorf("the deleted row carries %d, want %d: the timestamp was lost", marks[1].DeletedAt, want)
	}
}

// Running it twice must be safe: DDL does not roll back on MySQL, so an
// operator whose first attempt failed halfway has nothing to do but run it
// again.
func TestConversionIsRepeatable(t *testing.T) {
	db := openWithOldSchema(t)
	if err := db.Create(&oldUser{Username: "alice"}).Error; err != nil {
		t.Fatalf("seed: %v", err)
	}

	for i := 0; i < 3; i++ {
		if err := convertDeletedAt(db, "sys_user"); err != nil {
			t.Fatalf("convert %d: %v", i, err)
		}
	}
}

// The point of the whole exercise: the constraint has to reject a second live
// row and accept one whose predecessor was deleted.
func TestUniqueIndexBindsLiveRowsOnly(t *testing.T) {
	db := openWithOldSchema(t)
	if err := db.Create(&oldUser{Username: "alice"}).Error; err != nil {
		t.Fatalf("seed: %v", err)
	}
	if err := convertDeletedAt(db, "sys_user"); err != nil {
		t.Fatalf("convert: %v", err)
	}
	if err := createUniqueIndex(db, "sys_user", "username", "uk_sys_user_username"); err != nil {
		t.Fatalf("index: %v", err)
	}

	t.Run("a second live row is rejected", func(t *testing.T) {
		err := db.Exec("INSERT INTO sys_user (username, deleted_at) VALUES ('alice', 0)").Error
		if err == nil {
			t.Fatal("a duplicate username was accepted")
		}
	})

	t.Run("the name is free once the row is deleted", func(t *testing.T) {
		if err := db.Exec("UPDATE sys_user SET deleted_at = ? WHERE username = 'alice'", time.Now().UnixMilli()).Error; err != nil {
			t.Fatalf("delete: %v", err)
		}
		if err := db.Exec("INSERT INTO sys_user (username, deleted_at) VALUES ('alice', 0)").Error; err != nil {
			t.Errorf("the name stayed taken after its row was deleted: %v", err)
		}
	})

	t.Run("the same name can be deleted more than once", func(t *testing.T) {
		if err := db.Exec("UPDATE sys_user SET deleted_at = ? WHERE deleted_at = 0", time.Now().UnixMilli()+1).Error; err != nil {
			t.Errorf("a second deletion collided with the first: %v", err)
		}
	})
}

// A table that already holds duplicates cannot take the index, and the
// migration says which values rather than letting the index fail.
func TestRefusesWhenDuplicatesAlreadyExist(t *testing.T) {
	db := openWithOldSchema(t)
	for i := 0; i < 2; i++ {
		if err := db.Create(&oldUser{Username: "alice"}).Error; err != nil {
			t.Fatalf("seed: %v", err)
		}
	}

	err := refuseOnDuplicates(db)
	if err == nil {
		t.Fatal("the migration accepted a table that already holds duplicates")
	}
	if !contains(err.Error(), "alice") {
		t.Errorf("the error does not name the value: %v", err)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 || indexOf(s, sub) >= 0)
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}

// Why the column has to change at all, demonstrated rather than argued: the
// same index over the column as it was accepts both live rows, because each
// carries NULL and NULL is not equal to NULL. The constraint exists and binds
// nothing — which is worse than its absence, because it reads as protection.
func TestTheIndexOverANullableMarkerEnforcesNothing(t *testing.T) {
	db := openWithOldSchema(t)

	if err := db.Exec("CREATE UNIQUE INDEX uk_nullable ON sys_user (username, deleted_at)").Error; err != nil {
		t.Fatalf("index: %v", err)
	}

	for i := 0; i < 2; i++ {
		if err := db.Create(&oldUser{Username: "alice"}).Error; err != nil {
			t.Fatalf("the nullable marker rejected a duplicate after all, which would make this migration unnecessary: %v", err)
		}
	}

	var live int64
	if err := db.Raw("SELECT COUNT(*) FROM sys_user WHERE deleted_at IS NULL").Scan(&live).Error; err != nil {
		t.Fatalf("count: %v", err)
	}
	if live != 2 {
		t.Fatalf("got %d live rows, want 2", live)
	}
}
