package middleware

import (
	"database/sql"
	"errors"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"go-admin/common/config"
	"go-admin/common/global"
	"go-admin/common/middleware"
	"go-admin/tools"
)

var WithContextDb = middleware.WithContextDb

func getGormFromDb(driver string, db *sql.DB, config *gorm.Config) (*gorm.DB, error) {
	switch driver {
	case "mysql":
		return gorm.Open(mysql.New(mysql.Config{Conn: db}), config)
	case "postgres":
		return gorm.Open(postgres.New(postgres.Config{Conn: db}), config)
	default:
		return nil, errors.New("not support this db driver")
	}
}

func GetGormFromConfig(cfg config.Conf) map[string]*gorm.DB {
	gormDB := make(map[string]*gorm.DB)
	if cfg.GetSaas() {
		var err error
		for k, v := range cfg.GetDbs() {
			gormDB[k], err = getGormFromDb(v.Driver, v.DB, &gorm.Config{
				NamingStrategy: schema.NamingStrategy{
					SingularTable: true,
				},
			})
			if err != nil {
				global.Logger.Fatal(tools.Red(k+" connect error :"), err)
			}
		}
		return gormDB
	}
	c := cfg.GetDb()
	db, err := getGormFromDb(c.Driver, c.DB, &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		global.Logger.Fatal(tools.Red(c.Driver+" connect error :"), err)
	}
	gormDB["*"] = db
	return gormDB
}
