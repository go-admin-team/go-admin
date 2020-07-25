package database

import (
	_ "github.com/go-sql-driver/mysql" //加载mysql
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"go-admin/global"
	"go-admin/tools/config"
)

type Mysql struct {
}

func (e *Mysql) Setup() {
	var err error
	var db Database

	db = new(Mysql)
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

	global.Eloquent.LogMode(config.DatabaseConfig.Logger.Enabled)
	//global.Eloquent.SetLogger(global.DBLogger)
}

// 打开数据库连接
func (e *Mysql) Open(dbType string, conn string) (db *gorm.DB, err error) {
	return gorm.Open(dbType, conn)
}

// 获取数据库连接
func (e *Mysql) GetConnect() string {
	return config.DatabaseConfig.Source
}

func (e *Mysql) GetDriver() string {
	return config.DatabaseConfig.Driver
}
