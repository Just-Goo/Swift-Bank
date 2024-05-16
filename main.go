package main

import (
	"context"
	"log"

	"github.com/zde37/Swift_Bank/api"
	"github.com/zde37/Swift_Bank/config"
	"github.com/zde37/Swift_Bank/database"
	"github.com/zde37/Swift_Bank/repository"
	"github.com/zde37/Swift_Bank/service"
)

func main() {

	// load config
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}

	// Load database
	ctx := context.Background()
	PostgresClient := database.PostgresClient{}
	pool, err := PostgresClient.NewPostgresClient(ctx, config.Dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err = PostgresClient.PingDB(ctx); err != nil {
		log.Fatal(err)
	}

	defer pool.Close()

	repository := repository.NewRepository(pool)
	service := service.NewService(repository.R)
	handler, err := api.NewHandler(*config, service.S)
	if err != nil {
		log.Fatal("failed to load handler:", err)
	}

	err = handler.H.StartServer(config.Port)
	if err != nil {
		log.Fatal("failed to start server:", err)
	}
}
