package middleware

import (
	"log"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapped, r)

		timestamp := time.Now().Format(time.RFC3339)
		method := r.Method
		endpoint := r.URL.Path
		statusCode := wrapped.statusCode
		duration := time.Since(start)

		log.Printf("[%s] %s %s - Status: %d - Duration: %v", timestamp, method, endpoint, statusCode, duration)
	})
}

