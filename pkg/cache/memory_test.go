package cache

import (
	"fmt"
	"testing"
	"time"

	"github.com/matchstalk/go-admin-core/cache"
)

func TestInitMemory(t *testing.T) {
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
			if err := InitMemory(); (err != nil) != tt.wantErr {
				t.Errorf("InitRedis() error = %v, wantErr %v", err, tt.wantErr)
			}
			MemoryAdapter.Set("test", "1", 100)
			key, _ := MemoryAdapter.Get("test")
			message := &cache.MemoryMessage{}
			message.Stream = "queuetest"
			message.Values = map[string]interface{}{
				"key": "value",
			}
			MemoryAdapter.Append("queuetest", message)
			MemoryAdapter.Register("queuetest", func(message cache.Message) error {
				fmt.Println(message.GetValues())
				return nil
			})
			go func() {
				MemoryAdapter.Run()
			}()
			time.Sleep(time.Second)
			MemoryAdapter.Shutdown()
			fmt.Println(key)
		})
	}
}
