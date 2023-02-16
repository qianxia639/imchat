package ws

type Manager struct {
	Clients    map[int64]*Client
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
}

func NewManager() *Manager {
	return &Manager{
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
			manager.Clients[client.Id] = client
		case client := <-manager.Unregister:
			if _, ok := manager.Clients[client.Id]; ok {
				delete(manager.Clients, client.Id)
				close(client.Send)
			}
		case message := <-manager.Broadcast:
			for _, client := range manager.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					// delete(manager.Clients,)
				}
			}
		}
	}
}
