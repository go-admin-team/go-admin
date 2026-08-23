package version

import (
	"fmt"
	"runtime"
	"time"

	"gorm.io/gorm"

	"go-admin/cmd/migrate/migration"
	common "go-admin/common/models"
)

// Convert deleted_at from a nullable timestamp to a non-null millisecond
// marker, then put a unique index on the three natural keys.
//
// Those keys had no unique index at all. Uniqueness was a SELECT COUNT
// followed by an INSERT, which two concurrent requests both pass.
//
// The index cannot be on the key alone: a soft-deleted row keeps occupying the
// name, so a deleted user's username could never be used again. It has to
// include the delete marker — and the marker has to be non-null, because two
// live rows are (name, NULL) and (name, NULL), and NULL is not equal to NULL.
// An index over a nullable marker permits both rows. It looks like a
// constraint and enforces nothing.
//
// DDL does not roll back on MySQL, so this is written to be re-runnable rather
// than transactional: every step asks whether it has already been taken.
func init() {
	_, fileName, _, _ := runtime.Caller(0)
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _1786700003000SoftDeleteMarker)
}

// Every table whose model embeds common.ModelTime.
var softDeleteTables = []string{
	"sys_api", "sys_config", "sys_dept", "sys_dict_data", "sys_dict_type",
	"sys_job", "sys_menu", "sys_post", "sys_role", "sys_user", "demo_product",
}

// The keys that gain a constraint, and the column that makes it possible.
var naturalKeys = []struct {
	table, column, index string
}{
	{"sys_user", "username", "uk_sys_user_username"},
	{"sys_role", "role_key", "uk_sys_role_role_key"},
	{"sys_dict_type", "dict_type", "uk_sys_dict_type_dict_type"},
}

const tempColumn = "deleted_at_ms"

func _1786700003000SoftDeleteMarker(db *gorm.DB, version string) error {
	// Refused before anything is altered: creating the index on a table that
	// already holds duplicates fails halfway through, and the operator is left
	// guessing which rows to reconcile.
	if err := refuseOnDuplicates(db); err != nil {
		return err
	}

	for _, table := range softDeleteTables {
		if err := convertDeletedAt(db, table); err != nil {
			return fmt.Errorf("%s: %w", table, err)
		}
	}

	for _, k := range naturalKeys {
		if err := createUniqueIndex(db, k.table, k.column, k.index); err != nil {
			return fmt.Errorf("%s.%s: %w", k.table, k.column, err)
		}
	}

	return db.Create(&common.Migration{Version: version}).Error
}

// refuseOnDuplicates reports the values that would make the index impossible,
// rather than the index failing to build and saying only that it did.
func refuseOnDuplicates(db *gorm.DB) error {
	for _, k := range naturalKeys {
		if !db.Migrator().HasTable(k.table) {
			continue
		}

		var dupes []string
		q := fmt.Sprintf(
			"SELECT %s FROM %s WHERE deleted_at IS NULL GROUP BY %s HAVING COUNT(*) > 1",
			k.column, k.table, k.column,
		)
		if !db.Migrator().HasColumn(k.table, "deleted_at") {
			// Already converted; live rows carry zero rather than null.
			q = fmt.Sprintf(
				"SELECT %s FROM %s WHERE deleted_at = 0 GROUP BY %s HAVING COUNT(*) > 1",
				k.column, k.table, k.column,
			)
		}
		if err := db.Raw(q).Scan(&dupes).Error; err != nil {
			return fmt.Errorf("checking %s.%s for duplicates: %w", k.table, k.column, err)
		}
		if len(dupes) > 0 {
			return fmt.Errorf(
				"%s.%s already holds duplicates %v; reconcile them before this migration can add its unique index",
				k.table, k.column, dupes,
			)
		}
	}
	return nil
}

func convertDeletedAt(db *gorm.DB, table string) error {
	m := db.Migrator()
	if !m.HasTable(table) {
		return nil
	}

	converted, err := isConverted(db, table)
	if err != nil {
		return err
	}
	if converted {
		return nil
	}

	if !m.HasColumn(table, tempColumn) {
		if err := db.Exec(addBigIntColumn(db, table, tempColumn)).Error; err != nil {
			return fmt.Errorf("adding %s: %w", tempColumn, err)
		}
	}

	// Converted in Go rather than in SQL: turning a timestamp into epoch
	// milliseconds is spelled differently by every dialect this supports, and
	// the row counts here do not justify four versions of it.
	if err := copyTimestamps(db, table); err != nil {
		return err
	}

	// SQLite refuses to drop a column an index still refers to, where MySQL and
	// PostgreSQL drop the dependent indexes along with it. Drop them first, on
	// every dialect: the column is about to be replaced by one of the same name,
	// and gorm recreates the index from the model's tag.
	if err := dropIndexesOn(db, table, "deleted_at"); err != nil {
		return err
	}

	if err := db.Exec(dropColumn(db, table, "deleted_at")).Error; err != nil {
		return fmt.Errorf("dropping deleted_at: %w", err)
	}
	if err := db.Exec(renameColumn(db, table, tempColumn, "deleted_at")).Error; err != nil {
		return fmt.Errorf("renaming %s: %w", tempColumn, err)
	}
	return nil
}

// isConverted reports whether deleted_at already holds the marker. A table
// mid-conversion still has both columns, and is finished rather than skipped.
func isConverted(db *gorm.DB, table string) (bool, error) {
	types, err := db.Migrator().ColumnTypes(table)
	if err != nil {
		return false, err
	}
	for _, c := range types {
		if c.Name() != "deleted_at" {
			continue
		}
		if nullable, ok := c.Nullable(); ok && !nullable {
			return !db.Migrator().HasColumn(table, tempColumn), nil
		}
		return false, nil
	}
	// No such column: nothing to convert.
	return true, nil
}

func copyTimestamps(db *gorm.DB, table string) error {
	type row struct {
		Id        int64
		DeletedAt *time.Time
	}

	// The tables walked here do not agree on what the key is called: sys_dept
	// keys on dept_id, sys_user on user_id, and only some on id.
	key, err := primaryKeyOf(db, table)
	if err != nil {
		return err
	}

	var rows []row
	q := fmt.Sprintf("SELECT %s AS id, deleted_at FROM %s WHERE deleted_at IS NOT NULL", key, table)
	if err := db.Raw(q).Scan(&rows).Error; err != nil {
		return fmt.Errorf("reading deleted rows: %w", err)
	}

	for _, r := range rows {
		if r.DeletedAt == nil {
			continue
		}
		u := fmt.Sprintf("UPDATE %s SET %s = ? WHERE %s = ?", table, tempColumn, key)
		if err := db.Exec(u, r.DeletedAt.UnixMilli(), r.Id).Error; err != nil {
			return fmt.Errorf("marking row %d deleted: %w", r.Id, err)
		}
	}
	return nil
}

// primaryKeyOf asks the database which column is the primary key, rather than
// assuming the name.
func primaryKeyOf(db *gorm.DB, table string) (string, error) {
	columns, err := db.Migrator().ColumnTypes(table)
	if err != nil {
		return "", fmt.Errorf("reading columns of %s: %w", table, err)
	}
	for _, c := range columns {
		if isKey, ok := c.PrimaryKey(); ok && isKey {
			return c.Name(), nil
		}
	}
	return "", fmt.Errorf("no primary key on %s", table)
}

func createUniqueIndex(db *gorm.DB, table, column, name string) error {
	if db.Migrator().HasIndex(table, name) {
		return nil
	}
	return db.Exec(fmt.Sprintf(
		"CREATE UNIQUE INDEX %s ON %s (%s, deleted_at)", name, table, column,
	)).Error
}

// dropIndexesOn removes every index on table that mentions column, so the
// column can be dropped. Named indexes are asked of the migrator rather than
// guessed, because the names differ between a schema gorm created and one a
// hand-written migration did.
func dropIndexesOn(db *gorm.DB, table, column string) error {
	names, err := indexNamesFor(db, table, column)
	if err != nil {
		return fmt.Errorf("listing indexes on %s.%s: %w", table, column, err)
	}
	m := db.Migrator()
	for _, name := range names {
		if !m.HasIndex(table, name) {
			continue
		}
		if err := m.DropIndex(table, name); err != nil {
			return fmt.Errorf("dropping index %s: %w", name, err)
		}
	}
	return nil
}

// indexNamesFor asks the database which indexes cover column.
func indexNamesFor(db *gorm.DB, table, column string) ([]string, error) {
	indexes, err := db.Migrator().GetIndexes(table)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, idx := range indexes {
		for _, c := range idx.Columns() {
			if c == column {
				names = append(names, idx.Name())
				break
			}
		}
	}
	return names, nil
}

// The three statements every dialect spells differently.

func addBigIntColumn(db *gorm.DB, table, column string) string {
	if db.Dialector.Name() == "sqlserver" {
		return fmt.Sprintf("ALTER TABLE %s ADD %s BIGINT NOT NULL DEFAULT 0", table, column)
	}
	return fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s BIGINT NOT NULL DEFAULT 0", table, column)
}

func dropColumn(db *gorm.DB, table, column string) string {
	return fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s", table, column)
}

func renameColumn(db *gorm.DB, table, from, to string) string {
	switch db.Dialector.Name() {
	case "sqlserver":
		return fmt.Sprintf("EXEC sp_rename '%s.%s', '%s', 'COLUMN'", table, from, to)
	case "mysql":
		return fmt.Sprintf("ALTER TABLE %s CHANGE %s %s BIGINT NOT NULL DEFAULT 0", table, from, to)
	default:
		return fmt.Sprintf("ALTER TABLE %s RENAME COLUMN %s TO %s", table, from, to)
	}
}
