package tools

import (
	"regexp"
	"strconv"
	"strings"

	"go-admin/app/other/models/tools"
)

// columnLengthPattern pulls the first parenthesized integer out of a MySQL
// COLUMN_TYPE string - the "(255)" in "varchar(255)", the "(10" in
// "decimal(10,2)". Works regardless of trailing modifiers such as
// "unsigned" or a charset clause, since it only looks for the first digits
// after the first '('.
var columnLengthPattern = regexp.MustCompile(`\((\d+)`)

// InferColumnWidth backs R2's fallback path: when a column's colWidth is
// left at its 0 sentinel (unconfigured), this reads sys_columns.column_type
// - MySQL's information_schema.COLUMNS.COLUMN_TYPE, which carries length,
// e.g. "varchar(255)", "int(11)", "decimal(10,2)", "tinyint(1)" - and
// returns a px width sized to fit inside go-admin-ui's ~580px text-column
// budget for a 1280px viewport (its AGENTS.md "列宽" section).
//
// The judgment has to be columnType, not goType: sys_tables.go:323-338
// gives every non-primary-key int/tinyint/bigint/decimal column goType
// "string" (a bare substring match on "int" that also catches "tinyint"/
// "bigint", intentional at import time but useless for telling a boolean
// flag from a bigint), so goType alone cannot distinguish a switch column
// from a price column from a name column. This is the same judgment call
// API契约.md §1.1 made, reversing the PRD's original "GoType" reading of R2.
// GoType is not consulted anywhere in this function, including for
// datetime/timestamp columns - those are matched on columnType too.
//
// Exported and pure (string in, int out) so QA can pin an exact input/output
// table against it directly (测试用例.md §2.5's own recommendation), rather
// than only being able to assert "the rendered page happens not to overflow".
func InferColumnWidth(columnType string) int {
	ct := strings.ToLower(strings.TrimSpace(columnType))

	switch {
	case strings.HasPrefix(ct, "tinyint(1)"):
		// MySQL's own shape for a boolean/status flag - a tag or a switch,
		// not text, so it wants less room than a general numeric column.
		return 70

	case strings.Contains(ct, "datetime"), strings.Contains(ct, "timestamp"),
		strings.Contains(ct, "date"), strings.Contains(ct, "time"):
		return 110

	case strings.HasPrefix(ct, "tinyint"), strings.HasPrefix(ct, "smallint"),
		strings.HasPrefix(ct, "mediumint"), strings.HasPrefix(ct, "int"),
		strings.HasPrefix(ct, "bigint"), strings.HasPrefix(ct, "decimal"),
		strings.HasPrefix(ct, "float"), strings.HasPrefix(ct, "double"):
		// API契约.md §1.1: "decimal/bigint/int 类给数字型窄宽度" groups these
		// together rather than sizing each individually - none of them need
		// more than a handful of digits' worth of width.
		return 90

	case strings.HasPrefix(ct, "varchar"), strings.HasPrefix(ct, "char"):
		return varcharWidth(columnLength(ct))

	case strings.Contains(ct, "text"), strings.Contains(ct, "blob"):
		// longtext/mediumtext/text/blob: no declared length to size against,
		// and content here is free-form, so this errs wide rather than
		// guessing a number the actual content will not respect.
		return 260

	default:
		// Unrecognized column_type (an enum, a json column, a driver this
		// codebase does not special-case, ...). Matches the flat fallback
		// vue.go.template already used for every non-datetime column before
		// this function existed, so a type this does not recognize is no
		// worse off than the old blanket default.
		return 120
	}
}

// varcharWidth tiers a char/varchar column by its declared length. The
// tiers are deliberately coarse - R2 only asks for "common tables land in
// the 580px budget", not pixel-perfect sizing per character.
func varcharWidth(n int) int {
	switch {
	case n <= 0:
		// Length did not parse (unexpected shape) - mid tier, not the
		// narrowest, since an un-lengthed varchar is unlikely to be a
		// short code column.
		return 150
	case n <= 10:
		return 90
	case n <= 20:
		return 110
	case n <= 50:
		return 150
	case n <= 100:
		return 200
	default:
		return 240
	}
}

// columnLength extracts the first parenthesized integer, or 0 if the type
// string does not have one (already-lowercased input expected).
func columnLength(columnType string) int {
	m := columnLengthPattern.FindStringSubmatch(columnType)
	if m == nil {
		return 0
	}
	n, err := strconv.Atoi(m[1])
	if err != nil {
		return 0
	}
	return n
}

// applyInferredColumnWidths fills in InferColumnWidth's result for every
// column still at the 0 "unconfigured" sentinel, in place, before the
// template that reads .ColWidth runs. A column the user (or F6's config
// page) already gave an explicit width is left untouched.
func applyInferredColumnWidths(columns []tools.SysColumns) {
	for i := range columns {
		if columns[i].ColWidth == 0 {
			columns[i].ColWidth = InferColumnWidth(columns[i].ColumnType)
		}
	}
}
