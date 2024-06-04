package repository

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/zde37/Swift_Bank/database"
)

var testRepo *Repository

func TestMain(m *testing.M) {
	ctx := context.Background()

	// create a new postgres instance
	// config, err := config.LoadConfig("..")
	// if err != nil {
	// 	log.Fatal("failed to load config: ", err)
	// }
	dsn := "postgresql://root:4713a4cd628778cd1c37a95518f3eaf3@localhost:5432/Swift_Bank_DB?sslmode=disable"
	postgresClient := database.PostgresClient{}
	db, err := postgresClient.NewPostgresClient(ctx, dsn)
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
