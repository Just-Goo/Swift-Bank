package main

import (
	"context"
	"log"
	"net"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zde37/Swift_Bank/api"
	"github.com/zde37/Swift_Bank/config"
	"github.com/zde37/Swift_Bank/database"
	"github.com/zde37/Swift_Bank/gapi"
	"github.com/zde37/Swift_Bank/pb"
	"github.com/zde37/Swift_Bank/repository"
	"github.com/zde37/Swift_Bank/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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

	// runGrpcServer(config, pool)

	runGinServer(config, pool)
}

func runGrpcServer(config config.Config, pool *pgxpool.Pool) {	
	repository := repository.NewRepository(pool)
	service := service.NewService(repository.R)
	server, err := gapi.NewServer(config, service.S)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSwiftBankServer(grpcServer, server)
	reflection.Register(grpcServer) // it allows the grpc client to explore what RPCs are available on the server and how to call them (i.e self documentation for the server)

	// create a new listener
	listener, err := net.Listen("tcp", config.GrpcServerAddress)
	if err != nil {
		log.Fatal("cannot create listener:", err)
	}

	log.Printf("start gRPC server at: %s", listener.Addr().String())
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal("cannot start gRPC server:", err)
	}
}

func runGinServer(config config.Config, pool *pgxpool.Pool) {
	repository := repository.NewRepository(pool)
	service := service.NewService(repository.R)
	handler, err := api.NewHandler(config, service.S)
	if err != nil {
		log.Fatal("failed to load handler:", err)
	}

	err = handler.H.StartServer(config.HttpServerAddress)
	if err != nil {
		log.Fatal("failed to start server:", err)
	}
}
