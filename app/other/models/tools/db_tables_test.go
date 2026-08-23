package tools

import (
	"testing"

	"github.com/glebarez/sqlite"
	config2 "github.com/go-admin-team/go-admin-core/v2/sdk/config"
	"gorm.io/gorm"
)

const generatorSchema = "go_admin_test"

// newCandidateDB stands in for MySQL: sqlite is given an attached database
// called information_schema so the same query runs, and sys_tables lives on the
// connection the way it does in production.
func newCandidateDB(t *testing.T, tables ...string) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.Exec(`ATTACH DATABASE ':memory:' AS information_schema`).Error; err != nil {
		t.Fatalf("attach information_schema: %v", err)
	}
	if err := db.Exec("CREATE TABLE information_schema.`tables` (" +
		"TABLE_NAME text, TABLE_SCHEMA text, `ENGINE` text, TABLE_ROWS text," +
		"TABLE_COLLATION text, CREATE_TIME text, UPDATE_TIME text, TABLE_COMMENT text)").Error; err != nil {
		t.Fatalf("create information_schema.tables: %v", err)
	}
	for _, name := range tables {
		if err := db.Exec("INSERT INTO information_schema.`tables` (TABLE_NAME, TABLE_SCHEMA) VALUES (?, ?)",
			name, generatorSchema).Error; err != nil {
			t.Fatalf("insert %s: %v", name, err)
		}
	}
	if err := db.AutoMigrate(new(SysTables)); err != nil {
		t.Fatalf("migrate sys_tables: %v", err)
	}

	previousDriver := config2.DatabaseConfig.Driver
	previousName := config2.GenConfig.DBName
	config2.DatabaseConfig.Driver = "mysql"
	config2.GenConfig.DBName = generatorSchema
	t.Cleanup(func() {
		config2.DatabaseConfig.Driver = previousDriver
		config2.GenConfig.DBName = previousName
	})
	return db
}

func candidateNames(t *testing.T, db *gorm.DB) []string {
	t.Helper()
	found, _, err := new(DBTables).GetPage(db, 100, 1)
	if err != nil {
		t.Fatalf("GetPage: %v", err)
	}
	names := make([]string, 0, len(found))
	for _, row := range found {
		names = append(names, row.TableName)
	}
	return names
}

func contains(names []string, want string) bool {
	for _, name := range names {
		if name == want {
			return true
		}
	}
	return false
}

// An empty sys_tables must not filter everything out - that is what a bare
// NOT IN (empty set) does, and a fresh install is exactly the case where the
// list matters most.
func TestGetPageListsEveryTableWhenNoneAreRegistered(t *testing.T) {
	db := newCandidateDB(t, "sys_user", "sys_role")

	names := candidateNames(t, db)
	if len(names) != 2 {
		t.Fatalf("want both tables offered on a fresh install, got %v", names)
	}
}

func TestGetPageSkipsAlreadyRegisteredTables(t *testing.T) {
	db := newCandidateDB(t, "sys_user", "sys_role")
	if err := db.Create(&SysTables{TBName: "sys_user"}).Error; err != nil {
		t.Fatalf("register sys_user: %v", err)
	}

	names := candidateNames(t, db)
	if contains(names, "sys_user") {
		t.Errorf("sys_user is already registered and was offered again: %v", names)
	}
	if !contains(names, "sys_role") {
		t.Errorf("sys_role is not registered and was withheld: %v", names)
	}
}

// Deleting the generator entry has to hand the table back, which the raw
// subquery never did: it read the row whether or not it was soft-deleted.
func TestGetPageOffersTablesWhoseEntryWasDeleted(t *testing.T) {
	db := newCandidateDB(t, "sys_user")
	registered := SysTables{TBName: "sys_user"}
	if err := db.Create(&registered).Error; err != nil {
		t.Fatalf("register sys_user: %v", err)
	}
	if err := db.Delete(&registered).Error; err != nil {
		t.Fatalf("delete the entry: %v", err)
	}

	if names := candidateNames(t, db); !contains(names, "sys_user") {
		t.Errorf("the generator entry is deleted, sys_user should be a candidate again: %v", names)
	}
}
