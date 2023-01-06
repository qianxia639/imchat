package gapi

import (
	db "IMChat/db/pg/sqlc"
	"IMChat/pb"
)

type Server struct {
	pb.UnimplementedUserServer
	store db.Store
}

func NewServer(store db.Store) *Server {
	return &Server{store: store}
}
