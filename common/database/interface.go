package database

import "gorm.io/gorm"

// Database 数据库配置
type Database interface {
	Setup()
	Open(conn string, cfg *gorm.Config) (db *gorm.DB, err error)
	GetConnect() string
	GetDriver() string
}
