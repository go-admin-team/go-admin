package middleware

import (
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"go-admin/common/config"
	"go-admin/common/global"
	"go-admin/common/middleware"
	"go-admin/tools"
)

var WithContextDb = middleware.WithContextDb

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
