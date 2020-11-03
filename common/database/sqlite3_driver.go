// +build sqlite3

package database

import (
	"database/sql"
	. "log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"go-admin/common/config"
	"go-admin/common/global"
	"go-admin/common/log"
	"go-admin/tools"
	toolsConfig "go-admin/tools/config"
)

type SqLite struct {
}

func (e *SqLite) Setup() {
	var err error

	global.Source = e.GetConnect()
	log.Info(global.Source)
	db, err := sql.Open("sqlite3", global.Source)
	if err != nil {
		global.Logger.Fatal(tools.Red(e.GetDriver()+" connect error :"), err)
	}
	global.Cfg.SetDb(&config.DBConfig{
		Driver: "sqlite3",
		DB:     db,
	})
	global.Eloquent, err = e.Open(e.GetConnect(), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if err != nil {
		log.Fatalf("%s connect error %v", e.GetDriver(), err)
	} else {
		log.Infof("%s connect success!", e.GetDriver())
	}

	if global.Eloquent.Error != nil {
		log.Fatalf("database error %v", global.Eloquent.Error)
	}

	if toolsConfig.LoggerConfig.EnabledDB {
		global.Eloquent.Logger = logger.New(
			New(os.Stdout, "\r\n", LstdFlags), logger.Config{
				SlowThreshold: time.Second,
				Colorful:      true,
				LogLevel:      logger.Info,
			},
		)
	}
}

// 打开数据库连接
func (*SqLite) Open(conn string, cfg *gorm.Config) (db *gorm.DB, err error) {
	return gorm.Open(sqlite.Open(conn), cfg)
}

func (e *SqLite) GetConnect() string {
	return toolsConfig.DatabaseConfig.Source
}

func (e *SqLite) GetDriver() string {
	return toolsConfig.DatabaseConfig.Driver
}
