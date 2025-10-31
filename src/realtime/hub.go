package realtime

import "sync"

type Hub struct {
	Clients    map[int]map[int]*Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan Message
	Mutex      sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		Clients: make(map[int]map[int]*Client),
		Register: make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast: make(chan Message),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Mutex.Lock()
			if h.Clients[client.MeetingID] == nil {
				h.Clients[client.MeetingID] = make(map[int]*Client)
			}
			h.Clients[client.MeetingID][client.UserID] = client
			h.Mutex.Unlock()

		case client := <-h.Unregister:
			h.Mutex.Lock()
			if _, ok := h.Clients[client.MeetingID][client.UserID]; ok {
				delete(h.Clients[client.MeetingID], client.UserID)
				close(client.Send)
			}
			h.Mutex.Unlock()

		case msg := <-h.Broadcast:
			h.Mutex.Lock()
			for _, client := range h.Clients[msg.MeetingID] {
				if client.UserID != msg.UserID {
					select {
					case client.Send <- msg.Data:
					default:
						close(client.Send)
						delete(h.Clients[msg.MeetingID], client.UserID)
					}
				}
			}
			h.Mutex.Unlock()
		}
	}
}