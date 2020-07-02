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

func (*SqLite) Open(dbType string, conn string) (db *gorm.DB, err error) {
	eloquent, err := gorm.Open(dbType, conn)
	return eloquent, err
}

func (e *SqLite) GetConnect() string {
	return config.DatabaseConfig.SqLite.MasterConn
}

func (e *SqLite) Setup() {

	var err error
	var db Database

	db = new(SqLite)
	orm.SqLiteConn = db.GetConnect()
	log.Info(orm.SqLiteConn)
	orm.Eloquent, err = db.Open(config.DatabaseConfig.DbType, orm.SqLiteConn)

	if err != nil {
		log.Fatalf("%s connect error %v", config.DatabaseConfig.DbType, err)
	} else {
		log.Printf("%s connect success!", config.DatabaseConfig.DbType)
	}

	if orm.Eloquent.Error != nil {
		log.Fatalf("database error %v", orm.Eloquent.Error)
	}

	orm.Eloquent.LogMode(true)
}