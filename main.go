package main

import (
	"context"
	"log"

	"github.com/Just-Goo/Swift_Bank/config"
	"github.com/Just-Goo/Swift_Bank/controller/handler"
	"github.com/Just-Goo/Swift_Bank/controller/routes"
	"github.com/Just-Goo/Swift_Bank/database"
	"github.com/Just-Goo/Swift_Bank/repository"
	"github.com/Just-Goo/Swift_Bank/service"
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
	handler := handler.NewHandler(service.S)
	routes.RegisterRoutes(handler.H.GetGin(), handler)

	err = handler.H.StartServer(config.Port)
	if err != nil {
		log.Fatal("failed to start server:", err)
	}
}
