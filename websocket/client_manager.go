package ws

type ClientManager struct {
	Clients    map[uint64]*Client
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
}

func NewClientManager() *ClientManager {
	return &ClientManager{
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[uint64]*Client),
	}
}

func (manager *ClientManager) Run() {
	for {
		select {
		case client := <-manager.Register:
			manager.Clients[client.UserId] = client
		case client := <-manager.Unregister:
			if _, ok := manager.Clients[client.UserId]; ok {
				delete(manager.Clients, client.UserId)
				close(client.Send)
			}
		case message := <-manager.Broadcast:
			for _, client := range manager.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(manager.Clients, client.UserId)
				}
			}
		}
	}
}
