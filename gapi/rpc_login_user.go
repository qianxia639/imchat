package gapi

import (
	db "IMChat/db/pg/sqlc"
	"IMChat/pb"
	"IMChat/utils"
	"context"
	"database/sql"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {

	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "用户名不存在: %s", err)
		}
		return nil, status.Errorf(codes.NotFound, "user not found: %s", err)
	}

	err = utils.CheckPassword(req.Password, user.Password)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "incorrect password")
	}

	user2, err := server.store.LoginUser(ctx, db.LoginUserParams{
		Username: user.Username,
		Password: user.Password,
	})

	if err != nil {
		return nil, status.Errorf(codes.NotFound, "用户名或密码错误: %s", err)
	}

	resp := &pb.LoginUserResponse{
		Token: user2.Username,
	}

	return resp, nil
}
