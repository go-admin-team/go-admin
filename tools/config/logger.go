package config

import "github.com/spf13/viper"

type Logger struct {
	Path    string
	Level   string
	Stdout  bool
	Enabled bool
}

func InitLog(cfg *viper.Viper) *Logger {
	return &Logger{
		Path:    cfg.GetString("path"),
		Level:   cfg.GetString("level"),
		Enabled: cfg.GetBool("enabled"),
		Stdout:  cfg.GetBool("stdout"),
	}
}

var LoggerConfig = new(Logger)
