package config

import (
	"github.com/spf13/viper"
	"time"
)

type Ssl struct {
	KeyStr      string
	Pem         string
	Enable      bool
	Domain      string
	CertCache   string
	RenewBefore time.Duration
}

func InitSsl(cfg *viper.Viper) *Ssl {
	return &Ssl{
		KeyStr:      cfg.GetString("key"),
		Pem:         cfg.GetString("pem"),
		Enable:      cfg.GetBool("enable"),
		Domain:      cfg.GetString("domain"),
		CertCache:   cfg.GetString("cert"),
		RenewBefore: time.Duration(cfg.GetUint32("renew")),
	}
}

var SslConfig = new(Ssl)
