package config

import "github.com/spf13/viper"

type SSL struct {
	Enable   bool
	CacheDir string
	Domain   string
}

func InitSSL(cfg *viper.Viper) *SSL {
	return &SSL{
		Enable:   cfg.GetBool("enable"),
		CacheDir: cfg.GetString("cache_dir"),
		Domain:   cfg.GetString("domain"),
	}
}

var SSLConfig = new(SSL)
