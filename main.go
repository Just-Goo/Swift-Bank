package main

import (
	"context"
	"log"
	"net"
	"net/http"

	_ "github.com/zde37/Swift_Bank/doc/statik"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rakyll/statik/fs"
	"github.com/zde37/Swift_Bank/api"
	"github.com/zde37/Swift_Bank/config"
	"github.com/zde37/Swift_Bank/database"
	"github.com/zde37/Swift_Bank/gapi"
	"github.com/zde37/Swift_Bank/pb"
	"github.com/zde37/Swift_Bank/repository"
	"github.com/zde37/Swift_Bank/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
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

	go runGatewayServer(config, pool)
	runGrpcServer(config, pool)

	// runGinServer(config, pool)
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

func runGatewayServer(config config.Config, pool *pgxpool.Pool) {
	repository := repository.NewRepository(pool)
	service := service.NewService(repository.R)
	server, err := gapi.NewServer(config, service.S)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	// set json response to use snake case
	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})
	grpcMux := runtime.NewServeMux(jsonOption)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := pb.RegisterSwiftBankHandlerServer(ctx, grpcMux, server); err != nil {
		log.Fatal("cannot register handler server", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	// serve swagger file with statik
	statikFs, err := fs.New()
	if err != nil {
		log.Fatal("cannot create statik fs:", err)
	}

	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFs))
	mux.Handle("/swagger/", swaggerHandler)

	// create a new listener
	listener, err := net.Listen("tcp", config.HttpServerAddress)
	if err != nil {
		log.Fatal("cannot create listener:", err)
	}

	log.Printf("start HTTP gateway server at: %s", listener.Addr().String())
	if err := http.Serve(listener, mux); err != nil {
		log.Fatal("cannot start HTTP gateway server:", err)
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
