package gapi

import (
	db "IMChat/db/pg/sqlc"
	"IMChat/pb"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func converUser(user db.User) *pb.Users {
	return &pb.Users{
		Username:     user.Username,
		Email:        user.Email,
		Nickname:     user.Nickname,
		Gender:       int32(user.Gender),
		Avatar:       user.Avatar,
		RegisterTime: timestamppb.New(user.RegisterTime),
	}
}
