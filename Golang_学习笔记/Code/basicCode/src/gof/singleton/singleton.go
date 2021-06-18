package singleton

import (
	"sync"
	"sync/atomic"
)

/**
单例模式
 */

var (
	instance *singleton
	initialized uint32
	mu sync.Mutex
)

type singleton struct{
	name string
}

func GetSingletonInstance(name string) *singleton{
	if atomic.LoadUint32(&initialized) == 1{
		return instance
	}

	mu.Lock()
	defer mu.Unlock()

	if instance == nil{
		instance = &singleton{name: name}
		defer atomic.StoreUint32(&initialized, 1)
	}
	return instance
}
