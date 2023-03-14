package gapi

import (
	"IMChat/pb"
	"context"
	"database/sql"
	"fmt"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {

	authPayload, err := server.authorization(ctx)
	if err != nil {
		return nil, unauthenticateError(err)
	}

	violations := validateDeleteUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	user, err := server.store.GetUser(ctx, authPayload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "用户名不存在: %v", err)
		}
		return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
	}

	if user.ID != req.GetId() {
		return nil, status.Errorf(codes.PermissionDenied, "cannot delete user info: %v", err)
	}

	err = server.store.DeleteUser(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "faild delete user: %v", err)
	}

	resp := &pb.DeleteUserResponse{
		Message: "Delete User Successfully",
	}

	return resp, nil
}

func validateDeleteUserRequest(req *pb.DeleteUserRequest) (violation []*errdetails.BadRequest_FieldViolation) {

	if req.GetId() <= 0 {
		violation = append(violation, fieldViolation("id", fmt.Errorf("invaid parameter")))
	}

	return violation
}
