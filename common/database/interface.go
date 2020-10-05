package database

import "gorm.io/gorm"

type Database interface {
	Setup()
	Open(conn string, cfg *gorm.Config) (db *gorm.DB, err error)
	GetConnect() string
	GetDriver() string
}
