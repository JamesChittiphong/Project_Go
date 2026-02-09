package ws

import (
	"sync"

	"github.com/gofiber/websocket/v2"
)

type Hub struct {
	clients    map[uint]*websocket.Conn
	register   chan *Client
	unregister chan uint
	mu         sync.RWMutex
}

type Client struct {
	UserID uint
	Conn   *websocket.Conn
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[uint]*websocket.Conn),
		register:   make(chan *Client),
		unregister: make(chan uint),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.UserID] = client.Conn
			h.mu.Unlock()
		case userID := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[userID]; ok {
				delete(h.clients, userID)
			}
			h.mu.Unlock()
		}
	}
}

func (h *Hub) Register(userID uint, conn *websocket.Conn) {
	h.register <- &Client{UserID: userID, Conn: conn}
}

func (h *Hub) Unregister(userID uint) {
	h.unregister <- userID
}

func (h *Hub) BroadcastToUser(userID uint, message interface{}) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if conn, ok := h.clients[userID]; ok {
		conn.WriteJSON(message)
	}
}
