package repository

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/Just-Goo/Swift_Bank/config"
	"github.com/Just-Goo/Swift_Bank/database" 
)

var testRepo *Repository

func TestMain(m *testing.M) {
	ctx := context.Background()

	// create a new postgres instance
	config, err := config.LoadConfig("../../Swift Bank")
	if err != nil {
		log.Fatal("failed to load config: ", err)
	}

	postgresClient := database.PostgresClient{}
	db, err := postgresClient.NewPostgresClient(ctx, config.Dsn)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	testRepo = NewRepository(db)

	defer db.Close()

	if err = postgresClient.PingDB(ctx); err != nil {
		log.Fatal("cannot ping db: ", err)
	}

	os.Exit(m.Run())
}
