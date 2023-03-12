package gapi

import (
	db "IMChat/db/sqlc"
	"IMChat/pb"
	"IMChat/utils"
	"context"
	"database/sql"
	"fmt"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {

	authPayload, err := server.authorization(ctx)
	if err != nil {
		return nil, unauthenticateError(err)
	}

	violations := validateUpdateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	if authPayload.Username != req.GetUsername() {
		return nil, status.Errorf(codes.PermissionDenied, "权限不足")
	}

	arg := db.UpdateUserParams{
		Username: req.GetUsername(),
		Email: sql.NullString{
			String: req.GetEmail(),
			Valid:  req.Email != nil,
		},
		Nickname: sql.NullString{
			String: req.GetNickname(),
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

	_, err = server.store.UpdateUser(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found: %s", err)
		}
		return nil, status.Errorf(codes.Internal, "failed update user: %s", err)
	}

	resp := &pb.UpdateUserResponse{
		Message: "Update User Successfully",
	}

	return resp, nil
}

func validateUpdateUserRequest(req *pb.UpdateUserRequest) (violation []*errdetails.BadRequest_FieldViolation) {

	if err := utils.ValidateUsername(req.GetUsername()); err != nil {
		violation = append(violation, fieldViolation("username", err))
	}

	if req.Email != nil {
		if err := utils.ValidateEmail(req.GetEmail()); err != nil {
			violation = append(violation, fieldViolation("email", err))
		}
	}

	if req.Nickname != nil {
		if err := utils.ValidateUsername(req.GetNickname()); err != nil {
			violation = append(violation, fieldViolation("nickname", err))
		}
	}

	if req.Password != nil {
		if err := utils.ValidateLen(req.GetPassword(), 3, 20); err != nil {
			violation = append(violation, fieldViolation("password", err))
		}
	}

	if req.Gender != nil {
		if err := utils.ValidateGender(req.GetGender()); err != nil {
			violation = append(violation, fieldViolation("gender", err))
		}
	}

	if req.Avatar != nil {
		if ok := utils.IsEmpty(req.GetAvatar()); ok {
			violation = append(violation, fieldViolation("avatar", fmt.Errorf("cannot empty")))
		}
	}

	return violation
}
