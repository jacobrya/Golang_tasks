package main

import (
	"context"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"time"
)

func IsRetryable(resp *http.Response, err error) bool {
	if err != nil {
		return true
	}

	switch resp.StatusCode {
	case 429, 500, 502, 503, 504:
		return true
	default:
		return false
	}
}

func CalculateBackoff(attempt int) time.Duration {
	baseDelay := 500 * time.Millisecond
	maxDelay := 5 * time.Second

	backoff := float64(baseDelay) * math.Pow(2, float64(attempt))
	if backoff > float64(maxDelay) {
		backoff = float64(maxDelay)
	}

	return time.Duration(rand.Int63n(int64(backoff)))
}

func ExecutePayment(ctx context.Context, url string, maxRetries int) error {
	client := &http.Client{}

	for attempt := 0; attempt < maxRetries; attempt++ {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(`{"amount": 1000}`))
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		resp, err := client.Do(req)

		if err == nil {
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				body, _ := io.ReadAll(resp.Body)
				fmt.Printf("Attempt %d: Success! Response: %s\n", attempt+1, string(body))
				return nil
			}
		}

		if !IsRetryable(resp, err) {
			if resp != nil {
				return fmt.Errorf("non-retryable error: status %d", resp.StatusCode)
			}
			return fmt.Errorf("non-retryable error: %w", err)
		}

		if attempt == maxRetries-1 {
			return fmt.Errorf("max retries reached")
		}

		wait := CalculateBackoff(attempt)
		fmt.Printf("Attempt %d failed: waiting %v...\n", attempt+1, wait)

		select {
		case <-time.After(wait):
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}

func main() {
	var mu sync.Mutex
	count := 0

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		count++
		current := count
		mu.Unlock()

		if current <= 3 {
			fmt.Printf("[Server] Request #%d → 503 Service Unavailable\n", current)
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		fmt.Printf("[Server] Request #%d → 200 OK\n", current)
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `{"status": "success"}`)
	}))
	defer ts.Close()

	fmt.Println("=== Starting Payment Execution ===")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := ExecutePayment(ctx, ts.URL, 5)
	if err != nil {
		fmt.Printf("Payment failed: %v\n", err)
	}
}
