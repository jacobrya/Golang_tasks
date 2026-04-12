package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func startServer(ctx context.Context, name string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Duration(rand.Intn(500)) * time.Millisecond):
				out <- fmt.Sprintf("[%s] metric: %d", name, rand.Intn(100))
			}
		}
	}()
	return out
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

	alphaStream := startServer(ctx, "Alpha")
	betaStream := startServer(ctx, "Beta")
	gammaStream := startServer(ctx, "Gamma")

	output := FanIn(ctx, alphaStream, betaStream, gammaStream)

	for m := range output {
		fmt.Println(m)
	}
}
