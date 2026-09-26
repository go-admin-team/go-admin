package tools

import (
	"errors"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// A named, shared in-memory database: with a bare ":memory:" every pooled
// connection opens its own empty database, and a row written inside a
// transaction would be invisible to the query that checks for it.
func openAtomicDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(new(SysTables), new(SysColumns)); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

// rejectNthColumnWrite makes the n-th write of a sys_columns row fail, the way
// it does when the database refuses it.
func rejectNthColumnWrite(register func(name string, fn func(*gorm.DB)) error, n int) error {
	seen := 0
	return register("test:reject_column", func(tx *gorm.DB) {
		if tx.Statement.Table != "sys_columns" {
			return
		}
		seen++
		if seen == n {
			_ = tx.AddError(errors.New("column rejected"))
		}
	})
}

func threeColumns() []SysColumns {
	return []SysColumns{{ColumnName: "id"}, {ColumnName: "name"}, {ColumnName: "note"}}
}

func TestSysTablesCreateWritesAllColumnsOrNone(t *testing.T) {
	db := openAtomicDB(t)
	if err := rejectNthColumnWrite(db.Callback().Create().Before("gorm:create").Register, 2); err != nil {
		t.Fatal(err)
	}

	table := SysTables{TBName: "orders", Columns: threeColumns()}
	if _, err := table.Create(db); err == nil {
		t.Fatal("Create reported success although a column was rejected")
	}

	var tables, columns int64
	db.Model(new(SysTables)).Where("table_name = ?", "orders").Count(&tables)
	db.Model(new(SysColumns)).Count(&columns)
	if tables != 0 || columns != 0 {
		t.Errorf("left %d table row(s) and %d column row(s) behind; want none", tables, columns)
	}
}

func TestSysTablesUpdateWritesAllColumnsOrNone(t *testing.T) {
	db := openAtomicDB(t)
	table := SysTables{TBName: "orders", TableComment: "before", Columns: threeColumns()}
	created, err := table.Create(db)
	if err != nil {
		t.Fatal(err)
	}
	if err := rejectNthColumnWrite(db.Callback().Update().Before("gorm:update").Register, 2); err != nil {
		t.Fatal(err)
	}

	edited, err := (&SysTables{TableId: created.TableId}).Get(db, false)
	if err != nil {
		t.Fatal(err)
	}
	edited.TableComment = "after"
	for i := range edited.Columns {
		edited.Columns[i].ColumnComment = "after"
	}
	if _, err := edited.Update(db); err == nil {
		t.Fatal("Update reported success although a column was rejected")
	}

	stored, err := (&SysTables{TableId: created.TableId}).Get(db, false)
	if err != nil {
		t.Fatal(err)
	}
	if stored.TableComment != "before" {
		t.Errorf("table comment = %q; the table row was kept although a column failed", stored.TableComment)
	}
	for _, c := range stored.Columns {
		if c.ColumnComment != "" {
			t.Errorf("column %s comment = %q; a column write was kept although another failed", c.ColumnName, c.ColumnComment)
		}
	}
}
