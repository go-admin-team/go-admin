package version

import (
	"strings"
	"testing"
)

// The index has to mean the same thing on every driver this repository
// registers, and the drivers do not agree about NULL.
//
// MySQL, PostgreSQL and SQLite treat two NULLs as different values, so any
// number of rows missing one of these columns coexist under the index. SQL
// Server treats them as equal and permits exactly one, so the unfiltered
// statement fails there on any database with two rows lacking a seed_code -
// which is every database, a brand-new one included, because 1786700001000
// seeds five menus and none of them carries one.
func TestUniqueIndexOverNullableFiltersOnlyWhereItHasTo(t *testing.T) {
	const plain = "CREATE UNIQUE INDEX uk ON sys_menu (app_code, seed_code, deleted_at)"

	for _, dialect := range []string{"mysql", "postgres", "sqlite"} {
		got := uniqueIndexOverNullable(dialect, "uk", "sys_menu", "app_code, seed_code, deleted_at", "seed_code")
		if got != plain {
			t.Errorf("%s: %q\n  want %q", dialect, got, plain)
		}
	}

	got := uniqueIndexOverNullable("sqlserver", "uk", "sys_menu", "app_code, seed_code, deleted_at", "seed_code")
	want := plain + " WHERE seed_code IS NOT NULL"
	if got != want {
		t.Errorf("sqlserver: %q\n  want %q", got, want)
	}
}

// sys_api's key has two nullable columns, and either one being NULL is enough
// to collide on SQL Server.
func TestUniqueIndexOverNullableCoversEveryNullableColumn(t *testing.T) {
	got := uniqueIndexOverNullable("sqlserver", "uk", "sys_api",
		"app_code, path, action, deleted_at", "path", "action")
	if !strings.HasSuffix(got, " WHERE path IS NOT NULL AND action IS NOT NULL") {
		t.Errorf("got %q", got)
	}
}

// A key with nothing nullable in it needs no filter anywhere, or SQL Server
// would get a WHERE clause naming no column.
func TestUniqueIndexOverNullableWithoutNullableColumns(t *testing.T) {
	got := uniqueIndexOverNullable("sqlserver", "uk", "sys_menu", "app_code, deleted_at")
	if strings.Contains(got, "WHERE") {
		t.Errorf("got %q", got)
	}
}
