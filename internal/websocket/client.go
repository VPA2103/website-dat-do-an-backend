package websocket

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	UserID uint
	Role   string
	RoomID uint
	Conn   *websocket.Conn
	Send   chan []byte
	Hub    *Hub
}

func NewClient(conn *websocket.Conn, hub *Hub, userID uint, role string) *Client {
	return &Client{
		UserID: userID,
		Role:   role,
		Conn:   conn,
		Send:   make(chan []byte, 256),
		Hub:    hub,
	}
}

func (c *Client) WritePump() {
	defer c.Conn.Close()
	for msg := range c.Send {
		if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			return
		}
	}
}
