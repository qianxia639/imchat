package main

import (
	db "IMChat/db/pg/sqlc"
	"IMChat/gapi"
	"IMChat/pb"
	"IMChat/utils/config"
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "github.com/lib/pq"
)

func main() {
	conf, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("load config err: ", err)
	}

	conn, err := sql.Open(conf.Postgres.Driver, conf.Postgres.Source)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	runDBMigrate(conf.Postgres.Migration.MigrateUrl, conf.Postgres.Source)

	store := db.NewStore(conn)
	go runGatewayServer(conf, store)
	runGrpcServer(conf, store)
}

func runDBMigrate(migrationUrl, dbSource string) {
	migration, err := migrate.New(migrationUrl, dbSource)
	if err != nil {
		log.Fatal("cannot create new migrate instance: ", err)
	}

	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("failed to run migrate up: ", err)
	}

	log.Println("db migrated successfully")
}

func runGrpcServer(conf config.Config, store db.Store) {
	server, err := gapi.NewServer(conf, store)
	if err != nil {
		log.Fatal("cannot create server: ", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", conf.Server.GrpcServerAddress)
	if err != nil {
		log.Fatal("cannot create listener: ", err)
	}
	log.Printf("start gRPC server %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}

func runGatewayServer(conf config.Config, store db.Store) {
	server, err := gapi.NewServer(conf, store)
	if err != nil {
		log.Fatal("cannot create server: ", err)
	}

	grpcMux := runtime.NewServeMux()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterUserHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("conot register handler server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	listener, err := net.Listen("tcp", conf.Server.HttpServerAddress)
	if err != nil {
		log.Fatal("cannot create listener: ", err)
	}
	log.Printf("start HTTP gateway server %s", listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
