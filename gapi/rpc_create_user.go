package gapi

import (
	db "IMChat/db/pg/sqlc"
	"IMChat/pb"
	"IMChat/utils"
	"IMChat/validate"
	"context"

	"github.com/lib/pq"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	violations := validateCreateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	hashPassword, err := utils.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
	}

	arg := db.CreateUserParams{
		Email:    req.GetEmail(),
		Username: req.GetUsername(),
		Nickname: req.GetUsername(),
		Password: hashPassword,
		Gender:   int16(req.GetGender()),
	}

	_, err = server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "username already exists: %s", err)
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
	}

	rsp := &pb.CreateUserResponse{
		Message: "Create User Successfully",
	}
	//
	return rsp, nil
}

func validateCreateUserRequest(req *pb.CreateUserRequest) (violation []*errdetails.BadRequest_FieldViolation) {
	if err := validate.ValidateEmail(req.GetEmail()); err != nil {
		violation = append(violation, fieldViolation("email", err))
	}

	if err := validate.ValidateUsername(req.GetUsername()); err != nil {
		violation = append(violation, fieldViolation("username", err))
	}

	if err := validate.ValidateLen(req.GetPassword(), 3, 20); err != nil {
		violation = append(violation, fieldViolation("password", err))
	}

	if err := validate.ValidateGender(req.GetGender()); err != nil {
		violation = append(violation, fieldViolation("gender", err))
	}

	return violation
}
