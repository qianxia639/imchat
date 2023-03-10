package api

import (
	db "IMChat/db/sqlc"
	"IMChat/token"
	"IMChat/utils"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createUserReequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Email    string `json:"email" binding:"required,email"`
	Nickname string `json:"nickname" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Gender   int16  `json:"gender" binding:"required,gender"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserReequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.SecureJSON(http.StatusBadRequest, err.Error())
		return
	}

	hashPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		ctx.SecureJSON(http.StatusInternalServerError, err.Error())
		return
	}

	arg := db.CreateUserParams{
		Email:    req.Email,
		Username: req.Username,
		Nickname: req.Nickname,
		Password: hashPassword,
		Gender:   req.Gender,
	}

	_, err = server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.SecureJSON(http.StatusForbidden, fmt.Sprintf("username already exists: %s", err.Error()))
				return
			}
		}
		ctx.SecureJSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.SecureJSON(http.StatusOK, "Create User Successfully")

}

type loginUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	Token string  `json:"token"`
	User  db.User `jsn:"user"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.SecureJSON(http.StatusNotFound, err.Error())
			return
		}
		ctx.SecureJSON(http.StatusInternalServerError, err.Error())
		return
	}

	err = utils.CheckPassword(req.Password, user.Password)
	if err != nil {
		ctx.SecureJSON(http.StatusUnauthorized, err.Error())
		return
	}

	token, err := server.tokenMaker.CreateToken(user.Username, server.conf.Token.AccessTokenDuration)
	if err != nil {
		ctx.SecureJSON(http.StatusInternalServerError, err.Error())
		return
	}

	key := fmt.Sprintf("useer:%d_%s", user.ID, user.Username)
	err = server.cache.SetTtlCache(ctx, key, &user, 24*time.Hour)
	if err != nil {
		ctx.SecureJSON(http.StatusInternalServerError, err.Error())
		return
	}

	resp := loginUserResponse{
		Token: token,
		User: db.User{
			ID:       user.ID,
			Username: user.Username,
			Nickname: user.Nickname,
			Gender:   user.Gender,
			Avatar:   user.Avatar,
		},
	}

	ctx.SecureJSON(http.StatusOK, resp)
}

type UpdateUserRequeest struct {
	Username string  `json:"username" binding:"required,alphanum"`
	Email    *string `json:"email"`
	Nickname *string `json:"nickname"`
	Password *string `json:"password"`
	Gender   *int16  `json:"gender"`
	Avatar   *string `json:"avatar"`
}

func (server *Server) updateUser(ctx *gin.Context) {
	var req UpdateUserRequeest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.SecureJSON(http.StatusBadRequest, err.Error())
		return
	}

	val, ok := ctx.Get(authorizationHeader)
	if !ok {
		ctx.SecureJSON(http.StatusInternalServerError, "missing key")
		return
	}

	payload, ok := val.(*token.Payload)
	if !ok {
		ctx.SecureJSON(http.StatusBadRequest, "?????????????????????")
		return
	}

	if req.Username != payload.Username {
		ctx.SecureJSON(http.StatusBadRequest, "???????????????")
		return
	}

	arg := db.UpdateUserParams{
		Username: req.Username,
	}

	if req.Email != nil {
		arg.Email = sql.NullString{
			String: *req.Email,
			Valid:  req.Email != nil,
		}
	}

	if req.Password != nil {
		hashPassword, err := utils.HashPassword(*req.Password)
		if err != nil {
			ctx.SecureJSON(http.StatusInternalServerError, err.Error())
			return
		}
		arg.Password = sql.NullString{
			String: hashPassword,
			Valid:  true,
		}
	}

	if req.Nickname != nil {
		arg.Nickname = sql.NullString{
			String: *req.Nickname,
			Valid:  req.Nickname != nil,
		}
	}

	if req.Gender != nil {
		arg.Gender = sql.NullInt16{
			Int16: *req.Gender,
			Valid: req.Gender != nil,
		}
	}

	if req.Avatar != nil {
		arg.Avatar = sql.NullString{
			String: *req.Avatar,
			Valid:  req.Avatar != nil,
		}
	}

	user, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		ctx.SecureJSON(http.StatusInternalServerError, err.Error())
		return
	}

	key := fmt.Sprintf("user:%d_%s", user.ID, user.Username)
	_ = server.cache.SetCache(ctx, key, user)

	ctx.SecureJSON(http.StatusOK, "Update Usere Successfully")
}

func (server *Server) deleteUser(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		ctx.SecureJSON(http.StatusBadRequest, "invalid parameter")
		return
	}

	val, ok := ctx.Get(authorizationPayloadKey)
	if !ok {
		ctx.SecureJSON(http.StatusInternalServerError, "missing key")
		return
	}

	payload, ok := val.(*token.Payload)
	if !ok {
		ctx.SecureJSON(http.StatusBadRequest, "?????????????????????")
		return
	}

	user, err := server.store.GetUser(ctx, payload.Username)
	if err != nil || id != user.ID {
		ctx.SecureJSON(http.StatusInternalServerError, err.Error())
		return
	}

	_ = server.store.DeleteUser(ctx, id)

	ctx.SecureJSON(http.StatusOK, "Delete Usere Successfully")
}
