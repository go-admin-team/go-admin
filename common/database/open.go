// +build !sqlite3

package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var opens = map[string]func(string) gorm.Dialector{
	"mysql":    mysql.Open,
	"postgres": postgres.Open,
}
