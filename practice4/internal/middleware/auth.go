package middleware

import (
	"log"
	"net/http"
	"os"
)

const validAPIKey = "123"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-KEY")

		if apiKey == "" {
			log.Printf("Unauthorized: Missing X-API-KEY header")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error":"X-API-KEY header is required"}`))
			return
		}

		expectedKey := os.Getenv("API_KEY")
		if expectedKey == "" {
			expectedKey = validAPIKey
		}

		if apiKey != expectedKey {
			log.Printf("Unauthorized: Invalid X-API-KEY header")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error":"Invalid X-API-KEY"}`))
			return
		}

		next.ServeHTTP(w, r)
	})
}

