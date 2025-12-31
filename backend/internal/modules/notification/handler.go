package notification

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/school-management/backend/internal/domain/models"
)

// Handler handles HTTP requests for notification management
// Requirements: 17.3, 17.4 - Notification CRUD and mark as read
type Handler struct {
	service Service
}

// NewHandler creates a new notification handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers notification routes
func (h *Handler) RegisterRoutes(router fiber.Router) {
	notifications := router.Group("/notifications")

	// Notification operations
	notifications.Get("", h.GetNotifications)
	notifications.Get("/summary", h.GetNotificationSummary)
	notifications.Get("/:id", h.GetNotificationByID)
	notifications.Post("/:id/read", h.MarkAsRead)
	notifications.Post("/read", h.MarkMultipleAsRead)
	notifications.Post("/read-all", h.MarkAllAsRead)
	notifications.Delete("/:id", h.DeleteNotification)

	// FCM Token operations
	fcm := router.Group("/fcm")
	fcm.Post("/tokens", h.RegisterFCMToken)
	fcm.Get("/tokens", h.GetFCMTokens)
	fcm.Delete("/tokens/:token", h.DeactivateFCMToken)
}

// ==================== Notification Handlers ====================

// GetNotifications handles listing notifications for the current user
// @Summary List notifications
// @Description Get a paginated list of notifications for the current user
// @Tags Notifications
// @Produce json
// @Param type query string false "Filter by notification type"
// @Param is_read query bool false "Filter by read status"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} NotificationListResponse
// @Failure 401 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/notifications [get]
func (h *Handler) GetNotifications(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return h.authRequiredError(c)
	}

	filter := h.parseNotificationFilter(c)

	response, err := h.service.GetUserNotifications(c.Context(), userID, filter)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GetNotificationSummary handles getting notification summary for the current user
// @Summary Get notification summary
// @Description Get notification summary (total and unread count) for the current user
// @Tags Notifications
// @Produce json
// @Success 200 {object} NotificationSummary
// @Failure 401 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/notifications/summary [get]
func (h *Handler) GetNotificationSummary(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return h.authRequiredError(c)
	}

	summary, err := h.service.GetNotificationSummary(c.Context(), userID)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    summary,
	})
}

// GetNotificationByID handles getting a single notification
// @Summary Get notification by ID
// @Description Get detailed information about a specific notification
// @Tags Notifications
// @Produce json
// @Param id path int true "Notification ID"
// @Success 200 {object} NotificationResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/notifications/{id} [get]
func (h *Handler) GetNotificationByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return h.invalidIDError(c, "notification")
	}

	response, err := h.service.GetNotificationByID(c.Context(), uint(id))
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}


// MarkAsRead handles marking a single notification as read
// @Summary Mark notification as read
// @Description Mark a specific notification as read
// @Tags Notifications
// @Produce json
// @Param id path int true "Notification ID"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/notifications/{id}/read [post]
func (h *Handler) MarkAsRead(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return h.invalidIDError(c, "notification")
	}

	if err := h.service.MarkAsRead(c.Context(), uint(id)); err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Notification marked as read",
	})
}

// MarkMultipleAsRead handles marking multiple notifications as read
// @Summary Mark multiple notifications as read
// @Description Mark multiple notifications as read
// @Tags Notifications
// @Accept json
// @Produce json
// @Param request body MarkAsReadRequest true "Notification IDs to mark as read"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/notifications/read [post]
func (h *Handler) MarkMultipleAsRead(c *fiber.Ctx) error {
	var req MarkAsReadRequest
	if err := c.BodyParser(&req); err != nil {
		return h.invalidBodyError(c)
	}

	if len(req.NotificationIDs) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "notification_ids is required",
			},
		})
	}

	if err := h.service.MarkMultipleAsRead(c.Context(), req.NotificationIDs); err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Notifications marked as read",
	})
}

// MarkAllAsRead handles marking all notifications as read for the current user
// @Summary Mark all notifications as read
// @Description Mark all notifications as read for the current user
// @Tags Notifications
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/notifications/read-all [post]
func (h *Handler) MarkAllAsRead(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return h.authRequiredError(c)
	}

	if err := h.service.MarkAllAsRead(c.Context(), userID); err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "All notifications marked as read",
	})
}

// DeleteNotification handles deleting a notification
// @Summary Delete notification
// @Description Delete a specific notification
// @Tags Notifications
// @Produce json
// @Param id path int true "Notification ID"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/notifications/{id} [delete]
func (h *Handler) DeleteNotification(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return h.invalidIDError(c, "notification")
	}

	if err := h.service.DeleteNotification(c.Context(), uint(id)); err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Notification deleted successfully",
	})
}

// ==================== FCM Token Handlers ====================

// RegisterFCMToken handles registering an FCM token for push notifications
// @Summary Register FCM token
// @Description Register or update an FCM token for push notifications
// @Tags FCM
// @Accept json
// @Produce json
// @Param request body RegisterFCMTokenRequest true "FCM token data"
// @Success 201 {object} FCMTokenResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/fcm/tokens [post]
func (h *Handler) RegisterFCMToken(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return h.authRequiredError(c)
	}

	var req RegisterFCMTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return h.invalidBodyError(c)
	}

	response, err := h.service.RegisterFCMToken(c.Context(), userID, req)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    response,
		"message": "FCM token registered successfully",
	})
}

// GetFCMTokens handles getting FCM tokens for the current user
// @Summary Get FCM tokens
// @Description Get all FCM tokens for the current user
// @Tags FCM
// @Produce json
// @Success 200 {object} []FCMTokenResponse
// @Failure 401 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/fcm/tokens [get]
func (h *Handler) GetFCMTokens(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return h.authRequiredError(c)
	}

	tokens, err := h.service.GetUserFCMTokens(c.Context(), userID)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    tokens,
	})
}

// DeactivateFCMToken handles deactivating an FCM token
// @Summary Deactivate FCM token
// @Description Deactivate a specific FCM token
// @Tags FCM
// @Produce json
// @Param token path string true "FCM token"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/fcm/tokens/{token} [delete]
func (h *Handler) DeactivateFCMToken(c *fiber.Ctx) error {
	token := c.Params("token")
	if token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "token is required",
			},
		})
	}

	if err := h.service.DeactivateFCMToken(c.Context(), token); err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "FCM token deactivated successfully",
	})
}

// ==================== Helper Methods ====================

func (h *Handler) parseNotificationFilter(c *fiber.Ctx) NotificationFilter {
	filter := NotificationFilter{
		Page:     c.QueryInt("page", 1),
		PageSize: c.QueryInt("page_size", 20),
	}

	if typeStr := c.Query("type"); typeStr != "" {
		notifType := models.NotificationType(typeStr)
		filter.Type = &notifType
	}

	if isReadStr := c.Query("is_read"); isReadStr != "" {
		isRead := isReadStr == "true"
		filter.IsRead = &isRead
	}

	return filter
}

func (h *Handler) authRequiredError(c *fiber.Ctx) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"success": false,
		"error": fiber.Map{
			"code":    "AUTH_REQUIRED",
			"message": "Authentication required",
		},
	})
}

func (h *Handler) invalidBodyError(c *fiber.Ctx) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"success": false,
		"error": fiber.Map{
			"code":    "VAL_INVALID_FORMAT",
			"message": "Invalid request body",
		},
	})
}

func (h *Handler) invalidIDError(c *fiber.Ctx, resource string) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"success": false,
		"error": fiber.Map{
			"code":    "VAL_INVALID_FORMAT",
			"message": "Invalid " + resource + " ID",
		},
	})
}

func (h *Handler) handleError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, ErrNotificationNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "NOT_FOUND_NOTIFICATION",
				"message": "Notification not found",
			},
		})
	case errors.Is(err, ErrFCMTokenNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "NOT_FOUND_FCM_TOKEN",
				"message": "FCM token not found",
			},
		})
	case errors.Is(err, ErrUserNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "NOT_FOUND_USER",
				"message": "User not found",
			},
		})
	case errors.Is(err, ErrUserIDRequired),
		errors.Is(err, ErrTypeRequired),
		errors.Is(err, ErrTitleRequired),
		errors.Is(err, ErrMessageRequired),
		errors.Is(err, ErrTokenRequired),
		errors.Is(err, ErrPlatformRequired),
		errors.Is(err, ErrInvalidPlatform),
		errors.Is(err, ErrNotificationIDsRequired):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": err.Error(),
			},
		})
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "INTERNAL_ERROR",
				"message": "An internal error occurred",
			},
		})
	}
}
