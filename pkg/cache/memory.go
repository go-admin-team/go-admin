package cache

import (
	"github.com/go-admin-team/go-admin-core/cache"
)

var MemoryAdapter Adapter

func InitMemory() error {
	MemoryAdapter = &cache.Memory{
		PoolNum: 100,
	}
	err := MemoryAdapter.Connect()
	return err
}
