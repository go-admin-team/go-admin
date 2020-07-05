package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	log "github.com/sirupsen/logrus"
	"go-admin/global/orm"
	"go-admin/tools/config"
)

type SqLite struct {
}

func (e *SqLite) Setup() {

	var err error
	var db Database

	db = new(SqLite)
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
func (*SqLite) Open(dbType string, conn string) (db *gorm.DB, err error) {
	eloquent, err := gorm.Open(dbType, conn)
	return eloquent, err
}

func (e *SqLite) GetConnect() string {
	return config.DatabaseConfig.Source
}