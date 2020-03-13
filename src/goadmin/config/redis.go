package config

import "github.com/spf13/viper"

type RedisConn struct {
	Addr     string
	Password string
	DB       int
}

func InitRedisConn(cfg *viper.Viper) *RedisConn {
	return &RedisConn{
		Addr:     cfg.GetString("addr"),
		Password: cfg.GetString("password"),
		DB:       cfg.GetInt("db"),
	}
}

var RedisConnConfig = new(RedisConn)
