package ws

import (
	"log"
	"sync"
)

type Hub struct {
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.token] = client
			h.mu.Unlock()
			log.Printf("WebSocket client registered for token: %s", client.token)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.token]; ok {
				delete(h.clients, client.token)
				close(client.send)
			}
			h.mu.Unlock()
			log.Printf("WebSocket client unregistered for token: %s", client.token)
		}
	}
}

func (h *Hub) Register(client *Client) {
	h.register <- client
}

func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}

func (h *Hub) SendToToken(token string, message []byte) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if client, ok := h.clients[token]; ok {
		select {
		case client.send <- message:
			return true
		default:
			return false
		}
	}
	return false
}
