package database

func Setup(driver string) {
	dbType := driver
	if dbType == "mysql" {
		var db = new(Mysql)
		db.Setup()
	}

	//TODO： 如果需要sqlite3请开启下面注释
	//if dbType == "sqlite3" {
	//	var db = new(SqLite)
	//	db.Setup()
	//}

	if dbType == "postgres" {
		var db = new(PgSql)
		db.Setup()
	}
}
