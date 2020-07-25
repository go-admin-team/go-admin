package config

import "github.com/spf13/viper"

type Database struct {
	Driver string
	Source string
	DBName string
	Logger *Logger
}

func InitDatabase(cfg *viper.Viper) *Database {

	db := &Database{
		Driver: cfg.GetString("driver"),
		Source: cfg.GetString("source"),
		DBName: cfg.GetString("dbname"),
		Logger: &Logger{
			Path:    cfg.GetString("logger.path"),
			Level:   cfg.GetString("logger.level"),
			Stdout:  cfg.GetBool("logger.stdout"),
			Enabled: cfg.GetBool("logger.enabled"),
		},
	}
	return db
}

var DatabaseConfig = new(Database)
