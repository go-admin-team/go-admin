package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	log "github.com/sirupsen/logrus"
	"go-admin/global/orm"
)

type PgSql struct {
}

func (*PgSql) Open(dbType string, conn string) (db *gorm.DB, err error) {
	eloquent, err := gorm.Open(dbType, conn)
	return eloquent, err
}

func (e *PgSql) GetConnect() string {

	return ""
}

func (e *PgSql) Setup() {

	var err error
	var db Database

	db = new(PgSql)
	orm.PgSqlConn = db.GetConnect()
	log.Info(orm.PgSqlConn)
	orm.Eloquent, err = db.Open(DbType, orm.PgSqlConn)

	if err != nil {
		log.Fatalf("%s connect error %v", DbType, err)
	} else {
		log.Printf("%s connect success!", DbType)
	}

	if orm.Eloquent.Error != nil {
		log.Fatalf("database error %v", orm.Eloquent.Error)
	}

	orm.Eloquent.LogMode(true)
}