package cache

import (
	"time"

	"github.com/go-admin-team/go-admin-core/cache"
	"github.com/go-redis/redis/v7"
	"github.com/robinjoseph08/redisqueue/v2"
)

var RedisAdapter Adapter

func InitRedis() error {
	RedisAdapter = &cache.Redis{
		ConnectOption: &redis.Options{
			Addr: "127.0.0.1:6379",
		},
		ConsumerOptions: &redisqueue.ConsumerOptions{
			VisibilityTimeout: 60 * time.Second,
			BlockingTimeout:   5 * time.Second,
			ReclaimInterval:   1 * time.Second,
			BufferSize:        100,
			Concurrency:       10,
		},
		ProducerOptions: &redisqueue.ProducerOptions{
			StreamMaxLength:      100,
			ApproximateMaxLength: true,
		},
	}
	err := RedisAdapter.Connect()
	return err
}
