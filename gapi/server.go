package gapi

import (
	db "IMChat/db/pg/sqlc"
	"IMChat/pb"
	"IMChat/token"
)

type Server struct {
	pb.UnimplementedUserServer
	store      db.Store
	tokenMaker token.Maker
}

const key = "plokmnjiuhbvgytfcxdreszawq564738"

func NewServer(store db.Store) (*Server, error) {

	tokenMaker, err := token.NewPasetoMaker(key)
	if err != nil {
		return nil, err
	}

	return &Server{store: store, tokenMaker: tokenMaker}, nil
}
