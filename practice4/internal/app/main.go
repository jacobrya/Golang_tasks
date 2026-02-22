package main

import (
	"context"
	"fmt"
	"prac4/internal/repository"
	"prac4/internal/repository/_postgres"
	"prac4/pkg/modules"
	"time"
)

func Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	dbConfig := initPostgreConfig()
	_postgre := _postgres.NewPGXDialect(ctx, dbConfig)
	repositories := repository.NewRepositories(_postgre)
	users, err := repositories.GetUsers()
	if err != nil {
		fmt.Printf("Error fetching users: %v\n", err)
		return
	}
	fmt.Printf("Users: %+v\n", users)
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

func main() {
	Run()
}
