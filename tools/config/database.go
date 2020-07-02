package config

import "github.com/spf13/viper"

type Database struct {
	DbType string
	SqLite *SqLite
	Mysql  *Mysql
	PgSql  *PgSql
}

type SqLite struct {
	MasterConn string
	DBName     string
}

type Mysql struct {
	MasterConn string
	DBName     string
}

type PgSql struct {
	MasterConn string
	DBName     string
}

func InitDatabase(cfg *viper.Viper) *Database {

	dbType := cfg.GetString("dbType")

	db := &Database{
		DbType: cfg.GetString("dbType"),
	}

	if dbType == "sqlite3" {
		sqlit := cfg.Sub("sqlite")
		if sqlit == nil {
			panic("config not found settings.database.sqlite")
		}
		db.SqLite = InitSqlite(sqlit)
	} else if dbType == "mysql" {
		mysql := cfg.Sub("mysql")
		if mysql == nil {
			panic("config not found settings.database.mysql")
		}
		db.Mysql = InitMysql(mysql)
	} else if dbType == "pgsql" {
		pgsql := cfg.Sub("pgsql")
		if pgsql == nil {
			panic("config not found settings.database.pgsql")
		}
		db.PgSql = InitPgsql(pgsql)
	} else {
		panic("unknown dbtype")
	}

	return db
}

var DatabaseConfig = new(Database)

func InitSqlite(cfg *viper.Viper) *SqLite {
	return &SqLite{
		MasterConn: cfg.GetString("masterconn"),
		DBName:     cfg.GetString("dbname"),
	}
}

func InitMysql(cfg *viper.Viper) *Mysql {
	return &Mysql{
		MasterConn: cfg.GetString("masterconn"),
		DBName:     cfg.GetString("dbname"),
	}
}

func InitPgsql(cfg *viper.Viper) *PgSql {
	return &PgSql{
		MasterConn: cfg.GetString("masterconn"),
		DBName:     cfg.GetString("dbname"),
	}
}
