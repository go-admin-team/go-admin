// +build sqlite3

package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"go-admin/global"
	"go-admin/tools/config"
)

type SqLite struct {
}

func (e *SqLite) Setup() {
	var err error
	var db Database

	db = new(SqLite)
	global.Source = db.GetConnect()
	log.Info(global.Source)
	global.Eloquent, err = db.Open(db.GetDriver(), db.GetConnect())

	if err != nil {
		log.Fatalf("%s connect error %v", db.GetDriver(), err)
	} else {
		log.Printf("%s connect success!", db.GetDriver())
	}

	if global.Eloquent.Error != nil {
		log.Fatalf("database error %v", global.Eloquent.Error)
	}

	global.Eloquent.LogMode(true)
}

// 打开数据库连接
func (*SqLite) Open(dbType string, conn string) (db *gorm.DB, err error) {
	eloquent, err := gorm.Open(dbType, conn)
	return eloquent, err
}

func (e *SqLite) GetConnect() string {
	return config.DatabaseConfig.Source
}

func (e *SqLite) GetDriver() string {
	return config.DatabaseConfig.Driver
}
