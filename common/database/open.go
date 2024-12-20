//go:build !sqlite3

package database

import (
	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var opens = map[string]func(string) gorm.Dialector{
	"mysql":     mysql.Open,
	"postgres":  postgres.Open,
	"sqlserver": sqlserver.Open,
	"sqllite3":  sqlite.Open,
}
