package websocket

import (
	"encoding/json"
	"sync"

	"github.com/vpa/quanlynhahang-backend/internal/dto"
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

	h.mu.RLock()
	defer h.mu.RUnlock()

	for c := range h.Clients {
		if c.UserID == userID {
			c.Send <- data
		}
	}
}

func (h *Hub) BroadcastToRoom(roomID uint, msg dto.WSMessage) {
	data, _ := json.Marshal(msg)
	h.mu.RLock()
	defer h.mu.RUnlock()
	for c := range h.Clients {
		if c.RoomID == roomID {
			safeSend(c, data)
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
