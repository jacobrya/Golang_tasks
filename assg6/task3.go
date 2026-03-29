package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func launchProducer(ctx context.Context, label string) <-chan string {
	msgChan := make(chan string)

	go func() {
		defer close(msgChan)

		for {
			delay := time.Duration(rand.Intn(500)) * time.Millisecond

			select {
			case <-ctx.Done():
				return
			case <-time.After(delay):
				msg := fmt.Sprintf("[%s] metric: %d", label, rand.Intn(100))
				msgChan <- msg
			}
		}
	}()

	return msgChan
}

func FanIn(ctx context.Context, sources ...<-chan string) <-chan string {
	combined := make(chan string)
	wg := new(sync.WaitGroup)
	wg.Add(len(sources))

	for _, src := range sources {
		go func(input <-chan string) {
			defer wg.Done()

			for v := range input {
				select {
				case combined <- v:
				case <-ctx.Done():
					return
				}
			}
		}(src)
	}

	go func() {
		wg.Wait()
		close(combined)
	}()

	return combined
}

func main() {
	timeout := 2 * time.Second
	ctx, stop := context.WithTimeout(context.Background(), timeout)
	defer stop()

	alphaStream := launchProducer(ctx, "Alpha")
	betaStream := launchProducer(ctx, "Beta")
	gammaStream := launchProducer(ctx, "Gamma")

	output := FanIn(ctx, alphaStream, betaStream, gammaStream)

	for m := range output {
		fmt.Println(m)
	}
}
