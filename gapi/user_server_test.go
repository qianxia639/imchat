package gapi

import (
	"IMChat/pb"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

func login(t *testing.T, req *pb.LoginUserRequest) *pb.LoginUserResponse {
	resp, err := server.LoginUser(context.Background(), &pb.LoginUserRequest{
		Username: req.GetUsername(),
		Password: req.GetPassword(),
	})
	require.NoError(t, err)
	require.NotEmpty(t, resp)

	return resp
}

func TestLogin(t *testing.T) {

	// req := &pb.CreateUserRequest{
	// 	Username: utils.RandomString(6),
	// 	Email:    utils.RandomEmail(),
	// 	Nickname: utils.RandomString(6),
	// 	Password: utils.RandomString(6),
	// 	Gender:   int32(utils.RandomGender()),
	// }

	// server.CreateUser(context.Background(), req)

	// login(t, &pb.LoginUserRequest{
	// 	Username: req.GetUsername(),
	// 	Password: req.GetPassword(),
	// })

	resp, err := userClient.LoginUser(context.Background(), &pb.LoginUserRequest{
		Username: "alice",
		Password: "alice",
	})
	require.NoError(t, err)
	require.NotEmpty(t, resp)
}

func TestDeleteUser(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	resp, err := userClient.LoginUser(ctx, &pb.LoginUserRequest{
		Username: "alice",
		Password: "alice",
	})
	require.NoError(t, err)
	require.NotEmpty(t, resp)

	md, ok := metadata.FromOutgoingContext(ctx)
	require.True(t, !ok)

	md.Set(authorizationHeader, resp.GetToken())

	resp2, err := userClient.DeleteUser(context.Background(), &pb.DeleteUserRequest{
		Id: 1,
	})

	require.NoError(t, err)
	require.Empty(t, resp2)
}
