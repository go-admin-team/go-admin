// +build !sqlite3

package database

// Setup 配置数据库
func Setup(driver string) {
	dbType := driver
	if dbType == "mysql" {
		var db = new(Mysql)
		db.Setup()
	}

	if dbType == "postgres" {
		var db = new(PgSql)
		db.Setup()
	}
}
