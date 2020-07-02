package database

import "go-admin/tools/config"


func Setup() {
	dbType := config.DatabaseConfig.DbType
	if dbType == "mysql" {
		var db = new(Mysql)
		db.Setup()
	}

	if dbType == "sqlite3" {
		var db = new(SqLite)
		db.Setup()
	}

	if dbType == "pgsql" {
		var db = new(PgSql)
		db.Setup()
	}
}
