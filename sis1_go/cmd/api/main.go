package main

import (
	"log"
	"net/http"
	"sis1_go/handlers"
	"sis1_go/middleware"
)

func main() {

	handler := middleware.LoggingMiddleware(
		middleware.APIKeyMiddleware(
			handlers.HandleTasks,
		),
	)

	http.HandleFunc("/tasks", handler)

	log.Println("Server starting on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
