package config

import "github.com/spf13/viper"

type Ssl struct {
	KeyStr string
	Pem    string
	Enable bool
	Domain string
}

func InitSsl(cfg *viper.Viper) *Ssl {
	return &Ssl{
		KeyStr: cfg.GetString("key"),
		Pem:    cfg.GetString("pem"),
		Enable: cfg.GetBool("enable"),
		Domain: cfg.GetString("domain"),
	}
}

var SslConfig = new(Ssl)
