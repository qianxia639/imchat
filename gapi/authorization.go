package gapi

import (
	"IMChat/token"
	"context"
	"fmt"

	"google.golang.org/grpc/metadata"
)

const (
	authorizationHeader = "authorization"
)

func (server *Server) authorization(ctx context.Context) (*token.Payload, error) {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}

	values := md.Get(authorizationHeader)
	if len(values) == 0 {
		return nil, fmt.Errorf("missing authorization header")
	}

	accessToken := values[0]
	payload, err := server.tokenMaker.VerifyToken(accessToken)
	if err != nil {
		return nil, fmt.Errorf("无效的token: %s", err)
	}
	return payload, nil
}
