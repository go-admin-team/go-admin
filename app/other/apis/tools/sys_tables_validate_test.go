package tools

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"go-admin/app/other/models/tools"
)

func TestValidateAndSanitizeColumns_JsonFieldFormat(t *testing.T) {
	cases := []struct {
		name      string
		jsonField string
		wantErr   bool
	}{
		{"lower camelCase", "userName", false},
		{"two-letter lowercase", "id", false},
		// The importer's own output (sys_tables.go's namelist/JsonField
		// loop), not made up: a single-letter column ("x"), and a column
		// whose last name segment ends in a digit ("address2", "a1") both
		// produce a jsonField with no separator left to re-capitalize.
		// These three used to be rejected - the whole point of this fix.
		{"single letter, real importer output for a column named x", "x", false},
		{"letters then a trailing digit, real importer output for address2", "address2", false},
		{"two letters then a digit, real importer output for a1", "a1", false},
		{"leading underscore rejected", "_id", true},
		{"leading digit rejected (not a legal identifier start)", "1name", true},
		{"snake_case rejected (importer never emits an underscore)", "user_name", true},
		{"dot rejected, would break the gen/{pkg}/{biz}.ts key path", "user.name", true},
		{"empty rejected", "", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateAndSanitizeColumns([]tools.SysColumns{{JsonField: tc.jsonField}})
			if tc.wantErr && err == nil {
				t.Errorf("jsonField %q: want error, got nil", tc.jsonField)
			}
			if !tc.wantErr && err != nil {
				t.Errorf("jsonField %q: want no error, got %v", tc.jsonField, err)
			}
		})
	}
}

func TestValidateAndSanitizeColumns_JsonFieldUniqueWithinTable(t *testing.T) {
	err := validateAndSanitizeColumns([]tools.SysColumns{
		{JsonField: "name"},
		{JsonField: "name"},
	})
	if err == nil {
		t.Fatal("want error for a jsonField repeated in the same table, got nil")
	}
}

func TestValidateAndSanitizeColumns_ColWidthOutOfRangeIsSanitizedNotRejected(t *testing.T) {
	cases := []struct {
		name  string
		width int
		want  int
	}{
		{"zero (unconfigured) is left alone", 0, 0},
		{"in range is left alone", 150, 150},
		{"lower bound is left alone", colWidthMin, colWidthMin},
		{"upper bound is left alone", colWidthMax, colWidthMax},
		{"too small falls back to the sentinel", colWidthMin - 1, 0},
		{"too large falls back to the sentinel", colWidthMax + 1, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cols := []tools.SysColumns{{JsonField: "name", ColWidth: tc.width}}
			if err := validateAndSanitizeColumns(cols); err != nil {
				t.Fatalf("colWidth %d: want no error (out-of-range sanitizes, it does not reject), got %v", tc.width, err)
			}
			if cols[0].ColWidth != tc.want {
				t.Errorf("colWidth %d: want sanitized to %d, got %d", tc.width, tc.want, cols[0].ColWidth)
			}
		})
	}
}

func TestValidateAndSanitizeColumns_DefaultValueExpressionRejected(t *testing.T) {
	cases := []struct {
		name         string
		defaultValue string
		wantErr      bool
	}{
		{"plain literal", "0", false},
		{"plain string literal", "active", false},
		{"empty (unconfigured)", "", false},
		{"function call rejected", "Date.now()", true},
		{"template literal rejected", "`x`", true},
		{"arrow function rejected", "() => 1", true},
		{"statement separator rejected", "1; drop", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateAndSanitizeColumns([]tools.SysColumns{{JsonField: "name", DefaultValue: tc.defaultValue}})
			if tc.wantErr && err == nil {
				t.Errorf("defaultValue %q: want error, got nil", tc.defaultValue)
			}
			if !tc.wantErr && err != nil {
				t.Errorf("defaultValue %q: want no error, got %v", tc.defaultValue, err)
			}
		})
	}
}

func newBusinessNameTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(new(tools.SysTables)); err != nil {
		t.Fatalf("migrate sys_tables: %v", err)
	}
	return db
}

func TestValidateBusinessNameUnique(t *testing.T) {
	db := newBusinessNameTestDB(t)

	existing := tools.SysTables{TBName: "sys_widget", PackageName: "biz", BusinessName: "widget"}
	if err := db.Table("sys_tables").Create(&existing).Error; err != nil {
		t.Fatalf("seed: %v", err)
	}

	t.Run("same package, same businessName, different table: rejected", func(t *testing.T) {
		other := tools.SysTables{TBName: "sys_widget_copy", PackageName: "biz", BusinessName: "widget"}
		if err := db.Table("sys_tables").Create(&other).Error; err != nil {
			t.Fatalf("seed second row: %v", err)
		}
		// Unscoped: a plain Delete only soft-deletes (SysTables carries
		// common.ModelTime), which would leave this row's businessName
		// looking taken for the next subtest - production's own delete path
		// (SysTables.BatchDelete) hard-deletes for the same reason.
		defer db.Table("sys_tables").Unscoped().Delete(&other)

		if err := validateBusinessNameUnique(db, "biz", "widget", other.TableId); err == nil {
			t.Error("want error for a businessName already used by another table in the same package, got nil")
		}
	})

	t.Run("different package, same businessName: allowed", func(t *testing.T) {
		if err := validateBusinessNameUnique(db, "other-pkg", "widget", 0); err != nil {
			t.Errorf("want no error across different packages, got %v", err)
		}
	})

	t.Run("a table checking against its own current name: allowed", func(t *testing.T) {
		if err := validateBusinessNameUnique(db, "biz", "widget", existing.TableId); err != nil {
			t.Errorf("want no error when the only match is the row being saved itself, got %v", err)
		}
	})
}
