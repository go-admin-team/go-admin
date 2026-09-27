package version

import (
	"os"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	adminmodels "go-admin/app/admin/models"
)

// bindGenMenuApis writes casbin_rule with INSERT ... SELECT ... WHERE NOT
// EXISTS and sys_menu_api_rule through gorm's association code; both are
// dialect-sensitive, and the other tests for it run on SQLite only. This runs
// the same grant on every database CI has.

const mysqlDSNEnv = "GO_ADMIN_TEST_MYSQL_DSN"

func mysqlDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := os.Getenv(mysqlDSNEnv)
	if dsn == "" {
		if os.Getenv("CI") != "" {
			t.Fatalf("%s is not set while CI is: the MySQL migration tests must not skip here", mysqlDSNEnv)
		}
		t.Skipf("%s is not set; skipping the MySQL migration tests", mysqlDSNEnv)
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("connecting to %s: %v", mysqlDSNEnv, err)
	}
	return db
}

// casbinRuleRow is the shape gorm-adapter gives casbin_rule, declared here so
// each dialect builds the table its own way.
type casbinRuleRow struct {
	ID    uint   `gorm:"primaryKey;autoIncrement"`
	Ptype string `gorm:"size:100"`
	V0    string `gorm:"size:100"`
	V1    string `gorm:"size:100"`
	V2    string `gorm:"size:100"`
	V3    string `gorm:"size:100"`
	V4    string `gorm:"size:100"`
	V5    string `gorm:"size:100"`
}

func (casbinRuleRow) TableName() string { return "casbin_rule" }

func grantOn(t *testing.T, db *gorm.DB) {
	t.Helper()
	tables := []any{&casbinRuleRow{}, &adminmodels.SysApi{}, &adminmodels.SysMenu{}, &adminmodels.SysRole{}}
	drop := func() {
		for _, join := range []string{"sys_menu_api_rule", "sys_role_menu", "sys_role_dept"} {
			if db.Migrator().HasTable(join) {
				if err := db.Migrator().DropTable(join); err != nil {
					t.Fatal(err)
				}
			}
		}
		for _, m := range tables {
			if db.Migrator().HasTable(m) {
				if err := db.Migrator().DropTable(m); err != nil {
					t.Fatal(err)
				}
			}
		}
	}
	drop()
	t.Cleanup(drop)
	if err := db.AutoMigrate(tables...); err != nil {
		t.Fatal(err)
	}
	for _, m := range []adminmodels.SysMenu{
		{MenuName: "Gen", Component: "/dev-tools/gen/index", MenuType: "C"},
		{MenuName: "EditTable", Component: "/dev-tools/gen/editTable", MenuType: "C"},
	} {
		if err := db.Create(&m).Error; err != nil {
			t.Fatal(err)
		}
	}
	var gen adminmodels.SysMenu
	if err := db.Where("component = ?", "/dev-tools/gen/index").First(&gen).Error; err != nil {
		t.Fatal(err)
	}
	role := adminmodels.SysRole{RoleKey: "developer", RoleName: "developer"}
	if err := db.Create(&role).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Exec("INSERT INTO sys_role_menu (role_id, menu_id) VALUES (?, ?)", role.RoleId, gen.MenuId).Error; err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 2; i++ {
		if err := bindGenMenuApis(db); err != nil {
			t.Fatalf("run %d: %v", i+1, err)
		}
	}

	var n int64
	db.Model(&casbinRuleRow{}).Where("ptype = 'p' AND v0 = ?", "developer").Count(&n)
	if n != 9 {
		t.Errorf("developer holds %d policies after two runs, want 9", n)
	}
	if bound := db.Model(&gen).Association("SysApi").Count(); bound != 9 {
		t.Errorf("代码生成 is bound to %d APIs after two runs, want 9", bound)
	}
}

func TestGenMenuApisOnMySQL(t *testing.T)     { grantOn(t, mysqlDB(t)) }
func TestGenMenuApisOnPostgres(t *testing.T)  { grantOn(t, postgresDB(t)) }
func TestGenMenuApisOnSQLServer(t *testing.T) { grantOn(t, sqlserverDB(t)) }
