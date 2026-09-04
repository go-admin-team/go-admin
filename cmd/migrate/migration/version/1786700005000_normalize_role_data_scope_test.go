package version

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type roleDataScopeRow struct {
	RoleId    int    `gorm:"column:role_id;primaryKey;autoIncrement"`
	DataScope string `gorm:"column:data_scope"`
}

func (roleDataScopeRow) TableName() string { return "sys_role" }

func openRoleTable(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	if err := db.AutoMigrate(&roleDataScopeRow{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

// The migration exists because the shipped admin role is exactly this case:
// config/db.sql's role_id 1 carries an empty data_scope. Reproduces the seed
// data literally rather than a made-up example.
func TestNormalizesTheEmptyDataScopeTheSeedDataShips(t *testing.T) {
	db := openRoleTable(t)
	if err := db.Create(&roleDataScopeRow{RoleId: 1, DataScope: ""}).Error; err != nil {
		t.Fatalf("seed: %v", err)
	}

	if err := normalizeRoleDataScope(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	var row roleDataScopeRow
	if err := db.First(&row, 1).Error; err != nil {
		t.Fatalf("read back: %v", err)
	}
	if row.DataScope != "1" {
		t.Fatalf("data_scope = %q, want %q", row.DataScope, "1")
	}
}

// The five recognized values must survive untouched - this migration
// normalizes what Permission cannot make sense of, not what it already can.
func TestLeavesRecognizedScopesAlone(t *testing.T) {
	db := openRoleTable(t)
	valid := []string{"1", "2", "3", "4", "5"}
	for i, scope := range valid {
		if err := db.Create(&roleDataScopeRow{RoleId: i + 1, DataScope: scope}).Error; err != nil {
			t.Fatalf("seed %d: %v", i, err)
		}
	}

	if err := normalizeRoleDataScope(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	var rows []roleDataScopeRow
	if err := db.Order("role_id").Find(&rows).Error; err != nil {
		t.Fatalf("read back: %v", err)
	}
	for i, row := range rows {
		if row.DataScope != valid[i] {
			t.Errorf("role %d: data_scope = %q, want %q (untouched)", row.RoleId, row.DataScope, valid[i])
		}
	}
}

// A garbage value (not just empty) must be normalized the same way as empty -
// both are "not one of the five", and the migration's WHERE clause has to
// catch both.
func TestNormalizesGarbageScopesToo(t *testing.T) {
	db := openRoleTable(t)
	if err := db.Create(&roleDataScopeRow{RoleId: 1, DataScope: "6"}).Error; err != nil {
		t.Fatalf("seed: %v", err)
	}

	if err := normalizeRoleDataScope(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	var row roleDataScopeRow
	if err := db.First(&row, 1).Error; err != nil {
		t.Fatalf("read back: %v", err)
	}
	if row.DataScope != "1" {
		t.Fatalf("data_scope = %q, want %q", row.DataScope, "1")
	}
}

// Running it twice must be safe: it is a plain UPDATE, not DDL, but
// sys_migration only records success once, and an operator who reruns
// `migrate` on a partially-applied database has to be able to trust that.
func TestNormalizeRoleDataScopeIsRepeatable(t *testing.T) {
	db := openRoleTable(t)
	if err := db.Create(&roleDataScopeRow{RoleId: 1, DataScope: ""}).Error; err != nil {
		t.Fatalf("seed: %v", err)
	}

	for i := 0; i < 3; i++ {
		if err := normalizeRoleDataScope(db); err != nil {
			t.Fatalf("migrate %d: %v", i, err)
		}
	}

	var row roleDataScopeRow
	if err := db.First(&row, 1).Error; err != nil {
		t.Fatalf("read back: %v", err)
	}
	if row.DataScope != "1" {
		t.Fatalf("data_scope = %q, want %q", row.DataScope, "1")
	}
}
