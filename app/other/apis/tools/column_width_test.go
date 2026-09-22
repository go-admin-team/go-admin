package tools

import (
	"testing"

	"go-admin/app/other/models/tools"
)

// Input/output pins for InferColumnWidth, per 测试用例.md §2.5's own
// recommendation ("QA 才能在阶段 4 补一张精确的输入→输出对照表断言, 而不是只测
// 结果凑巧没溢出这种弱结论") - this is that table, kept next to the function
// it pins rather than only living in a later QA-owned suite.
func TestInferColumnWidth(t *testing.T) {
	cases := []struct {
		name       string
		columnType string
		want       int
	}{
		{"boolean/status flag", "tinyint(1)", 70},
		{"boolean flag, case-insensitive", "TINYINT(1)", 70},
		{"datetime", "datetime", 110},
		{"timestamp", "timestamp", 110},
		{"date only", "date", 110},
		{"time only", "time", 110},

		{"plain tinyint (not the (1) boolean shape)", "tinyint(4)", 90},
		{"smallint", "smallint(6)", 90},
		{"mediumint", "mediumint(9)", 90},
		{"int", "int(11)", 90},
		{"bigint", "bigint(20)", 90},
		{"decimal", "decimal(10,2)", 90},
		{"float", "float", 90},
		{"double", "double", 90},

		{"varchar short code", "varchar(8)", 90},
		{"varchar at the 10 boundary", "varchar(10)", 90},
		{"varchar just past the 10 boundary", "varchar(11)", 110},
		{"varchar at the 20 boundary", "varchar(20)", 110},
		{"varchar mid length", "varchar(32)", 150},
		{"varchar at the 50 boundary", "varchar(50)", 150},
		{"varchar just past the 50 boundary", "varchar(51)", 200},
		{"varchar(255), the common default", "varchar(255)", 240},
		{"char, fixed-width", "char(2)", 90},
		{"varchar with no parsed length", "varchar", 150},

		{"text, no length to size against", "text", 260},
		{"longtext", "longtext", 260},
		{"mediumtext", "mediumtext", 260},
		{"blob", "blob", 260},

		{"unrecognized type falls back to the old flat default", "json", 120},
		{"empty column_type falls back to the old flat default", "", 120},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := InferColumnWidth(tc.columnType); got != tc.want {
				t.Errorf("InferColumnWidth(%q) = %d, want %d", tc.columnType, got, tc.want)
			}
		})
	}
}

func TestApplyInferredColumnWidths(t *testing.T) {
	columns := []tools.SysColumns{
		{JsonField: "name", ColumnType: "varchar(64)", ColWidth: 0},
		{JsonField: "price", ColumnType: "decimal(10,2)", ColWidth: 300}, // already configured
	}

	applyInferredColumnWidths(columns)

	if columns[0].ColWidth == 0 {
		t.Error("unconfigured column: want an inferred non-zero width, still 0")
	}
	if want := InferColumnWidth("varchar(64)"); columns[0].ColWidth != want {
		t.Errorf("unconfigured column: want %d (InferColumnWidth's own answer), got %d", want, columns[0].ColWidth)
	}
	if columns[1].ColWidth != 300 {
		t.Errorf("already-configured column: want the user's 300 left untouched, got %d", columns[1].ColWidth)
	}
}
