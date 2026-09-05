// Package migration registers app-order's one migration: create its two
// tables and seed the menu/API entries the admin UI needs to expose them.
//
// It registers through contract/migration.ForApp - the package-level
// facade, not a private NewRegistry() - because that is the only registry a
// third-party app, which cannot reach into the host process, can register
// against and have any hope of the host's own execution engine picking up.
// Whether it actually does, today, is a different question: see this
// package's test file and the gap list in the accompanying report.
package migration

import (
	"gorm.io/gorm"

	contractmigration "github.com/go-admin-team/go-admin-core/v2/sdk/contract/migration"
	contractmodels "github.com/go-admin-team/go-admin-core/v2/sdk/contract/models"
	"github.com/go-admin-team/go-admin-core/v2/sdk/contract/seed"

	"github.com/go-admin-team/example-app-order/models"
)

// AppCode is app-order's migration.ForApp / seed.SeedMenus identity.
const AppCode = "order"

// version is this migration's sys_migration key before ForApp namespaces
// it (see contract/migration.ForApp's doc comment: the stored key becomes
// "order-" + version). It follows the framework's own 13-digit millisecond
// timestamp convention purely so a human reading sys_migration.version
// alongside the framework's own rows can still eyeball roughly when it was
// authored; contract/migration.ForApp does not require that shape, just
// uniqueness within this app's own namespace.
const version = "1793800000000"

func init() {
	contractmigration.ForApp(AppCode).SetVersion(version, createOrderSchema)
}

// createOrderSchema creates app_order/app_order_item and seeds the menu and
// API entries a host's Seeder turns into sys_menu/sys_api/sys_menu_api_rule
// rows (and, once an administrator grants the menu to a role through the
// ordinary admin UI, casbin_rule). See seed.Seeder's security note: this
// call does not sandbox anything, it only saves app-order from needing to
// know go-admin's own schema.
func createOrderSchema(db *gorm.DB, migrationVersion, appCode string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.AutoMigrate(&models.Order{}, &models.OrderItem{}); err != nil {
			return err
		}

		menus := []seed.MenuSpec{
			{
				Code: "dir", Kind: contractmodels.Directory,
				Title: "Order Example", Path: "/apps/order", Component: "Layout",
				Icon: "shopping", Sort: 20,
			},
			{
				Code: "list", Parent: "dir", Kind: contractmodels.Menu,
				Title: "Orders", Path: "list",
				// Component must start with "apps/<code>/" - see
				// seed.MenuSpec.Component's doc comment. This is the one
				// concrete rule the report's gap list has nothing bad to
				// say about: it is documented exactly where a caller
				// building a MenuSpec would look.
				Component: "apps/order/order/index",
				Sort:      1,
				ApiCodes:  []string{"list", "get", "create", "pay"},
			},
			{
				Code: "btn-create", Parent: "list", Kind: contractmodels.Button,
				Title: "Create", Permission: "order:order:create", Sort: 1,
			},
			{
				Code: "btn-pay", Parent: "list", Kind: contractmodels.Button,
				Title: "Pay", Permission: "order:order:pay", Sort: 2,
			},
		}
		apis := []seed.ApiSpec{
			{Code: "list", Title: "Order list", Path: "/api/v1/order", Method: "GET", Handle: "apis.Order.GetPage-fm"},
			{Code: "get", Title: "Order detail", Path: "/api/v1/order/:id", Method: "GET", Handle: "apis.Order.Get-fm"},
			{Code: "create", Title: "Create order", Path: "/api/v1/order", Method: "POST", Handle: "apis.Order.Create-fm"},
			{Code: "pay", Title: "Pay order", Path: "/api/v1/order/:id/pay", Method: "PUT", Handle: "apis.Order.Pay-fm"},
		}
		if err := seed.SeedMenus(tx, appCode, menus, apis); err != nil {
			return err
		}

		return tx.Create(&contractmodels.Migration{Version: migrationVersion, AppCode: appCode}).Error
	})
}
