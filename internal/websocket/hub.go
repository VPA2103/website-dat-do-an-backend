package websocket

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/vpa/quanlynhahang-backend/dto"
)

type Hub struct {
	Clients    map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	mu         sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client] = true
			h.mu.Unlock()

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}
			h.mu.Unlock()
		}
	}
}

func safeSend(c *Client, data []byte) {
	select {
	case c.Send <- data:
	default:
		// Client quá chậm hoặc đã chết
		close(c.Send)
	}
}

func (h *Hub) SendToUser(userID uint, msg dto.WSMessage) {

	data, _ := json.Marshal(msg)

	h.mu.Lock()
	defer h.mu.Unlock()

	for c := range h.Clients {

		if c.UserID == userID {

			select {

			case c.Send <- data:

			default:
				close(c.Send)
				delete(h.Clients, c)
			}
		}
	}
}

func (h *Hub) Broadcast(msg dto.WSMessage) {
	log.Printf("Broadcast: %+v\n", msg)
	data, _ := json.Marshal(msg)

	h.mu.Lock()
	defer h.mu.Unlock()

	for c := range h.Clients {

		select {

		case c.Send <- data:

		default:
			close(c.Send)
			delete(h.Clients, c)
		}
	}
}

func (h *Hub) SendToRole(role string, msg dto.WSMessage) {
	data, _ := json.Marshal(msg)
	h.mu.RLock()
	defer h.mu.RUnlock()
	for c := range h.Clients {
		if c.Role == role {
			c.Send <- data
		}
	}
}
