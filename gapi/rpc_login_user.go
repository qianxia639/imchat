package gapi

import (
	db "IMChat/db/sqlc"
	"IMChat/pb"
	"IMChat/utils"
	"context"
	"database/sql"
	"fmt"
	"time"

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
			return nil, status.Errorf(codes.NotFound, "用户名不存在: %v", err)
		}
		return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
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
		return nil, status.Errorf(codes.NotFound, "用户名或密码错误: %v", err)
	}

	token, err := server.tokenMaker.CreateToken(user2.Username, server.conf.Token.AccessTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create access token: %v", err)
	}

	key := fmt.Sprintf("useer:%d_%v", user2.ID, user2.Username)
	err = server.cache.SetTtlCache(ctx, key, &user, 24*time.Hour)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to set cache: %v", err)
	}

	resp := &pb.LoginUserResponse{
		Token: token,
		User:  converUser(user),
	}

	return resp, nil
}

func validateLoginUserRequest(req *pb.LoginUserRequest) (violation []*errdetails.BadRequest_FieldViolation) {

	if ok := utils.IsEmpty(req.GetUsername()); ok {
		violation = append(violation, fieldViolation("username", fmt.Errorf("cannot empty")))
	}

	if err := utils.ValidateLen(req.GetPassword(), 3, 20); err != nil {
		violation = append(violation, fieldViolation("password", err))
	}

	return violation
}
