package main

import (
	"fmt"
	"sync"
)

const numWorkers = 100

func main() {
	mp := &sync.Map{}
	wg := new(sync.WaitGroup)

	wg.Add(numWorkers)
	for n := 0; n < numWorkers; n++ {
		go func(val int) {
			mp.Store("key", val)
			wg.Done()
		}(n)
	}

	wg.Wait()

	if result, found := mp.Load("key"); found {
		fmt.Printf("Value: %d\n", result)
	}
}
