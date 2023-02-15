package api

import (
	db "IMChat/db/pg/sqlc"
	"IMChat/utils"
	"database/sql"
	"fmt"
	"net/http"

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

	ctx.SecureJSON(http.StatusOK, token)
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

	_, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		ctx.SecureJSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.SecureJSON(http.StatusOK, "Update Usere Successfully")
}
