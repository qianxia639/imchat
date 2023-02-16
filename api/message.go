package api

import (
	ws "IMChat/websocket"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func (server *Server) socketHandler(ctx *gin.Context) {
	// 协议升级，将http协议升级为websocket协议
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}).Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.SecureJSON(http.StatusInternalServerError, err.Error())
		return
	}

	fromId, err := strconv.ParseInt(ctx.Query("from_id"), 10, 64)
	if err != nil {
		ctx.SecureJSON(http.StatusBadRequest, "invalid parameter")
		return
	}
	// 创建用户实例
	client := &ws.Client{
		Manager: *server.manager,
		Conn:    conn,
		Send:    make(chan []byte),
		Id:      fromId,
	}

	server.manager.Register <- client

	go client.Read()

	go client.Writ()
}
