package realtime

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"

	"github.com/school-management/backend/internal/modules/auth"
)

// Handler handles HTTP and WebSocket requests for real-time attendance
type Handler struct {
	service    Service
	jwtManager *auth.JWTManager
}

// NewHandler creates a new real-time handler
func NewHandler(service Service, jwtManager *auth.JWTManager) *Handler {
	return &Handler{
		service:    service,
		jwtManager: jwtManager,
	}
}

// RegisterRoutes registers real-time routes
func (h *Handler) RegisterRoutes(router fiber.Router) {
	// REST endpoints for initial data load - routes are already under /realtime group
	router.Get("/live-feed", h.GetLiveFeed)
	router.Get("/stats", h.GetStats)
	router.Get("/leaderboard", h.GetLeaderboard)
}

// RegisterWebSocketRoutes registers WebSocket routes
// Requirements: 4.8 - THE System SHALL use WebSocket for real-time updates
func (h *Handler) RegisterWebSocketRoutes(app *fiber.App) {
	// WebSocket upgrade middleware
	app.Use("/api/v1/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	// WebSocket endpoint for authenticated users
	// Requirements: 4.8, 4.9 - WebSocket with reconnection support
	app.Get("/api/v1/ws/attendance", websocket.New(h.HandleWebSocket, websocket.Config{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}))
}

// GetLiveFeed handles getting live attendance feed
// @Summary Get live attendance feed
// @Description Get the 20 most recent attendance records for today
// @Tags Real-Time
// @Produce json
// @Param class_id query int false "Filter by class ID"
// @Success 200 {object} LiveFeedResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/realtime/live-feed [get]
func (h *Handler) GetLiveFeed(c *fiber.Ctx) error {
	schoolID, ok := c.Locals("school_id").(uint)
	if !ok {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTHZ_TENANT_REQUIRED",
				"message": "Konteks sekolah diperlukan",
			},
		})
	}

	var classID *uint
	if classIDStr := c.Query("class_id"); classIDStr != "" {
		if id, err := strconv.ParseUint(classIDStr, 10, 32); err == nil {
			cid := uint(id)
			classID = &cid
		}
	}

	// Requirements: 4.6 - Wali_Kelas only sees their assigned class by default
	role, _ := c.Locals("role").(string)
	if role == "wali_kelas" && classID == nil {
		// Get assigned class for wali_kelas
		userID, _ := c.Locals("userID").(uint)
		if userID > 0 {
			// Note: This would need to be implemented via repository
			// For now, we'll allow them to see all if no class filter
		}
	}

	response, err := h.service.GetLiveFeed(c.Context(), schoolID, classID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "INTERNAL_ERROR",
				"message": "Gagal mengambil data live feed",
			},
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GetStats handles getting attendance statistics
// @Summary Get attendance statistics
// @Description Get current day's attendance statistics
// @Tags Real-Time
// @Produce json
// @Param class_id query int false "Filter by class ID"
// @Success 200 {object} StatsResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/realtime/stats [get]
func (h *Handler) GetStats(c *fiber.Ctx) error {
	schoolID, ok := c.Locals("school_id").(uint)
	if !ok {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTHZ_TENANT_REQUIRED",
				"message": "Konteks sekolah diperlukan",
			},
		})
	}

	var classID *uint
	if classIDStr := c.Query("class_id"); classIDStr != "" {
		if id, err := strconv.ParseUint(classIDStr, 10, 32); err == nil {
			cid := uint(id)
			classID = &cid
		}
	}

	response, err := h.service.GetAttendanceStats(c.Context(), schoolID, classID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "INTERNAL_ERROR",
				"message": "Gagal mengambil statistik kehadiran",
			},
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GetLeaderboard handles getting attendance leaderboard
// @Summary Get attendance leaderboard
// @Description Get top 10 earliest arrivals for today
// @Tags Real-Time
// @Produce json
// @Success 200 {object} LeaderboardResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/realtime/leaderboard [get]
func (h *Handler) GetLeaderboard(c *fiber.Ctx) error {
	schoolID, ok := c.Locals("school_id").(uint)
	if !ok {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTHZ_TENANT_REQUIRED",
				"message": "Konteks sekolah diperlukan",
			},
		})
	}

	response, err := h.service.GetLeaderboard(c.Context(), schoolID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "INTERNAL_ERROR",
				"message": "Gagal mengambil leaderboard",
			},
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// HandleWebSocket handles WebSocket connections for real-time updates
// Requirements: 4.8 - THE System SHALL use WebSocket for real-time updates
// Requirements: 4.9 - IF connection is lost, THE System SHALL attempt to reconnect automatically
func (h *Handler) HandleWebSocket(c *websocket.Conn) {
	// Get token from query parameter or header
	token := c.Query("token")
	if token == "" {
		// Try to get from Sec-WebSocket-Protocol header
		token = c.Headers("Sec-WebSocket-Protocol")
	}

	// Validate token
	if token == "" {
		h.sendWSError(c, "AUTH_TOKEN_MISSING", "Token diperlukan")
		c.Close()
		return
	}

	// Remove "Bearer " prefix if present
	token = strings.TrimPrefix(token, "Bearer ")
	token = strings.TrimPrefix(token, "bearer ")

	claims, err := h.jwtManager.ValidateAccessToken(token)
	if err != nil {
		h.sendWSError(c, "AUTH_TOKEN_INVALID", "Token tidak valid")
		c.Close()
		return
	}

	// Get school ID from claims
	if claims.SchoolID == nil {
		h.sendWSError(c, "AUTHZ_TENANT_REQUIRED", "Konteks sekolah diperlukan")
		c.Close()
		return
	}

	// Create client
	hub := h.service.GetHub()
	client := &Client{
		Hub:      hub,
		Conn:     c,
		Send:     make(chan []byte, 256),
		SchoolID: *claims.SchoolID,
		UserID:   claims.UserID,
		IsPublic: false,
	}

	// Register client
	hub.Register(client)

	// Send connection success message
	h.sendWSMessage(c, "connected", WSConnectionStatus{
		Connected: true,
		Message:   "Terhubung ke server real-time",
	})

	// Start goroutines for reading and writing
	go h.writePump(client)
	h.readPump(client)
}

// readPump pumps messages from the WebSocket connection to the hub
func (h *Handler) readPump(client *Client) {
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

		// Handle incoming messages (e.g., subscription requests)
		var msg WSMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			continue
		}

		switch msg.Type {
		case "subscribe":
			// Handle subscription with class filter
			if payload, ok := msg.Payload.(map[string]interface{}); ok {
				if classID, ok := payload["class_id"].(float64); ok {
					cid := uint(classID)
					client.ClassID = &cid
				}
			}
			h.sendWSMessage(client.Conn, "subscribed", map[string]interface{}{
				"class_id": client.ClassID,
			})
		case "ping":
			h.sendWSMessage(client.Conn, "pong", map[string]interface{}{
				"timestamp": time.Now().Unix(),
			})
		}
	}
}

// writePump pumps messages from the hub to the WebSocket connection
func (h *Handler) writePump(client *Client) {
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

			if err := client.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			// Send ping to keep connection alive
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// sendWSMessage sends a WebSocket message
func (h *Handler) sendWSMessage(c *websocket.Conn, msgType string, payload interface{}) {
	msg := WSMessage{
		Type:    msgType,
		Payload: payload,
	}
	data, err := json.Marshal(msg)
	if err != nil {
		return
	}
	c.WriteMessage(websocket.TextMessage, data)
}

// sendWSError sends a WebSocket error message
func (h *Handler) sendWSError(c *websocket.Conn, code string, message string) {
	h.sendWSMessage(c, "error", map[string]string{
		"code":    code,
		"message": message,
	})
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Success bool `json:"success"`
	Error   struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}
