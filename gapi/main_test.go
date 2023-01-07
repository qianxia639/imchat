package gapi

import (
	"IMChat/pb"
	"IMChat/utils/config"
	"database/sql"
	"log"
	"os"
	"testing"

	db "IMChat/db/pg/sqlc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var server *Server
var clientConn *grpc.ClientConn
var userClient pb.UserClient

func TestMain(m *testing.M) {

	conf, err := config.LoadConfig("../")
	if err != nil {
		log.Fatal("load config err: ", err)
	}

	conn, err := sql.Open(conf.Postgres.Driver, conf.Postgres.Source)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	store := db.NewStore(conn)

	server, err = NewServer(conf, store)
	if err != nil {
		log.Fatal("connot create server: ", err)
	}

	clientConn, err = grpc.Dial(":9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("conot create client: ", err)
	}
	defer clientConn.Close()
	userClient = pb.NewUserClient(clientConn)

	os.Exit(m.Run())

}
