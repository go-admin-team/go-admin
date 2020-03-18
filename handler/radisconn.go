package handler

import (
	"github.com/go-redis/redis"
	"go-admin/config"
)

var RedisClient = new(redis.Client)

func init() {
	RedisNewClient(config.RedisConnConfig.Addr, config.RedisConnConfig.Password, config.RedisConnConfig.DB)
}

func RedisNewClient(addr string, password string, db int) {
	//timeout := time.Duration(readTimeout)
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // no password set
		DB:       db,       // use default DB
		//ReadTimeout: ,
	})
}
