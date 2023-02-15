package api

import (
	"IMChat/token"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader     = "authorization"
	authorizationPayloadKey = "authorization_payload"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorization := ctx.GetHeader(authorizationHeader)

		if len(authorization) == 0 {
			err := fmt.Errorf("missing authorization header")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
			return
		}

		payload, err := tokenMaker.VerifyToken(authorization)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}
