// +build sqlite3

package middleware

import (
	"database/sql"
	"errors"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"go-admin/common/global"
)

func getGormFromDb(driver string, db *sql.DB, config *gorm.Config) (*gorm.DB, error) {
	switch driver {
	case "mysql":
		return gorm.Open(mysql.New(mysql.Config{Conn: db}), config)
	case "postgres":
		return gorm.Open(postgres.New(postgres.Config{Conn: db}), config)
	case "sqlite3":
		return gorm.Open(sqlite.Open(global.Source), config)
	default:
		return nil, errors.New("not support this db driver")
	}
}
