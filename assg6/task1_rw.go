package main

import (
	"fmt"
	"sync"
)

const goroutineCount = 100

type ConcurrentMap struct {
	lock sync.RWMutex
	data map[string]int
}

func CreateConcurrentMap() *ConcurrentMap {
	return &ConcurrentMap{
		data: make(map[string]int),
	}
}

func (cm *ConcurrentMap) Set(k string, v int) {
	cm.lock.Lock()
	cm.data[k] = v
	cm.lock.Unlock()
}

func (cm *ConcurrentMap) Get(k string) (int, bool) {
	cm.lock.RLock()
	result, exists := cm.data[k]
	cm.lock.RUnlock()
	return result, exists
}

func main() {
	storage := CreateConcurrentMap()
	wg := &sync.WaitGroup{}

	wg.Add(goroutineCount)
	for n := 0; n < goroutineCount; n++ {
		go func(val int) {
			storage.Set("key", val)
			wg.Done()
		}(n)
	}

	wg.Wait()

	if result, exists := storage.Get("key"); exists {
		fmt.Printf("Value: %d\n", result)
	}
}
