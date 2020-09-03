package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"

	"go-admin/common/global"
	"go-admin/tools"
	"go-admin/tools/config"
)

type Mysql struct {
}

func (e *Mysql) Setup() {
	var err error

	global.Source = e.GetConnect()
	global.Logger.Info(tools.Green(global.Source))
	global.Eloquent, err = e.Open(e.GetConnect(), &gorm.Config{
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

	if config.LoggerConfig.EnabledDB {
		global.Eloquent.Logger = logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold: time.Second,
			Colorful:      true,
			LogLevel:      logger.Info,
		})
	}
}

// 打开数据库连接
func (e *Mysql) Open(conn string, cfg *gorm.Config) (db *gorm.DB, err error) {
	return gorm.Open(mysql.Open(conn), cfg)
}

// 获取数据库连接
func (e *Mysql) GetConnect() string {
	return config.DatabaseConfig.Source
}

func (e *Mysql) GetDriver() string {
	return config.DatabaseConfig.Driver
}
