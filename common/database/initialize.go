package database

import (
	. "log"
	"time"

	logCore "github.com/go-admin-team/go-admin-core/logger"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"go-admin/common/global"
	"go-admin/common/log"
	mycasbin "go-admin/pkg/casbin"
	"go-admin/tools"
	toolsConfig "go-admin/tools/config"
)

// Setup 配置数据库
func Setup() {
	for k := range toolsConfig.DatabasesConfig {
		setupSimpleDatabase(k, toolsConfig.DatabasesConfig[k])
	}
}

func setupSimpleDatabase(host string, c *toolsConfig.Database) {
	if global.Driver == "" {
		global.Driver = c.Driver
	}
	log.Infof("%s => %s", host, tools.Green(c.Source))
	db, err := gorm.Open(open[c.Driver](c.Source), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: logger.New(
			New(logCore.DefaultLogger.Options().Out, "\r\n", LstdFlags),
			logger.Config{
				SlowThreshold: time.Second,
				Colorful:      true,
				LogLevel: logger.LogLevel(
					logCore.DefaultLogger.Options().Level.LevelForGorm()),
			},
		),
	})
	if err != nil {
		log.Fatal(tools.Red(c.Driver+" connect error :"), err)
	} else {
		log.Info(tools.Green(c.Driver + " connect success !"))
	}

	e := mycasbin.Setup(db, "sys_")

	if host == "*" {
		global.Eloquent = db
	}

	global.Cfg.SetDb(host, db)
	global.Cfg.SetCasbin(host, e)
}
