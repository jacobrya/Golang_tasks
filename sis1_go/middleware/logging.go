package middleware

import (
	"log"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		startTime := time.Now()

		next(w, r)

		log.Printf("%s %s %s Request processed",
			startTime.Format("2006-01-02T15:04:05"),
			r.Method,
			r.URL.Path,
		)
	}
}
