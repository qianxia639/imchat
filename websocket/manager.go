package ws

import (
	"IMChat/pb"
	"IMChat/utils"
	"log"
	"sync"

	"google.golang.org/protobuf/encoding/protojson"
)

type Manager struct {
	mutex      *sync.RWMutex
	Clients    map[int64]*Client
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
}

func NewManager() *Manager {
	return &Manager{
		mutex:      &sync.RWMutex{},
		Clients:    make(map[int64]*Client),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (manager *Manager) Start() {
	for {
		select {
		case client := <-manager.Register:
			manager.mutex.Lock()
			manager.Clients[client.Id] = client
			manager.mutex.Unlock()
		case client := <-manager.Unregister:
			manager.mutex.RLock()
			_, ok := manager.Clients[client.Id]
			manager.mutex.RUnlock()
			if ok {
				delete(manager.Clients, client.Id)
				close(client.Send)
			}
		case message := <-manager.Broadcast:
			msg := &pb.Message{}
			_ = protojson.Unmarshal(message, msg)

			if msg.To > 0 {
				manager.dispatch(message, msg)
			} else {
				for _, client := range manager.Clients {
					select {
					case client.Send <- message:
					default:
						close(client.Send)
						delete(manager.Clients, client.Id)
					}
				}
			}
		}
	}
}

func (manager *Manager) dispatch(message []byte, msg *pb.Message) {
	// 普通消息，如文本消息，文件消息等
	if msg.ContentType >= utils.TEXT && msg.ContentType <= utils.IMAGE {
		manager.mutex.RLock()
		_, exits := manager.Clients[msg.From]
		manager.mutex.RUnlock()
		if exits {
			// 保存消息
			log.Print("保存成功")
		}

		// 判断单聊还是群聊
		if msg.MessageType == utils.MESSAGE_TYPE_USER { // 单聊
			manager.mutex.RLock()
			clint, ok := manager.Clients[msg.To]
			manager.mutex.RUnlock()
			if ok {
				msgByte, err := protojson.Marshal(msg)
				if err == nil {
					log.Printf("单聊: %s", string(msgByte))
					clint.Send <- msgByte
				}
			}
		} else if msg.MessageType == utils.MESSAGE_TYPE_GROUP { // 群聊
			// log.Printf("单聊: %s", string(msgByte))
			log.Printf("群聊")
		}

	} else {
		// 语音、视频等通话，不保存文件,直接进行转发
		manager.mutex.RLock()
		client, ok := manager.Clients[msg.To]
		manager.mutex.RUnlock()
		if ok {
			client.Send <- message
		}
	}
}
