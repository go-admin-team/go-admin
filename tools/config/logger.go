package config

import "github.com/spf13/viper"

type Logger struct {
	Path       string
	Level      string
	Stdout     bool
	EnabledBUS bool
	EnabledREQ bool
	EnabledDB  bool
	EnabledJOB bool `default:"false"`
}

func InitLog(cfg *viper.Viper) *Logger {
	return &Logger{
		Path:       cfg.GetString("path"),
		Level:      cfg.GetString("level"),
		Stdout:     cfg.GetBool("stdout"),
		EnabledBUS: cfg.GetBool("enabledbus"),
		EnabledREQ: cfg.GetBool("enabledreq"),
		EnabledDB:  cfg.GetBool("enableddb"),
		EnabledJOB: cfg.GetBool("enabledjob"),
	}
}

var LoggerConfig = new(Logger)
