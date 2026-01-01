package publicdisplay

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"

	"github.com/school-management/backend/internal/modules/realtime"
)

// Handler handles HTTP and WebSocket requests for public display
type Handler struct {
	service        Service
	realtimeHub    *realtime.Hub
	realtimeRepo   realtime.Repository
}

// NewHandler creates a new public display handler
func NewHandler(service Service, realtimeHub *realtime.Hub, realtimeRepo realtime.Repository) *Handler {
	return &Handler{
		service:        service,
		realtimeHub:    realtimeHub,
		realtimeRepo:   realtimeRepo,
	}
}

// RegisterPublicRoutes registers public display routes (no auth required)
// Requirements: 5.3 - Accessing public display URL with valid token SHALL show attendance data without login
func (h *Handler) RegisterPublicRoutes(router fiber.Router) {
	public := router.Group("/public/display")

	// REST endpoint for public display data
	// GET /api/v1/public/display/:token
	public.Get("/:token", h.GetPublicDisplayData)

	// WebSocket upgrade middleware for public display
	public.Use("/:token/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	// WebSocket endpoint for public display real-time updates
	// GET /api/v1/public/display/:token/ws
	public.Get("/:token/ws", websocket.New(h.HandlePublicWebSocket, websocket.Config{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}))
}

// GetPublicDisplayData handles getting public display data via REST
// @Summary Get public display data
// @Description Get attendance data for public display using display token
// @Tags Public Display
// @Produce json
// @Param token path string true "Display Token"
// @Success 200 {object} PublicDisplayData
// @Failure 401 {object} PublicDisplayError "Token invalid or revoked"
// @Failure 500 {object} PublicDisplayError "Internal server error"
// @Router /api/v1/public/display/{token} [get]
func (h *Handler) GetPublicDisplayData(c *fiber.Ctx) error {
	token := c.Params("token")
	if token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(PublicDisplayError{
			Error:   "TOKEN_REQUIRED",
			Message: "Token diperlukan",
		})
	}

	data, err := h.service.GetPublicDisplayData(c.Context(), token)
	if err != nil {
		// Check if it's an access error
		if accessErr, ok := err.(*PublicDisplayAccessError); ok {
			return c.Status(fiber.StatusUnauthorized).JSON(PublicDisplayError{
				Error:   accessErr.Code,
				Message: accessErr.Message,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(PublicDisplayError{
			Error:   "INTERNAL_ERROR",
			Message: "Gagal mengambil data display",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

// HandlePublicWebSocket handles WebSocket connections for public display real-time updates
// Requirements: 5.9 - WHEN a new attendance is recorded, THE System SHALL update the public display within 3 seconds
func (h *Handler) HandlePublicWebSocket(c *websocket.Conn) {
	token := c.Params("token")
	if token == "" {
		h.sendPublicWSError(c, "TOKEN_REQUIRED", "Token diperlukan")
		c.Close()
		return
	}

	// Validate token and get school ID
	validation, err := h.service.ValidateAndUpdateAccess(context.Background(), token)
	if err != nil || !validation.Valid {
		errorMsg := "Token tidak valid"
		if validation != nil && validation.Error != "" {
			errorMsg = validation.Error
		}
		h.sendPublicWSError(c, "TOKEN_INVALID", errorMsg)
		c.Close()
		return
	}

	// Create client for public display
	client := &realtime.Client{
		Hub:      h.realtimeHub,
		Conn:     c,
		Send:     make(chan []byte, 256),
		SchoolID: validation.SchoolID,
		IsPublic: true,
		Token:    token,
	}

	// Register client
	h.realtimeHub.Register(client)

	// Send connection success message
	h.sendPublicWSMessage(c, "connected", PublicWSConnectionStatus{
		Connected: true,
		Message:   "Terhubung ke server real-time",
	})

	// Start goroutines for reading and writing
	go h.writePublicPump(client)
	h.readPublicPump(client)
}

// readPublicPump pumps messages from the WebSocket connection
func (h *Handler) readPublicPump(client *realtime.Client) {
	defer func() {
		client.Hub.Unregister(client)
		client.Conn.Close()
	}()

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				// Log error if needed
			}
			break
		}

		// Handle incoming messages
		var msg PublicWSMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			continue
		}

		switch msg.Type {
		case "ping":
			h.sendPublicWSMessage(client.Conn, "pong", map[string]interface{}{
				"timestamp": time.Now().Unix(),
			})
		case "refresh":
			// Client requests a data refresh
			data, err := h.service.GetPublicDisplayData(context.Background(), client.Token)
			if err == nil {
				h.sendPublicWSMessage(client.Conn, "refresh_data", data)
			}
		}
	}
}

// writePublicPump pumps messages from the hub to the WebSocket connection
func (h *Handler) writePublicPump(client *realtime.Client) {
	ticker := time.NewTicker(30 * time.Second) // Ping interval
	defer func() {
		ticker.Stop()
		client.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.Send:
			if !ok {
				// Hub closed the channel
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// Transform the message for public display (remove sensitive data)
			transformedMsg := h.transformForPublicDisplay(message)
			if transformedMsg != nil {
				if err := client.Conn.WriteMessage(websocket.TextMessage, transformedMsg); err != nil {
					return
				}
			}
		case <-ticker.C:
			// Send ping to keep connection alive
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// transformForPublicDisplay transforms attendance event for public display
// Requirements: 5.14 - SHALL NOT expose sensitive student information (only name, class, and attendance time)
func (h *Handler) transformForPublicDisplay(message []byte) []byte {
	var event realtime.AttendanceEvent
	if err := json.Unmarshal(message, &event); err != nil {
		return nil
	}

	// Create public-safe event
	publicEvent := struct {
		Type        string                    `json:"type"`
		Attendance  *PublicLiveFeedEntry      `json:"attendance,omitempty"`
		Stats       *PublicAttendanceStats    `json:"stats,omitempty"`
		Leaderboard []PublicLeaderboardEntry  `json:"leaderboard,omitempty"`
	}{
		Type: string(event.Type),
	}

	// Transform attendance entry (remove StudentID)
	if event.Attendance != nil {
		publicEvent.Attendance = &PublicLiveFeedEntry{
			StudentName: event.Attendance.StudentName,
			ClassName:   event.Attendance.ClassName,
			Time:        event.Attendance.Time,
			Status:      string(event.Attendance.Status),
			Type:        event.Attendance.Type,
		}
	}

	// Transform stats
	if event.Stats != nil {
		publicEvent.Stats = &PublicAttendanceStats{
			TotalStudents: event.Stats.TotalStudents,
			Present:       event.Stats.Present,
			Late:          event.Stats.Late,
			VeryLate:      event.Stats.VeryLate,
			Absent:        event.Stats.Absent,
			Percentage:    event.Stats.Percentage,
		}
	}

	// Transform leaderboard (remove StudentID)
	if event.Leaderboard != nil {
		publicEvent.Leaderboard = make([]PublicLeaderboardEntry, len(event.Leaderboard))
		for i, entry := range event.Leaderboard {
			publicEvent.Leaderboard[i] = PublicLeaderboardEntry{
				Rank:        entry.Rank,
				StudentName: entry.StudentName,
				ClassName:   entry.ClassName,
				ArrivalTime: entry.ArrivalTime,
			}
		}
	}

	result, err := json.Marshal(publicEvent)
	if err != nil {
		return nil
	}

	return result
}

// sendPublicWSMessage sends a WebSocket message for public display
func (h *Handler) sendPublicWSMessage(c *websocket.Conn, msgType string, payload interface{}) {
	msg := PublicWSMessage{
		Type:    msgType,
		Payload: payload,
	}
	data, err := json.Marshal(msg)
	if err != nil {
		return
	}
	c.WriteMessage(websocket.TextMessage, data)
}

// sendPublicWSError sends a WebSocket error message for public display
func (h *Handler) sendPublicWSError(c *websocket.Conn, code string, message string) {
	h.sendPublicWSMessage(c, "error", PublicDisplayError{
		Error:   code,
		Message: message,
	})
}
