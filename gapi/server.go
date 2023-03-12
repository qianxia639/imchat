package gapi

import (
	db "IMChat/db/sqlc"
	"IMChat/pb"
	"IMChat/token"
	"IMChat/utils/config"
)

type Server struct {
	pb.UnimplementedUserServer
	pb.UnimplementedContactServer
	store      db.Store
	tokenMaker token.Maker
	conf       config.Config
}

func NewServer(conf config.Config, store db.Store) (*Server, error) {

	tokenMaker, err := token.NewPasetoMaker(conf.Token.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}

	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
		conf:       conf,
	}

	return server, nil
}
