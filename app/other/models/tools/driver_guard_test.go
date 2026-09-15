package tools

import (
	"strings"
	"testing"

	"github.com/go-admin-team/go-admin-core/v2/sdk/config"
	"gorm.io/gorm"
)

// The generator reads MySQL's information_schema and nothing else. Every entry
// point says so, but the guard used to be written pkg.Assert(true, ...), which
// never fires: pkg.Assert panics when its condition is false. A non-mysql
// deployment therefore fell through to a query built on a zero-value *gorm.DB.
//
// tx is nil on purpose - the assertion has to come before anything touches it.
func TestCodegenModelsRefuseNonMySQLDrivers(t *testing.T) {
	previous := config.DatabaseConfig.Driver
	config.DatabaseConfig.Driver = "postgres"
	t.Cleanup(func() { config.DatabaseConfig.Driver = previous })

	cases := map[string]func(*gorm.DB){
		"DBTables.GetPage": func(tx *gorm.DB) {
			_, _, _ = new(DBTables).GetPage(tx, 10, 1)
		},
		"DBTables.Get": func(tx *gorm.DB) {
			_, _ = (&DBTables{TableName: "sys_user"}).Get(tx)
		},
		"DBColumns.GetPage": func(tx *gorm.DB) {
			_, _, _ = (&DBColumns{TableName: "sys_user"}).GetPage(tx, 10, 1)
		},
		"DBColumns.GetList": func(tx *gorm.DB) {
			_, _ = (&DBColumns{TableName: "sys_user"}).GetList(tx)
		},
	}

	for name, call := range cases {
		t.Run(name, func(t *testing.T) {
			defer func() {
				raised := recover()
				if raised == nil {
					t.Fatal("driver is not mysql and the call went through anyway")
				}
				msg, ok := raised.(string)
				if !ok || !strings.Contains(msg, "目前只支持mysql数据库") {
					t.Fatalf("want the mysql-only assertion, got %v", raised)
				}
			}()
			call(nil)
		})
	}
}
