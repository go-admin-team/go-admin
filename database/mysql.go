package database

import (
	_ "github.com/go-sql-driver/mysql" //加载mysql
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"go-admin/global/orm"
	"go-admin/tools/config"
)

func (e *Mysql) Setup() {

	var err error
	var db Database

	db = new(Mysql)
	orm.MysqlConn = db.GetConnect()
	log.Info(orm.MysqlConn)
	orm.Eloquent, err = db.Open(config.DatabaseConfig.DbType, orm.MysqlConn)

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

type Mysql struct {
}

func (e *Mysql) Open(dbType string, conn string) (db *gorm.DB, err error) {
	return gorm.Open(dbType, conn)
}

func (e *Mysql) GetConnect() string {
	return config.DatabaseConfig.Mysql.MasterConn
}
