package config

import "github.com/spf13/viper"

type Database struct {
	Dbtype   string
	Host     string
	Port     int
	Database string
	Username string
	Password string
}

func InitDatabase(cfg *viper.Viper) *Database {
	return &Database{
		Port:     cfg.GetInt("port"),
		Dbtype:   cfg.GetString("dbType"),
		Host:     cfg.GetString("host"),
		Database: cfg.GetString("database"),
		Username: cfg.GetString("username"),
		Password: cfg.GetString("password"),
	}
}

var DatabaseConfig = new(Database)
