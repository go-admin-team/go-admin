package service

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"go-admin/app/admin/models"
)

// The admin branch of getSysMenuByRoleName carried "deleted_at is null" in its
// where clause. GORM adds that condition itself for a model with a DeletedAt
// field, so the clause was a duplicate — and one written in terms of a column
// being null, which stops being true the moment the column stops being
// nullable. This pins the behaviour the clause was there for.
func TestSoftDeletedMenusAreNotReturned(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	if err := db.AutoMigrate(&models.SysMenu{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	live := models.SysMenu{MenuName: "live", MenuType: "M"}
	gone := models.SysMenu{MenuName: "gone", MenuType: "M"}
	if err := db.Create(&live).Error; err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := db.Create(&gone).Error; err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := db.Delete(&gone).Error; err != nil {
		t.Fatalf("delete: %v", err)
	}

	// Through getByRoleName rather than a copy of its query: a test that
	// reissues the statement passes whether or not the production line still
	// says what it is supposed to, which is what the first version of this
	// test did.
	e := &SysMenu{}
	e.Orm = db
	got, err := e.getByRoleName("admin")
	if err != nil {
		t.Fatalf("getByRoleName: %v", err)
	}

	if len(got) != 1 {
		t.Fatalf("got %d rows, want 1", len(got))
	}
	if got[0].MenuName != "live" {
		t.Errorf("got %q, want the row that was not deleted", got[0].MenuName)
	}
}
