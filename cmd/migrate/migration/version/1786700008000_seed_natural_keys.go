package version

import (
	"fmt"
	"runtime"

	"gorm.io/gorm"

	adminmodels "go-admin/app/admin/models"
	"go-admin/cmd/migrate/migration"
	common "go-admin/common/models"
)

// Give seed.SeedMenus's two write paths (seedApis, seedMenuTree in
// app/admin/service/seed.go) a real natural key to check before inserting,
// so a retried, partially-failed migration (see the design doc
// docs-prd/008-应用清单与安装器/数据库变更.md §1.5/§1.6) does not insert the
// same row twice. This has already happened in production once (duplicate
// sys_menu/casbin_rule rows on the demo site), not a theoretical risk.
func init() {
	_, fileName, _, _ := runtime.Caller(0)
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _1786700008000SeedNaturalKeys)
}

func _1786700008000SeedNaturalKeys(db *gorm.DB, version string) error {
	if err := seedNaturalKeys(db); err != nil {
		return err
	}
	return db.Create(&common.Migration{Version: version}).Error
}

// seedNaturalKeys is split out from the wrapper above so tests can call it
// against a database that only has sys_menu/sys_api, without also standing
// up sys_migration - and so it can be called more than once in the same
// test to prove the re-run tolerance the doc comment above promises: DDL
// does not roll back on MySQL, so an operator whose first attempt failed
// partway through has nothing to do but run the whole migration again.
func seedNaturalKeys(db *gorm.DB) error {
	m := db.Migrator()

	// sys_menu.seed_code is a brand-new column: every existing row becomes
	// NULL, and NULL never collides in the unique index built below, so
	// this needs no pre-check.
	if !m.HasColumn(&adminmodels.SysMenu{}, "SeedCode") {
		if err := m.AddColumn(&adminmodels.SysMenu{}, "SeedCode"); err != nil {
			return err
		}
	}
	if !m.HasIndex(&adminmodels.SysMenu{}, "uk_sys_menu_app_seed_code_del") {
		if err := db.Exec(
			"CREATE UNIQUE INDEX uk_sys_menu_app_seed_code_del ON sys_menu (app_code, seed_code, deleted_at)",
		).Error; err != nil {
			return err
		}
	}

	// sys_api reuses existing, already-populated columns, which the demo
	// site has already proven can hold duplicates. Refuse rather than let
	// CREATE UNIQUE INDEX fail on an operator with no idea which rows to
	// reconcile - same shape as 1786700003000_soft_delete_marker.go's
	// refuseOnDuplicates.
	if err := refuseOnDuplicateApis(db); err != nil {
		return err
	}
	if !m.HasIndex(&adminmodels.SysApi{}, "uk_sys_api_app_path_action_del") {
		if err := db.Exec(
			"CREATE UNIQUE INDEX uk_sys_api_app_path_action_del ON sys_api (app_code, path, action, deleted_at)",
		).Error; err != nil {
			return err
		}
	}

	return nil
}

// refuseOnDuplicateApis reports the (app_code, path, action) values that
// would make the unique index impossible, rather than the index failing to
// build and saying only that it did. Only live rows count: a soft-deleted
// duplicate does not block the index it will never occupy a slot in.
//
// sys_api.path/action (app/admin/models/sys_api.go) carry no NOT NULL
// constraint, and that stays true here on purpose: tightening it is an
// independent, backward-incompatible change of its own - existing NULL
// rows in a real database would need reconciling or backfilling before
// ALTER TABLE ... NOT NULL could even run, which is a decision for
// whoever owns that data, not something this migration should force as a
// side effect of adding an unrelated index. So this function has to
// tolerate NULL path/action rather than assume they cannot occur - see the
// query below for how it does that without either crashing on them
// (MySQL's CONCAT) or wrongly flagging them (GROUP BY's NULL-equals-NULL).
//
// The two are independent bugs that happened to share one root cause, and
// SQLite's own test suite for this file would have caught neither on its
// own: MySQL's CONCAT() returns NULL if any argument is NULL, which turned
// a duplicate check against a NULL-holding library into "converting NULL
// to string is unsupported" instead of a report - but SQLite's (and
// PostgreSQL's) CONCAT() treats a NULL argument as an empty string
// instead, so the exact same query never errors there no matter how it is
// called. A suite that only ever ran on SQLite would report success for
// both defects; only a real MySQL server surfaces the first one at all -
// this migration's PostgreSQL-only sibling test file
// (1786700008000_seed_natural_keys_postgres_test.go) rules out one more
// dialect, but MySQL specifically has to be checked by hand, since this
// repository's test suite has no MySQL service to run against in CI.
func refuseOnDuplicateApis(db *gorm.DB) error {
	var dupes []string
	if err := db.Raw(
		// This has to agree with what the unique index it guards actually
		// enforces, not just with what looks like a duplicate at a glance.
		// Two different SQL rules collide on a NULL: GROUP BY treats two
		// NULLs as equal, so a naive query flags every pair of rows that
		// share a NULL path or action - even a pair with only one of the
		// two NULL, since GROUP BY's equality still holds on whichever
		// column both rows leave NULL - but a UNIQUE INDEX treats every
		// NULL as distinct from every other value, including another
		// NULL, so the index itself accepts every one of those pairs
		// without complaint. Excluding any row missing either column from
		// consideration entirely is what makes the two agree: a row
		// missing path, or missing action, or missing both, can never
		// violate the index no matter how many other rows are also
		// missing the same one, so none of them belong in this count.
		//
		// No COALESCE: with both columns excluded whenever either is
		// NULL, CONCAT here never receives a NULL argument for path or
		// action - app_code cannot be NULL at all (see its own NOT NULL
		// tag) - so there is nothing left for COALESCE to guard against,
		// and leaving it out is deliberate rather than an oversight. A
		// future regression that removed the two IS NOT NULL conditions
		// above would fail loudly on MySQL (the same Scan error this
		// query used to produce) instead of quietly reporting a made-up
		// "duplicate" whose path and action both print as empty - the
		// failure this function exists to prevent in the first place.
		`SELECT CONCAT(app_code, '|', path, '|', action) FROM sys_api
		 WHERE deleted_at = 0 AND path IS NOT NULL AND action IS NOT NULL
		 GROUP BY app_code, path, action HAVING COUNT(*) > 1`,
	).Scan(&dupes).Error; err != nil {
		return fmt.Errorf("checking sys_api for duplicates: %w", err)
	}
	if len(dupes) > 0 {
		return fmt.Errorf(
			"sys_api already holds duplicate (app_code,path,action) %v; reconcile them before this migration can add its unique index",
			dupes)
	}
	return nil
}
