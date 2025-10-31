package services

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gofiber/fiber/v3"
	"github.com/saveblush/gofiber3-contrib/websocket"
)

type Client struct {
	Conn      *websocket.Conn
	ProjectID string
}

type WebSocketHub struct {
	clients    map[string][]*Client // projectID -> clients
	broadcast  chan BroadcastMessage
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

type BroadcastMessage struct {
	ProjectID string
	Data      interface{}
}

func NewWebSocketHub() *WebSocketHub {
	hub := &WebSocketHub{
		clients:    make(map[string][]*Client),
		broadcast:  make(chan BroadcastMessage, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}

	go hub.run()

	log.Println("✅ WebSocket hub initialized")
	return hub
}

func (h *WebSocketHub) run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.ProjectID] = append(h.clients[client.ProjectID], client)
			h.mu.Unlock()
			log.Printf("📡 Client connected for project: %s", client.ProjectID)

		case client := <-h.unregister:
			h.mu.Lock()
			if clients, ok := h.clients[client.ProjectID]; ok {
				for i, c := range clients {
					if c == client {
						h.clients[client.ProjectID] = append(clients[:i], clients[i+1:]...)
						break
					}
				}
				if len(h.clients[client.ProjectID]) == 0 {
					delete(h.clients, client.ProjectID)
				}
			}
			h.mu.Unlock()
			log.Printf("📡 Client disconnected from project: %s", client.ProjectID)

		case message := <-h.broadcast:
			h.mu.RLock()
			clients := h.clients[message.ProjectID]
			h.mu.RUnlock()

			jsonData, err := json.Marshal(message.Data)
			if err != nil {
				log.Printf("❌ Failed to marshal message: %v", err)
				continue
			}

			for _, client := range clients {
				if err := client.Conn.WriteMessage(websocket.TextMessage, jsonData); err != nil {
					log.Printf("❌ Failed to send message to client: %v", err)
					h.unregister <- client
				}
			}
		}
	}
}

func (h *WebSocketHub) Register(client *Client) {
	h.register <- client
}

func (h *WebSocketHub) Unregister(client *Client) {
	h.unregister <- client
}

func (h *WebSocketHub) BroadcastToProject(projectID string, data interface{}) {
	h.broadcast <- BroadcastMessage{
		ProjectID: projectID,
		Data:      data,
	}
}

func (h *WebSocketHub) HandleWebSocket(c fiber.Ctx) error {
	projectID := c.Params("projectId")

	if projectID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Project ID is required",
		})
	}

	log.Printf("📡 WebSocket connection request for project: %s", projectID)

	// Upgrade to WebSocket
	if websocket.IsWebSocketUpgrade(c) {
		return websocket.New(func(conn *websocket.Conn) {
			client := &Client{
				Conn:      conn,
				ProjectID: projectID,
			}

			log.Printf("📡 WebSocket client connected for project: %s", projectID)
			h.Register(client)
			defer func() {
				log.Printf("📡 WebSocket client disconnecting from project: %s", projectID)
				h.Unregister(client)
			}()

			// Keep connection alive and listen for close
			for {
				_, _, err := conn.ReadMessage()
				if err != nil {
					log.Printf("📡 WebSocket read error for project %s: %v", projectID, err)
					break
				}
			}
		})(c)
	}

	return c.Status(426).JSON(fiber.Map{
		"error": "WebSocket upgrade required",
	})
}

