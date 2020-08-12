package cache

import (
	"github.com/go-redis/redis/v7"
	"github.com/matchstalk/redisqueue"
	"github.com/matchstalk/utils/cache"
	"time"
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
