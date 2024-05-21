package main

import (
	"context"
	"net"
	"net/http"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rakyll/statik/fs"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/zde37/Swift_Bank/api"
	"github.com/zde37/Swift_Bank/config"
	"github.com/zde37/Swift_Bank/database"
	_ "github.com/zde37/Swift_Bank/doc/statik"
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
		log.Fatal().Err(err).Msg("failed to load config")
	}

	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	// Load database
	ctx := context.Background()
	PostgresClient := database.PostgresClient{}
	pool, err := PostgresClient.NewPostgresClient(ctx, config.Dsn)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create new postgres client")
	}

	defer pool.Close()

	if err = PostgresClient.PingDB(ctx); err != nil {
		log.Fatal().Err(err).Msg("failed to ping DB")
	}

	runDBMigration(config.MigrationURL, config.Dsn)

	go runGatewayServer(config, pool)
	runGrpcServer(config, pool)

	// runGinServer(config, pool)
}

func runDBMigration(migrationURL, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create new migrate instance")
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msg("failed to run migrate up")
	}

	log.Info().Msg("DB migrated successfully")
}

func runGrpcServer(config config.Config, pool *pgxpool.Pool) {
	repository := repository.NewRepository(pool)
	service := service.NewService(repository.R)
	server, err := gapi.NewServer(config, service.S)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
	}

	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger) // create a unary server interceptor
	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterSwiftBankServer(grpcServer, server)
	reflection.Register(grpcServer) // it allows the grpc client to explore what RPCs are available on the server and how to call them (i.e self documentation for the server)

	// create a new listener
	listener, err := net.Listen("tcp", config.GrpcServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create listener")
	}

	log.Info().Msgf("start gRPC server at: %s", listener.Addr().String())
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal().Err(err).Msg("cannot start gRPC server")
	}
}

func runGatewayServer(config config.Config, pool *pgxpool.Pool) {
	repository := repository.NewRepository(pool)
	service := service.NewService(repository.R)
	server, err := gapi.NewServer(config, service.S)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
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
		log.Fatal().Err(err).Msg("cannot register handler server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	// serve swagger file with statik
	statikFs, err := fs.New()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create statik fs")
	}

	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFs))
	mux.Handle("/swagger/", swaggerHandler)

	// create a new listener
	listener, err := net.Listen("tcp", config.HttpServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create listener")
	}

	log.Info().Msgf("start HTTP gateway server at: %s", listener.Addr().String())
	handler := gapi.HttpLogger(mux) // add the http gateway logger

	if err := http.Serve(listener, handler); err != nil {
		log.Fatal().Err(err).Msg("cannot start HTTP gateway server")
	}
}

func runGinServer(config config.Config, pool *pgxpool.Pool) {
	repository := repository.NewRepository(pool)
	service := service.NewService(repository.R)
	handler, err := api.NewHandler(config, service.S)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load handler")
	}

	err = handler.H.StartServer(config.HttpServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start server")
	}
}
