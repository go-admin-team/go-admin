package config

import (
	"github.com/spf13/viper"
)

type Jwt struct {
	Secret  string
	Timeout int64
}

func InitJwt(cfg *viper.Viper) *Jwt {
	return &Jwt{
		Secret:  cfg.GetString("secret"),
		Timeout: cfg.GetInt64("timeout"),
	}
}

var JwtConfig = new(Jwt)
