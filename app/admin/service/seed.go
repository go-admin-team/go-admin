package service

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"gorm.io/gorm"

	contractmodels "github.com/go-admin-team/go-admin-core/v2/sdk/contract/models"
	"github.com/go-admin-team/go-admin-core/v2/sdk/contract/seed"

	"go-admin/app/admin/models"
)

// adminSeeder is go-admin's own implementation of seed.Seeder: it turns the
// MenuSpec/ApiSpec values a third-party application asks for into rows
// across the four tables a visible, working menu entry needs - sys_api,
// sys_menu, sys_menu_api_rule, and sys_role_menu/casbin_rule - following the
// same shape cmd/migrate/migration/version/1786700001000_demo_menu.go
// already hand-writes for the host's own demo module.
//
// See go-admin-core's docs/contract.md, "Application-supplied menu and API
// entries", for the requirements this satisfies, and the security note on
// seed.Seeder for what this boundary does and does not protect against: an
// application already holds the same *gorm.DB this receives and could write
// sys_menu/sys_api/casbin_rule directly, bypassing this entirely.
type adminSeeder struct{}

func init() {
	seed.RegisterSeeder(adminSeeder{})
}

// adminRoleKey is the role every seeded menu is granted to. This mirrors
// 1786700001000_demo_menu.go's own convention rather than inventing a
// second one: MenuSpec carries no "which roles should see this" field for a
// Seeder to consult instead, and admin is the one role guaranteed to exist
// once the framework's own seed data has run.
const adminRoleKey = "admin"

// menuSortRange is what sys_menu.sort's column type actually holds.
//
// sort is `gorm:"size:4"`, which MySQL builds as a tinyint (-128..127);
// sqlite ignores the width and accepts anything, so this only ever surfaces
// on a real install, mid-migration, as Error 1264 - by which point the
// migration has already run other, non-transactional DDL that will not be
// retried. tools/checksilent's menu-sort-overflow check catches this for
// every MenuSpec-shaped literal committed to this repository, but it walks
// the repository's own source tree: a third-party application living in the
// module cache is invisible to it. This is the equivalent check for that
// application, run when its migration actually calls SeedMenus rather than
// never.
const (
	menuSortMin = -128
	menuSortMax = 127
)

func (adminSeeder) SeedMenus(tx *gorm.DB, appCode string, menus []seed.MenuSpec, apis []seed.ApiSpec) error {
	apiRows, err := seedApis(tx, appCode, apis)
	if err != nil {
		return fmt.Errorf("seed: app %q: apis: %w", appCode, err)
	}

	menuIDs, err := seedMenuTree(tx, appCode, menus, apiRows)
	if err != nil {
		return fmt.Errorf("seed: app %q: menus: %w", appCode, err)
	}

	// Not `len(menuIDs) == 0`: grantToAdminRole grants two independent
	// things, and an application is free to register apis without menus -
	// endpoints another service calls, or a UI mounted somewhere else.
	// Skipping the whole call on an empty menu list wrote the sys_api rows
	// and then no casbin rule for them, so those endpoints were denied to
	// everyone, admin included, with a migration that reported success.
	if len(menuIDs) == 0 && len(apiRows) == 0 {
		return nil
	}
	return grantToAdminRole(tx, menuIDs, apiRows)
}

// seedApis writes one sys_api row per ApiSpec and returns them keyed by
// ApiSpec.Code, so seedMenuTree can resolve a MenuSpec's ApiCodes into the
// rows sys_menu_api_rule needs to reference.
//
// sys_api.id is left to autoincrement rather than assigned by the caller,
// unlike 1786700001000_demo_menu.go's hand-picked ids: that migration is
// the one file tools/checksilent's menu-id-collision check can see, because
// it lives in this repository; nothing plays that role for a third-party
// application's ids in the module cache. Never accepting a caller-chosen id
// here removes the collision this Seeder has no way to detect instead of
// trying to detect it after the fact.
func seedApis(tx *gorm.DB, appCode string, apis []seed.ApiSpec) (map[string]models.SysApi, error) {
	seen := make(map[string]bool, len(apis))
	rows := make(map[string]models.SysApi, len(apis))
	for _, a := range apis {
		if a.Code == "" {
			return nil, errors.New("ApiSpec.Code must not be empty")
		}
		if seen[a.Code] {
			return nil, fmt.Errorf("duplicate ApiSpec.Code %q", a.Code)
		}
		seen[a.Code] = true

		row := models.SysApi{
			Handle:  a.Handle,
			Title:   a.Title,
			Path:    a.Path,
			Action:  a.Method,
			Type:    "SYS",
			AppCode: appCode,
		}
		if err := tx.Create(&row).Error; err != nil {
			return nil, fmt.Errorf("api %q: %w", a.Code, err)
		}
		rows[a.Code] = row
	}
	return rows, nil
}

// seedMenuTree writes one sys_menu row per MenuSpec, resolving Parent/Code
// references into parent_id/paths, and returns every menu id created so the
// caller can grant them to a role.
//
// Specs do not have to be given in parent-before-child order: this makes
// repeated passes over the remaining specs, creating whichever ones have
// their Parent (if any) already created, until every spec is placed. A
// spec whose Parent never resolves - naming a Code missing from this call,
// or only reachable through a cycle - stops making progress and is reported
// rather than looping forever.
func seedMenuTree(tx *gorm.DB, appCode string, specs []seed.MenuSpec, apiRows map[string]models.SysApi) ([]int, error) {
	byCode := make(map[string]seed.MenuSpec, len(specs))
	for _, s := range specs {
		if s.Code == "" {
			return nil, errors.New("MenuSpec.Code must not be empty")
		}
		if _, dup := byCode[s.Code]; dup {
			return nil, fmt.Errorf("duplicate MenuSpec.Code %q", s.Code)
		}
		if err := validateMenuSpec(s); err != nil {
			return nil, fmt.Errorf("%q: %w", s.Code, err)
		}
		byCode[s.Code] = s
	}

	created := make(map[string]models.SysMenu, len(specs))
	ids := make([]int, 0, len(specs))

	for len(created) < len(specs) {
		progressed := false
		for _, s := range specs {
			if _, done := created[s.Code]; done {
				continue
			}

			var parentRow models.SysMenu
			if s.Parent != "" {
				parent, ok := created[s.Parent]
				if !ok {
					if _, exists := byCode[s.Parent]; !exists {
						return nil, fmt.Errorf("%q: Parent %q is not a Code in this call", s.Code, s.Parent)
					}
					continue // s.Parent exists but has not been created yet; retry next pass
				}
				parentRow = parent
			}

			row := models.SysMenu{
				MenuName:   menuName(appCode, s.Code),
				Title:      s.Title,
				Icon:       s.Icon,
				Path:       s.Path,
				MenuType:   s.Kind,
				Permission: s.Permission,
				ParentId:   parentRow.MenuId,
				Component:  s.Component,
				Sort:       s.Sort,
				// Hidden by default and marked as an external frame, the
				// same defaults 1786700001000_demo_menu.go seeds its own
				// menu with: a freshly installed application's menu should
				// not need an administrator to first find and unhide it.
				Visible: "0",
				IsFrame: "1",
				AppCode: appCode,
			}
			for _, code := range s.ApiCodes {
				api, ok := apiRows[code]
				if !ok {
					return nil, fmt.Errorf("%q: ApiCodes references %q, which is not an ApiSpec.Code in this call", s.Code, code)
				}
				// The full row, not just {Id: api.Id}: gorm's many2many
				// association save upserts an associated row whose primary
				// key is already set, so a stub carrying only Id would
				// overwrite every other column of an sys_api row this same
				// call just wrote with zero values.
				row.SysApi = append(row.SysApi, api)
			}

			if err := tx.Create(&row).Error; err != nil {
				return nil, fmt.Errorf("%q: %w", s.Code, err)
			}

			// paths is a materialized path from the root ("/0"), built from
			// ids that only exist once the row above is created - the same
			// two-step create-then-update 1786700001000_demo_menu.go's
			// hand-assigned ids let it do in one literal, sequenced here
			// instead.
			if s.Parent == "" {
				row.Paths = "/0/" + strconv.Itoa(row.MenuId)
			} else {
				row.Paths = parentRow.Paths + "/" + strconv.Itoa(row.MenuId)
			}
			if err := tx.Model(&models.SysMenu{}).Where("menu_id = ?", row.MenuId).
				Update("paths", row.Paths).Error; err != nil {
				return nil, fmt.Errorf("%q: writing paths: %w", s.Code, err)
			}

			created[s.Code] = row
			ids = append(ids, row.MenuId)
			progressed = true
		}
		if !progressed {
			return nil, fmt.Errorf("unresolved Parent reference(s) among %d remaining spec(s); check for a cycle", len(specs)-len(created))
		}
	}
	return ids, nil
}

// validateMenuSpec rejects the malformed input tools/checksilent's
// menu-sort-overflow and Kind-adjacent checks would catch for an in-tree
// seed but cannot for a third-party application's - see menuSortRange's doc
// comment.
func validateMenuSpec(s seed.MenuSpec) error {
	switch s.Kind {
	case contractmodels.Directory, contractmodels.Menu, contractmodels.Button:
	default:
		return fmt.Errorf("Kind %q is not one of Directory/Menu/Button", s.Kind)
	}
	if s.Sort < menuSortMin || s.Sort > menuSortMax {
		return fmt.Errorf("Sort %d does not fit sys_menu.sort's tinyint column (%d..%d)", s.Sort, menuSortMin, menuSortMax)
	}
	return nil
}

// menuName synthesizes sys_menu.menu_name from appCode and the spec's Code,
// since MenuSpec carries no field of its own for it - contract/seed's
// package doc says a MenuSpec is what rendering a menu and checking a
// button permission need, not a mirror of sys_menu's columns.
//
// PascalCasing both and concatenating them, rather than using Code alone,
// is what keeps two applications that both picked the plain word "list" as
// a Code from producing the identical menu_name: the frontend's keep-alive
// cache matches a route by this exact string, not by (appCode, Code), so a
// collision there is a UI bug, not a database error, and nothing else here
// would ever surface it.
func menuName(appCode, code string) string {
	return pascalCase(appCode) + pascalCase(code)
}

func pascalCase(s string) string {
	var b strings.Builder
	for _, part := range strings.FieldsFunc(s, func(r rune) bool { return r == '-' || r == '_' }) {
		b.WriteString(strings.ToUpper(part[:1]))
		b.WriteString(part[1:])
	}
	return b.String()
}

// grantToAdminRole is sys_role_menu and casbin_rule: the two tables
// go-admin-core's contract.md requires alongside sys_menu/sys_api, without
// which a seeded menu is invisible to every role and its apis are
// authorized for no one.
//
// It follows 1786700001000_demo_menu.go's exact pattern, including
// tolerating a missing admin role: a database that has not yet run the
// framework's own seed data (config/db.sql, inside 1599190683659_tables.go)
// has nothing to grant to yet, and namespacedKey's ordering guarantee - every
// framework migration sorts before every app-prefixed one - means that
// should not happen in practice, but failing this call over it would be
// worse than a menu with no grant yet.
func grantToAdminRole(tx *gorm.DB, menuIDs []int, apiRows map[string]models.SysApi) error {
	var role models.SysRole
	if err := tx.Where("role_key = ?", adminRoleKey).First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	for _, id := range menuIDs {
		if err := tx.Exec(
			"INSERT INTO sys_role_menu (role_id, menu_id) SELECT ?, ? WHERE NOT EXISTS (SELECT 1 FROM sys_role_menu WHERE role_id = ? AND menu_id = ?)",
			role.RoleId, id, role.RoleId, id,
		).Error; err != nil {
			return err
		}
	}

	for _, a := range apiRows {
		if err := tx.Exec(
			"INSERT INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5) SELECT 'p', ?, ?, ?, '', '', '' WHERE NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype='p' AND v0=? AND v1=? AND v2=?)",
			role.RoleKey, a.Path, a.Action, role.RoleKey, a.Path, a.Action,
		).Error; err != nil {
			return err
		}
	}
	return nil
}
