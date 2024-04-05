package main

import (
	"context"
	"log"

	"github.com/Just-Goo/Swift_Bank/config"
	"github.com/Just-Goo/Swift_Bank/controller"
	"github.com/Just-Goo/Swift_Bank/database" 
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

	controller.InitRouter(pool, config)
}
