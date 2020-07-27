package config

import "github.com/spf13/viper"

type Database struct {
	Driver string
	Source string
}

func InitDatabase(cfg *viper.Viper) *Database {

	db := &Database{
		Driver: cfg.GetString("driver"),
		Source: cfg.GetString("source"),
	}
	return db
}

var DatabaseConfig = new(Database)
