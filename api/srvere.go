package api

import (
	db "IMChat/db/pg/sqlc"
	"IMChat/token"
	"IMChat/utils/config"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	store      db.Store
	router     *gin.Engine
	tokenMaker token.Maker
	conf       config.Config
}

func NewServer(conf config.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(conf.Token.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}
	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
		conf:       conf,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("gender", validGender)
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/user", server.createUser)
	router.POST("/user/login", server.loginUser)

	authRuters := router.Group("").Use(authMiddleware(server.tokenMaker))
	authRuters.PUT("user", server.updateUser)

	server.router = router
}

func (server *Server) Start(addrss string) error {
	return server.router.Run(addrss)
}
