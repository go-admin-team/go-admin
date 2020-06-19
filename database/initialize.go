package database

import "go-admin/tools/config"

func Setup() {
	dbType := config.DatabaseConfig.Dbtype
	if dbType == "mysql" {
		var db = new(Mysql)
		db.Setup()
	}
}
