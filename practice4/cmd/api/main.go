package main

import (
	"context"
	"log"
	"net/http"
	"prac4/internal/handler"
	"prac4/internal/middleware"
	"prac4/internal/repository"
	"prac4/internal/repository/_postgres"
	"prac4/internal/usecase"
	"prac4/pkg/modules"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbConfig := initPostgreConfig()
	_postgre := _postgres.NewPGXDialect(ctx, dbConfig)
	repositories := repository.NewRepositories(_postgre)

	userUsecase := usecase.NewUserUsecase(repositories.UserRepository)
	userHandler := handler.NewUserHandler(userUsecase)

	router := mux.NewRouter()

	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.Use(middleware.AuthMiddleware)
	apiRouter.Use(middleware.LoggingMiddleware)

	apiRouter.HandleFunc("/users", userHandler.GetUsers).Methods("GET")
	apiRouter.HandleFunc("/users/{id}", userHandler.GetUserByID).Methods("GET")
	apiRouter.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	apiRouter.HandleFunc("/users/{id}", userHandler.UpdateUser).Methods("PUT")
	apiRouter.HandleFunc("/users/{id}", userHandler.UpdateUser).Methods("PATCH")
	apiRouter.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods("DELETE")

	router.HandleFunc("/health", userHandler.HealthCheck).Methods("GET")
	router.Use(middleware.LoggingMiddleware)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Println("Server starting on :8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func initPostgreConfig() *modules.PostgreConfig {
	return &modules.PostgreConfig{
		Host:        "localhost",
		Port:        "5432",
		Username:    "postgres",
		Password:    "Dakobay1994",
		DBName:      "mydb",
		SSLMode:     "disable",
		ExecTimeout: 5 * time.Second,
	}
}
