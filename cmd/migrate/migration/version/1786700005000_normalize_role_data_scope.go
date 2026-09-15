package version

import (
	"runtime"

	"gorm.io/gorm"

	"go-admin/cmd/migrate/migration"
	common "go-admin/common/models"
)

// Normalize sys_role.data_scope to one of the five values
// actions.Permission recognizes, ahead of PRD 006 F14/H2 making its
// unrecognized-scope branch fail closed instead of fail open.
//
// Before that change, an empty or unrecognized data_scope fell into
// Permission's default branch, which returned the query untouched - exactly
// the same SQL as data_scope "1" (全部数据权限). The seed data shipped
// precisely that: config/db.sql's built-in admin role (role_id 1) carries an
// empty data_scope rather than "1". Once the default starts matching no
// rows instead, that role would silently lose all visibility everywhere
// actions.Permission is used, the moment a deployment turns EnableDP on.
//
// Rewriting every value outside {1,2,3,4,5} to "1" keeps each such role's
// effective visibility exactly what it already was - a role that intended a
// tighter scope was never getting it under the old fail-open default either,
// so this does not tighten anything a deployment was relying on. Whether to
// tighten it further is left to whoever owns that role.
func init() {
	_, fileName, _, _ := runtime.Caller(0)
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _1786700005000NormalizeRoleDataScope)
}

func _1786700005000NormalizeRoleDataScope(db *gorm.DB, version string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := normalizeRoleDataScope(tx); err != nil {
			return err
		}
		return tx.Create(&common.Migration{Version: version}).Error
	})
}

// normalizeRoleDataScope is split out so tests can run it against a database
// that only has sys_role, without also standing up sys_migration.
//
// The explicit "IS NULL OR" matters: sys_role.data_scope has no NOT NULL
// constraint, and SQL's three-valued logic makes `NULL NOT IN (...)`
// evaluate to NULL rather than TRUE, so a bare NOT IN clause silently skips
// NULL rows instead of normalizing them.
func normalizeRoleDataScope(tx *gorm.DB) error {
	return tx.Exec(
		"UPDATE sys_role SET data_scope = '1' WHERE data_scope IS NULL OR data_scope NOT IN ('1', '2', '3', '4', '5')",
	).Error
}
