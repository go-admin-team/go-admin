package database

import (
	"bytes"
	"fmt"
	_ "github.com/go-sql-driver/mysql" //加载mysql
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	config2 "go-admin/config"
	"log"
	"strconv"
)

var Eloquent *gorm.DB

func init() {

	dbType := config2.DatabaseConfig.Dbtype
	host := config2.DatabaseConfig.Host
	port := config2.DatabaseConfig.Port
	database := config2.DatabaseConfig.Database
	username := config2.DatabaseConfig.Username
	password := config2.DatabaseConfig.Password

	if dbType != "mysql" && dbType != "sqlite3" {
		fmt.Println("db type unknow")
	}
	var err error

	var conn bytes.Buffer
	conn.WriteString(username)
	conn.WriteString(":")
	conn.WriteString(password)
	conn.WriteString("@tcp(")
	conn.WriteString(host)
	conn.WriteString(":")
	conn.WriteString(strconv.Itoa(port))
	conn.WriteString(")")
	conn.WriteString("/")
	conn.WriteString(database)
	//conn.WriteString("?charset=utf8&parseTime=True&loc=Local&timeout=1000ms")

	log.Println(conn.String())

	var db Database
	if dbType == "mysql" {
		db = new(Mysql)
		Eloquent, err = db.Open(dbType, conn.String())

	} else if dbType == "sqlite3" {
		db = new(SqlLite)
		Eloquent, err = db.Open(dbType, host)

	} else {
		panic("db type unknow")
	}

	Eloquent.LogMode(true)
	if err != nil {
		log.Fatalln("%s connect error %v",dbType, err)
	} else {
		log.Println("%s connect success!",dbType)
	}

	if Eloquent.Error != nil {
		log.Fatalln("database error %v", Eloquent.Error)
	}

}

type Database interface {
	Open(dbType string, conn string) (db *gorm.DB, err error)
}

type Mysql struct {
}

func (*Mysql) Open(dbType string, conn string) (db *gorm.DB, err error) {
	eloquent, err := gorm.Open(dbType, conn)
	return eloquent, err
}

type SqlLite struct {
}

func (*SqlLite) Open(dbType string, conn string) (db *gorm.DB, err error) {
	eloquent, err := gorm.Open(dbType, conn)
	return eloquent, err
}
