package api

import (
	db "IMChat/db/sqlc"
	"IMChat/token"
	"IMChat/utils"
	"IMChat/utils/config"

	"IMChat/cache"
	ws "IMChat/websocket"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type Server struct {
	store      db.Store
	cache      cache.Cache
	router     *gin.Engine
	tokenMaker token.Maker
	conf       config.Config
	manager    *ws.Manager
	logger     *zap.Logger
}

func NewServer(conf config.Config, store db.Store, cache cache.Cache) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(conf.Token.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}
	server := &Server{
		store:      store,
		cache:      cache,
		tokenMaker: tokenMaker,
		conf:       conf,
	}

	server.logger = utils.Zap(conf.Logger.Path, conf.Logger.Level)

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("gender", validGender)
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	manager := ws.NewManager()

	router.POST("/user", server.createUser)
	router.POST("/user/login", server.loginUser)

	authRuters := router.Group("/").Use(authMiddleware(server.tokenMaker))
	{
		authRuters.PUT("/user", server.updateUser)
		authRuters.DELETE("/user/:id", server.deleteUser)

		authRuters.GET("/ws", server.socketHandler)
	}

	server.router = router
	server.manager = manager
}

func (server *Server) Start(addrss string) error {
	return server.router.Run(addrss)
}
