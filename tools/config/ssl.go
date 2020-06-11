package config

import "github.com/spf13/viper"

type Ssl struct {
	KeyStr string
	Pem    string
}

func InitSsl(cfg *viper.Viper) *Ssl {
	return &Ssl{
		KeyStr: cfg.GetString("key"),
		Pem:    cfg.GetString("pem"),
	}
}

var SslConfig = new(Ssl)
