package database

import (
	"database/sql"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"go-admin/common/config"
	"go-admin/common/global"
	goAdminLogger "go-admin/logger"
	"go-admin/tools"
	toolsConfig "go-admin/tools/config"
)

type Mysql struct {
}

func (e *Mysql) Setup() {
	global.Source = e.GetConnect()
	global.Logger.Info(tools.Green(global.Source))
	db, err := sql.Open("mysql", global.Source)
	if err != nil {
		global.Logger.Fatal(tools.Red(e.GetDriver()+" connect error :"), err)
	}
	global.Cfg.SetDb(&config.DBConfig{
		Driver: "mysql",
		DB:     db,
	})
	global.Eloquent, err = e.Open(db, &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		global.Logger.Fatal(tools.Red(e.GetDriver()+" connect error :"), err)
	} else {
		global.Logger.Info(tools.Green(e.GetDriver() + " connect success !"))
	}

	if global.Eloquent.Error != nil {
		global.Logger.Fatal(tools.Red(" database error :"), global.Eloquent.Error)
	}

	if toolsConfig.LoggerConfig.EnabledDB {
		global.Eloquent.Logger = logger.New(
			log.New(goAdminLogger.DefaultLogger.Options().Out, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: time.Second,
				Colorful:      true,
				LogLevel: logger.LogLevel(
					goAdminLogger.DefaultLogger.Options().Level.LevelForGorm()),
			})
	}
}

// 打开数据库连接
func (e *Mysql) Open(db *sql.DB, cfg *gorm.Config) (*gorm.DB, error) {
	return gorm.Open(mysql.New(mysql.Config{Conn: db}), cfg)
}

// 获取数据库连接
func (e *Mysql) GetConnect() string {
	return toolsConfig.DatabaseConfig.Source
}

func (e *Mysql) GetDriver() string {
	return toolsConfig.DatabaseConfig.Driver
}
