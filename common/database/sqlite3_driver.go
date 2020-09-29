// +build sqlite3

package database

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"go-admin/common/global"
	"go-admin/tools/config"
)

type SqLite struct {
}

func (e *SqLite) Setup() {
	var err error

	global.Source = e.GetConnect()
	log.Println(global.Source)
	global.Eloquent, err = e.Open(e.GetDriver(), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if err != nil {
		log.Fatalf("%s connect error %v", e.GetDriver(), err)
	} else {
		log.Printf("%s connect success!", e.GetDriver())
	}

	if global.Eloquent.Error != nil {
		log.Fatalf("database error %v", global.Eloquent.Error)
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
func (*SqLite) Open(conn string, cfg *gorm.Config) (db *gorm.DB, err error) {
	eloquent, err := gorm.Open(sqlite.Open(conn), cfg)
	return eloquent, err
}

func (e *SqLite) GetConnect() string {
	return config.DatabaseConfig.Source
}

func (e *SqLite) GetDriver() string {
	return config.DatabaseConfig.Driver
}
