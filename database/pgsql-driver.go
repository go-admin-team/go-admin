package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	log "github.com/sirupsen/logrus"
	"go-admin/global/orm"
	"go-admin/tools/config"
)

type PgSql struct {
}

func (e *PgSql) Setup() {

	var err error
	var db Database

	db = new(PgSql)
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
func (*PgSql) Open(dbType string, conn string) (db *gorm.DB, err error) {
	eloquent, err := gorm.Open(dbType, conn)
	return eloquent, err
}

func (e *PgSql) GetConnect() string {

	return config.DatabaseConfig.Source
}
