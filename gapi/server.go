package gapi

import (
	"IMChat/cache"
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
	cache      cache.Cache
}

func NewServer(conf config.Config, store db.Store, cache cache.Cache) (*Server, error) {

	tokenMaker, err := token.NewPasetoMaker(conf.Token.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}

	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
		conf:       conf,
		cache:      cache,
	}

	return server, nil
}
