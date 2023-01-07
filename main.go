package main

import (
	db "IMChat/db/pg/sqlc"
	"IMChat/gapi"
	"IMChat/pb"
	"IMChat/utils/config"
	"context"
	"database/sql"
	"embed"
	"io/fs"
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/websocket"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"

	_ "github.com/lib/pq"

	ws "IMChat/websocket"
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
	// pb.RegisterUserServer(grpcServer, server)
	pb.RegisterMessageServer(grpcServer, server)
	reflection.Register(grpcServer)

	go server.ClientManager.Run()

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

	grpcMux := runtime.NewServeMux(
		// 下划线命明
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
	mux.Handle("/ws", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := (&websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		}).Upgrade(w, r, nil)
		if err != nil {
			// return err
			log.Fatal("http -> websocket failed: ", err)
		}
		defer conn.Close()

		id, _ := strconv.ParseUint(r.URL.Query().Get("id"), 10, 64)

		// 创建用户实例
		client := &ws.Client{
			Manager: server.ClientManager,
			Conn:    conn,
			Send:    make(chan []byte),
			UserId:  id,
		}

		// 注册用户到用户管理
		server.ClientManager.Register <- client
		// 发送消息
		go client.Write(ctx)
		// 接收消息
		go client.Read(ctx)
	}))

	staticFile, _ := fs.Sub(embedFs, "doc/swagger")
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", http.FileServer(http.FS(staticFile))))

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

func runWebsocketServer() {
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}).Upgrade(nil, nil, nil)
	if err != nil {
		log.Fatal("", err)
	}
	defer conn.Close()
}
