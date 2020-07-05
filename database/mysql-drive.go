package database

import (
	_ "github.com/go-sql-driver/mysql" //加载mysql
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"go-admin/global/orm"
	"go-admin/tools/config"
)

type Mysql struct {
}

func (e *Mysql) Setup() {
	var err error
	var db Database
	db = new(Mysql)
	orm.Source = db.GetConnect()
	log.Info(orm.Source)
	orm.Eloquent, err = db.Open(orm.Driver, orm.Source)
	if err != nil {
		log.Fatalf("%s connect error %v", orm.Driver, err)
	} else {
		log.Printf("%s connect success!", orm.Driver)
	}
	if orm.Eloquent.Error != nil {
		log.Fatalf("database error %v", orm.Eloquent.Error)
	}
	orm.Eloquent.LogMode(true)
}

// 打开数据库连接
func (e *Mysql) Open(dbType string, conn string) (db *gorm.DB, err error) {
	return gorm.Open(dbType, conn)
}

// 获取数据库连接
func (e *Mysql) GetConnect() string {
	return config.DatabaseConfig.Source
}
