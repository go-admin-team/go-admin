// +build sqlite3

package database

func Setup(driver string) {
	dbType := driver
	if dbType == "mysql" {
		var db = new(Mysql)
		db.Setup()
	}

	if dbType == "sqlite3" {
		var db = new(SqLite)
		db.Setup()
	}

	if dbType == "postgres" {
		var db = new(PgSql)
		db.Setup()
	}
}
