package cache

import (
	"fmt"
	"testing"
	"time"

	"github.com/matchstalk/go-admin-core/cache"
)

func TestInitRedis(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			"test01",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InitRedis(); (err != nil) != tt.wantErr {
				t.Errorf("InitRedis() error = %v, wantErr %v", err, tt.wantErr)
			}
			RedisAdapter.Set("test", "1", 100)
			key, _ := RedisAdapter.Get("test")
			message := &cache.RedisMessage{}
			message.Stream = "queuetest"
			message.Values = map[string]interface{}{
				"key": "value",
			}
			RedisAdapter.Append("queuetest", message)
			RedisAdapter.Register("queuetest", func(message cache.Message) error {
				fmt.Println(message.GetValues())
				return nil
			})
			go func() {
				RedisAdapter.Run()
			}()
			time.Sleep(time.Second)
			RedisAdapter.Shutdown()
			fmt.Println(key)
		})
	}
}
