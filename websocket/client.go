package ws

import (
	"context"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

// 客户端
type Client struct {
	Manager *ClientManager
	Conn    *websocket.Conn // 连接对象
	Send    chan []byte     //消息数据channel
	UserId  uint64          // 用户Id
	// Heartbeat uint64          // 心跳时间
}

func (c *Client) Read(ctx context.Context) {
	defer func() {
		c.Manager.Unregister <- c
		_ = c.Conn.Close()
	}()

	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(appData string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			log.Printf("read message error: %v\n", err)
			c.Manager.Unregister <- c
			_ = c.Conn.Close()
			break
		}
		log.Printf("read message: %s\n", message)

		c.Manager.Broadcast <- message
		log.Printf("[message] <<< %s\n", string(message))
	}
}

func (c *Client) Write(ctx context.Context) {

	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		_ = c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			log.Printf("writer send data: %s, ok: %t\n", message, ok)
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Printf("nextWriter message error: %v\n", err)
				return
			}
			w.Write(message)

			// 将排队的聊天消息添加到当前websocket消息中。
			for i := 0; i < len(c.Send); i++ {
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
