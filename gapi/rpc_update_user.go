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

func (server *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	arg := db.UpdateUserParams{
		Email: sql.NullString{
			String: req.GetEmail(),
			Valid:  req.Email != nil,
		},
		Nickname: sql.NullString{
			String: req.GetEmail(),
			Valid:  req.Nickname != nil,
		},
		Gender: sql.NullInt16{
			Int16: int16(req.GetGender()),
			Valid: req.Gender != nil,
		},
		Avatar: sql.NullString{
			String: req.GetAvatar(),
			Valid:  req.Avatar != nil,
		},
	}

	if req.Password != nil {
		hashPassword, err := utils.HashPassword(*req.Password)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
		}
		arg.Password = sql.NullString{
			String: hashPassword,
			Valid:  true,
		}
	}

	err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed update user: %s", err)
	}

	return &pb.UpdateUserResponse{}, nil
}
