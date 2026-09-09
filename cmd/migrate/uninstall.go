package migrate

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"gorm.io/gorm"

	adminmodels "go-admin/app/admin/models"
	"go-admin/cmd/migrate/migration"
	commonmodels "go-admin/common/models"
)

// policyKey is one casbin_rule row identified the way casbin_rule is unique:
// by its tuple, not by its id. Ids do not survive SysRole.Update, which
// removes a role's policies and adds them back.
type policyKey struct {
	Ptype string `gorm:"column:ptype"`
	V0    string `gorm:"column:v0"`
	V1    string `gorm:"column:v1"`
	V2    string `gorm:"column:v2"`
	V3    string `gorm:"column:v3"`
	V4    string `gorm:"column:v4"`
	V5    string `gorm:"column:v5"`
}

func (p policyKey) String() string {
	return fmt.Sprintf("%s %s %s %s", p.Ptype, p.V0, p.V1, p.V2)
}

// uninstallReport is what an uninstall removed, and what it deliberately did
// not.
type uninstallReport struct {
	Code string
	// Found says whether sys_app had a row. An application whose migrations
	// were applied by plain `migrate` rather than by `install` has its menus
	// and its permissions without ever having had one.
	Found      bool
	Version    string
	Menus      int64
	Apis       int64
	Bindings   int64
	RoleMenus  int64
	Policies   int64
	Migrations int64
	// Skipped are ledger entries whose casbin_rule row was not there any
	// more: something this install created and something else removed.
	Skipped []policyKey
	// Orphans are policies naming this application's paths that no ledger
	// entry claims - somebody granted this app's API to another role by
	// hand. Reported, never deleted.
	Orphans []policyKey
}

// uninstall removes one application's menus, APIs and permission grants.
//
// It does not touch the application's own tables. Removing an order module
// is not the same decision as destroying the orders, and nothing here can
// tell the operator apart from someone who will reinstall tomorrow.
//
// One transaction, and this one really is one: every statement below is DML
// or a SELECT, so unlike an install there is no DDL to commit it out from
// under itself. Child rows go first, while the ids that identify them can
// still be read from the parents.
//
// A sys_app row is not required. `migrate` with no subcommand applies every
// registered migration, an application's included, so an application can
// have all of its data without ever having gone through the installer.
func uninstall(db *gorm.DB, code string) (uninstallReport, error) {
	code = migration.NormalizeAppCode(code)
	rep := uninstallReport{Code: code}
	if code == "" {
		return rep, errors.New("no app code given")
	}
	if code == migration.FrameworkAppCode {
		return rep, fmt.Errorf("%q is the framework's own migrations; there is no uninstall for those", code)
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		row, found, err := loadApp(tx, code)
		if err != nil {
			return err
		}
		rep.Found = found
		if found {
			rep.Version = row.Version
		}

		// 1 and 2. Read before deleting: sys_api's rows are about to go, and
		// step 5b needs their paths.
		//
		// Unscoped throughout. A row this application wrote that somebody
		// soft-deleted from the UI is still this application's row, and
		// leaving it behind would leave its join rows pointing at it.
		var menuIDs []int
		if err := tx.Unscoped().Model(&adminmodels.SysMenu{}).
			Where("app_code = ?", code).Pluck("menu_id", &menuIDs).Error; err != nil {
			return fmt.Errorf("reading this app's menus: %w", err)
		}
		var apiIDs []int
		if err := tx.Unscoped().Model(&adminmodels.SysApi{}).
			Where("app_code = ?", code).Pluck("id", &apiIDs).Error; err != nil {
			return fmt.Errorf("reading this app's apis: %w", err)
		}
		var apiKeys []policyKey
		if err := tx.Unscoped().Model(&adminmodels.SysApi{}).
			Where("app_code = ?", code).
			Select("path as v1, action as v2").Scan(&apiKeys).Error; err != nil {
			return fmt.Errorf("reading this app's api paths: %w", err)
		}

		// 3. The many2many rows behind SysMenu.SysApi. Either side is enough
		// to make a row this application's.
		if len(menuIDs) > 0 || len(apiIDs) > 0 {
			q := tx.Table("sys_menu_api_rule")
			switch {
			case len(menuIDs) > 0 && len(apiIDs) > 0:
				q = q.Where("sys_menu_menu_id IN ? OR sys_api_id IN ?", menuIDs, apiIDs)
			case len(menuIDs) > 0:
				q = q.Where("sys_menu_menu_id IN ?", menuIDs)
			default:
				q = q.Where("sys_api_id IN ?", apiIDs)
			}
			res := q.Delete(nil)
			if res.Error != nil {
				return fmt.Errorf("removing menu/api bindings: %w", res.Error)
			}
			rep.Bindings = res.RowsAffected
		}

		// 4. Role assignments. menu_id is a surrogate key, so a row here can
		// only have come from a menu this application wrote - there is no
		// "looks like it but is not". That is why this needs no ledger, and
		// why a column on sys_role_menu would have been wrong: SysRole.Update
		// deletes a role's rows and writes them back through GORM's
		// many2many, which does not carry extra columns, so any such column
		// would be silently blanked the first time somebody edits a role.
		if len(menuIDs) > 0 {
			res := tx.Table("sys_role_menu").Where("menu_id IN ?", menuIDs).Delete(nil)
			if res.Error != nil {
				return fmt.Errorf("removing role assignments: %w", res.Error)
			}
			rep.RoleMenus = res.RowsAffected
		}

		// 5. Policies, by ledger, one at a time and by exact tuple.
		var grants []adminmodels.SysAppCasbinGrant
		if err := tx.Where("app_code = ?", code).Find(&grants).Error; err != nil {
			return fmt.Errorf("reading the grant ledger: %w", err)
		}
		for _, g := range grants {
			k := policyKey{Ptype: g.Ptype, V0: g.V0, V1: g.V1, V2: g.V2, V3: g.V3, V4: g.V4, V5: g.V5}
			res := tx.Table("casbin_rule").
				Where("ptype = ? AND v0 = ? AND v1 = ? AND v2 = ? AND v3 = ? AND v4 = ? AND v5 = ?",
					k.Ptype, k.V0, k.V1, k.V2, k.V3, k.V4, k.V5).
				Delete(nil)
			if res.Error != nil {
				return fmt.Errorf("removing policy %s: %w", k, res.Error)
			}
			if res.RowsAffected == 0 {
				// Something this install created is not there any more. Not
				// an error: the uninstall's job was to remove it and it is
				// gone. Reported because a policy this app created and did
				// not remove means something else rewrote casbin_rule.
				rep.Skipped = append(rep.Skipped, k)
				continue
			}
			rep.Policies += res.RowsAffected
		}
		// The ledger's job ends here whether or not each row matched. Left
		// behind it would only grow, and a reinstall writes its own entries.
		if err := tx.Where("app_code = ?", code).
			Delete(&adminmodels.SysAppCasbinGrant{}).Error; err != nil {
			return fmt.Errorf("clearing the grant ledger: %w", err)
		}

		// 5b. Read-only. By now every policy the ledger could speak for has
		// been dealt with, so a policy still matching one of this app's paths
		// is one the ledger never claimed - somebody granted this app's API
		// to another role by hand. Business rule 3 says do not delete what
		// is not ours; without this step nobody would ever learn it is
		// there, pointing at an API that is about to stop existing.
		orphans, err := findOrphanPolicies(tx, apiKeys)
		if err != nil {
			return err
		}
		rep.Orphans = orphans

		// 6 and 7.
		res := tx.Unscoped().Where("app_code = ?", code).Delete(&adminmodels.SysApi{})
		if res.Error != nil {
			return fmt.Errorf("removing this app's apis: %w", res.Error)
		}
		rep.Apis = res.RowsAffected

		res = tx.Unscoped().Where("app_code = ?", code).Delete(&adminmodels.SysMenu{})
		if res.Error != nil {
			return fmt.Errorf("removing this app's menus: %w", res.Error)
		}
		rep.Menus = res.RowsAffected

		// 8. Without this a reinstall finds every version already applied,
		// runs no migration, seeds nothing, and reports success. It is the
		// easiest step to leave out, because a migration record does not
		// look like the application's data.
		res = tx.Where("app_code = ?", code).Delete(&commonmodels.Migration{})
		if res.Error != nil {
			return fmt.Errorf("removing this app's migration records: %w", res.Error)
		}
		rep.Migrations = res.RowsAffected

		// 9.
		if found {
			if err := tx.Where("app_code = ?", code).
				Delete(&adminmodels.SysApp{}).Error; err != nil {
				return fmt.Errorf("removing the sys_app row: %w", err)
			}
		}
		return nil
	})
	if err != nil {
		return uninstallReport{Code: code}, err
	}
	return rep, nil
}

// findOrphanPolicies looks for policies naming any of this application's
// paths.
//
// Written as an OR chain rather than a row-value IN, which MySQL and modern
// SQLite accept and SQL Server does not; this repository supports all of
// them. Chunked because a driver's placeholder limit is reached long before
// an application runs out of endpoints.
func findOrphanPolicies(tx *gorm.DB, keys []policyKey) ([]policyKey, error) {
	const chunk = 100
	var out []policyKey
	for start := 0; start < len(keys); start += chunk {
		end := start + chunk
		if end > len(keys) {
			end = len(keys)
		}
		clauses := make([]string, 0, end-start)
		args := make([]any, 0, (end-start)*2)
		for _, k := range keys[start:end] {
			clauses = append(clauses, "(v1 = ? AND v2 = ?)")
			args = append(args, k.V1, k.V2)
		}
		var found []policyKey
		if err := tx.Table("casbin_rule").
			Where("ptype = ? AND ("+strings.Join(clauses, " OR ")+")", append([]any{"p"}, args...)...).
			Scan(&found).Error; err != nil {
			return nil, fmt.Errorf("looking for policies nothing claims: %w", err)
		}
		out = append(out, found...)
	}
	return out, nil
}

// reportUninstall prints what went and what stayed.
//
// The two lists are separate because they mean different things: one is
// something this application created that had already gone, the other is
// somebody else's grant that is now pointing at nothing. Merged into one
// "could not remove" list, neither would be actionable.
func reportUninstall(w io.Writer, rep uninstallReport) {
	if !rep.Found {
		fmt.Fprintf(w, "%s had no sys_app row; removed what its migrations had written\n", rep.Code)
	} else {
		fmt.Fprintf(w, "uninstalled %s %s\n", rep.Code, rep.Version)
	}
	fmt.Fprintf(w, "removed: %d menu(s), %d api(s), %d binding(s), %d role assignment(s), %d policy(ies), %d migration record(s)\n",
		rep.Menus, rep.Apis, rep.Bindings, rep.RoleMenus, rep.Policies, rep.Migrations)
	fmt.Fprintln(w, "the application's own tables were not touched.")

	if len(rep.Skipped) > 0 {
		fmt.Fprintf(w, "\n%d policy(ies) this install had created were already gone:\n", len(rep.Skipped))
		for _, k := range rep.Skipped {
			fmt.Fprintf(w, "  %s\n", k)
		}
	}
	if len(rep.Orphans) > 0 {
		fmt.Fprintf(w, "\n%d policy(ies) name this application's paths and were granted by somebody else, so they were left alone:\n", len(rep.Orphans))
		for _, k := range rep.Orphans {
			fmt.Fprintf(w, "  %s\n", k)
		}
		fmt.Fprintln(w, "they now point at APIs that no longer exist. Harmless to the running server, and yours to clear up.")
	}
}
