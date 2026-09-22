package tools

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// GORM's Updates(struct) skips zero-value fields, and PRD 010 F1/F2 chose 0 /
// "" as the sentinel for "unconfigured" (docs-prd/010-代码生成器前端模板迁移Vue3/
// 数据库变更.md §1.1). Put those together and Update can set ColWidth/
// DefaultValue but never clear them back to the sentinel: the struct-form
// Updates call silently drops the very values this feature needs to write.
func TestSysColumnsUpdateClearsSentinelFields(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(new(SysColumns)); err != nil {
		t.Fatalf("migrate sys_columns: %v", err)
	}

	col := SysColumns{TableId: 1, ColumnName: "status", ColWidth: 150, DefaultValue: "active"}
	if _, err := col.Create(db); err != nil {
		t.Fatalf("create: %v", err)
	}

	// Reset back to the sentinel - the UI action for "go back to inferred
	// width / no default", not merely "never configured".
	update := SysColumns{ColumnId: col.ColumnId, ColWidth: 0, DefaultValue: ""}
	if _, err := update.Update(db); err != nil {
		t.Fatalf("update: %v", err)
	}

	var got SysColumns
	if err := db.Table("sys_columns").First(&got, col.ColumnId).Error; err != nil {
		t.Fatalf("read back: %v", err)
	}
	if got.ColWidth != 0 {
		t.Errorf("colWidth: want 0 (cleared), got %d - Update() did not write the sentinel back", got.ColWidth)
	}
	if got.DefaultValue != "" {
		t.Errorf("defaultValue: want \"\" (cleared), got %q - Update() did not write the sentinel back", got.DefaultValue)
	}
}
