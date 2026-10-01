package models

import (
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/glebarez/sqlite"
	"github.com/go-admin-team/go-admin-core/v2/sdk"
	"github.com/go-admin-team/go-admin-core/v2/storage/queue"
	"gorm.io/gorm"
)

// SaveOperaLog shortens JsonResult before writing it. Cutting at a byte offset
// splits a multi-byte character, and PostgreSQL and strict-mode MySQL reject the
// resulting string, so the row is dropped and the only trace is a log line.
func TestSaveOperaLogTruncatesOnCharacterBoundary(t *testing.T) {
	const tenant = "opera-log-truncate-test"
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.AutoMigrate(&SysOperaLog{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	previous := sdk.Runtime.GetDbByTenant(tenant)
	sdk.Runtime.SetDbByTenant(tenant, db)
	t.Cleanup(func() { sdk.Runtime.SetDbByTenant(tenant, previous) })

	// 150 characters, 3 bytes each: byte 100 falls inside the 34th character.
	long := strings.Repeat("中", 150)

	msg := &queue.Message{}
	msg.SetValues(map[string]interface{}{"title": "t", "jsonResult": long})
	msg.SetPrefix(tenant)
	if err := SaveOperaLog(msg); err != nil {
		t.Fatalf("SaveOperaLog: %v", err)
	}

	var rows []SysOperaLog
	if err := db.Find(&rows).Error; err != nil {
		t.Fatalf("read back: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("want 1 row written, got %d", len(rows))
	}
	got := rows[0].JsonResult
	if !utf8.ValidString(got) {
		t.Errorf("stored JsonResult is not valid UTF-8: %q", got)
	}
	if n := utf8.RuneCountInString(got); n != 100 {
		t.Errorf("want 100 characters kept, got %d", n)
	}
}
