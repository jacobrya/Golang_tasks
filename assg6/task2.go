package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

const totalIterations = 1000

func mutexBasedCounter() {
	cnt := 0
	lock := sync.Mutex{}
	wg := new(sync.WaitGroup)

	wg.Add(totalIterations)
	for idx := 0; idx < totalIterations; idx++ {
		go func() {
			lock.Lock()
			cnt++
			lock.Unlock()
			wg.Done()
		}()
	}

	wg.Wait()
	fmt.Println("[Mutex] Counter:", cnt)
}

func atomicBasedCounter() {
	var cnt int64
	wg := &sync.WaitGroup{}

	wg.Add(totalIterations)
	for idx := 0; idx < totalIterations; idx++ {
		go func() {
			atomic.AddInt64(&cnt, 1)
			wg.Done()
		}()
	}

	wg.Wait()
	fmt.Println("[Atomic] Counter:", cnt)
}

func main() {
	mutexBasedCounter()
	atomicBasedCounter()
}
