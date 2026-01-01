package realtime

import (
	"encoding/json"
	"sync"

	"github.com/gofiber/websocket/v2"
)

// Client represents a WebSocket client connection
type Client struct {
	Hub      *Hub
	Conn     *websocket.Conn
	Send     chan []byte
	SchoolID uint
	ClassID  *uint  // Optional filter by class
	IsPublic bool   // For public display
	Token    string // Display token for public
	UserID   uint   // User ID for authenticated clients
}

// Hub maintains the set of active clients and broadcasts messages to them
// Requirements: 4.2 - Broadcast to school-specific clients
// Requirements: 4.5 - Filter by class_id if specified
type Hub struct {
	// Registered clients grouped by school ID
	clients map[uint]map[*Client]bool

	// Broadcast channel for attendance events
	broadcast chan *AttendanceEvent

	// Register requests from clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Mutex for thread-safe operations
	mu sync.RWMutex
}

// NewHub creates a new Hub instance
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[uint]map[*Client]bool),
		broadcast:  make(chan *AttendanceEvent, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)
		case client := <-h.unregister:
			h.unregisterClient(client)
		case event := <-h.broadcast:
			h.broadcastEvent(event)
		}
	}
}

// registerClient adds a client to the hub
func (h *Hub) registerClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.clients[client.SchoolID] == nil {
		h.clients[client.SchoolID] = make(map[*Client]bool)
	}
	h.clients[client.SchoolID][client] = true
}

// unregisterClient removes a client from the hub
func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if clients, ok := h.clients[client.SchoolID]; ok {
		if _, ok := clients[client]; ok {
			delete(clients, client)
			close(client.Send)
			if len(clients) == 0 {
				delete(h.clients, client.SchoolID)
			}
		}
	}
}

// broadcastEvent sends an event to all relevant clients
// Requirements: 4.2 - Update dashboard within 3 seconds without page refresh
// Requirements: 4.5 - Allow filtering by class
func (h *Hub) broadcastEvent(event *AttendanceEvent) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	clients, ok := h.clients[event.SchoolID]
	if !ok {
		return
	}

	message, err := json.Marshal(event)
	if err != nil {
		return
	}

	for client := range clients {
		// Filter by class if client has class filter and event has attendance data
		if client.ClassID != nil && event.Attendance != nil {
			if event.Attendance.ClassID != *client.ClassID {
				continue
			}
		}

		select {
		case client.Send <- message:
		default:
			// Client's send buffer is full, close connection
			h.mu.RUnlock()
			h.unregisterClient(client)
			h.mu.RLock()
		}
	}
}

// Register adds a client to the hub
func (h *Hub) Register(client *Client) {
	h.register <- client
}

// Unregister removes a client from the hub
func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}

// Broadcast sends an event to all relevant clients
func (h *Hub) Broadcast(event *AttendanceEvent) {
	h.broadcast <- event
}

// GetClientCount returns the number of connected clients for a school
func (h *Hub) GetClientCount(schoolID uint) int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if clients, ok := h.clients[schoolID]; ok {
		return len(clients)
	}
	return 0
}

// GetTotalClientCount returns the total number of connected clients
func (h *Hub) GetTotalClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	total := 0
	for _, clients := range h.clients {
		total += len(clients)
	}
	return total
}
