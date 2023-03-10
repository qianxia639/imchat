package main

import (
	"IMChat/api"
	"IMChat/cache"
	db "IMChat/db/sqlc"
	"IMChat/gapi"
	"IMChat/pb"
	"IMChat/utils"
	"IMChat/utils/config"
	"context"
	"database/sql"
	"embed"
	"io/fs"
	"log"
	"net"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/tmc/grpc-websocket-proxy/wsproxy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"

	_ "github.com/lib/pq"
)

//go:embed doc/swagger/*
var embedFs embed.FS

func main() {
	conf, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("load config err: ", err)
	}

	conn, err := sql.Open(conf.Postgres.Driver, conf.Postgres.Source)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	runDBMigrate(conf.Postgres.MigrateUrl, conf.Postgres.Source)

	store := db.NewStore(conn)

	cache := cache.NewRedisCache(utils.InitRedis(conf))

	go runGatewayServer(conf, store, cache)
	runGrpcServer(conf, store, cache)
	// runGinServer(conf, store, cache)
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

func runGrpcServer(conf config.Config, store db.Store, cache cache.Cache) {
	server, err := gapi.NewServer(conf, store, cache)
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

func runGatewayServer(conf config.Config, store db.Store, cache cache.Cache) {
	server, err := gapi.NewServer(conf, store, cache)
	if err != nil {
		log.Fatal("cannot create server: ", err)
	}

	grpcMux := runtime.NewServeMux(
		// ???????????????
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterUserHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("conot register handler server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	staticFile, err := fs.Sub(embedFs, "doc/swagger")
	if err != nil {
		log.Fatal("fs error: ", err)
	}
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", http.FileServer(http.FS(staticFile))))

	listener, err := net.Listen("tcp", conf.Server.HttpServerAddress)
	if err != nil {
		log.Fatal("cannot create listener: ", err)
	}
	log.Printf("start HTTP gateway server %s", listener.Addr().String())
	err = http.Serve(listener, wsproxy.WebsocketProxy(mux))
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}

func runGinServer(conf config.Config, store db.Store, cache cache.Cache) {
	server, err := api.NewServer(conf, store, cache)
	if err != nil {
		log.Fatal("cannot create server: ", err)
	}

	log.Fatal(server.Start(conf.Server.HttpServerAddress))
}
