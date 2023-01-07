package gapi

import (
	db "IMChat/db/pg/sqlc"
	"IMChat/pb"
	"IMChat/token"
	"IMChat/utils/config"
	ws "IMChat/websocket"
)

type Server struct {
	pb.UnimplementedUserServer
	pb.UnimplementedMessageServer
	store         db.Store
	tokenMaker    token.Maker
	conf          config.Config
	ClientManager *ws.ClientManager
}

func NewServer(conf config.Config, store db.Store) (*Server, error) {

	tokenMaker, err := token.NewPasetoMaker(conf.Token.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}

	server := &Server{
		store:         store,
		tokenMaker:    tokenMaker,
		conf:          conf,
		ClientManager: ws.NewClientManager(),
	}

	return server, nil
}
