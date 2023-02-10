package gapi

import (
	db "IMChat/db/pg/sqlc"
	"IMChat/pb"
	"IMChat/utils"
	"IMChat/validate"
	"context"
	"database/sql"
	"fmt"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	violations := validateLoginUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

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

	token, err := server.tokenMaker.CreateToken(user2.Username, server.conf.Token.AccessTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "filed to create access token: %s", err)
	}

	resp := &pb.LoginUserResponse{
		Token: token,
	}

	return resp, nil
}

func validateLoginUserRequest(req *pb.LoginUserRequest) (violation []*errdetails.BadRequest_FieldViolation) {

	if ok := validate.IsEmpty(req.GetUsername()); ok {
		violation = append(violation, fieldViolation("username", fmt.Errorf("cannot empty")))
	}

	if err := validate.ValidateLen(req.GetPassword(), 3, 20); err != nil {
		violation = append(violation, fieldViolation("password", err))
	}

	return violation
}
