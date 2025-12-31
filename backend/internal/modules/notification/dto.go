package notification

import (
	"time"

	"github.com/school-management/backend/internal/domain/models"
)

// ==================== Pagination ====================

// PaginationMeta represents pagination metadata
type PaginationMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// ==================== Notification DTOs ====================

// CreateNotificationRequest represents the request to create a notification
// Requirements: 17.3 - THE System SHALL store notification history in database
type CreateNotificationRequest struct {
	UserID  uint                    `json:"user_id" validate:"required"`
	Type    models.NotificationType `json:"type" validate:"required"`
	Title   string                  `json:"title" validate:"required"`
	Message string                  `json:"message" validate:"required"`
	Data    map[string]interface{}  `json:"data,omitempty"`
}

// NotificationResponse represents a notification in responses
// Requirements: 17.4 - THE System SHALL display all notifications with read/unread status
type NotificationResponse struct {
	ID        uint                    `json:"id"`
	UserID    uint                    `json:"user_id"`
	Type      models.NotificationType `json:"type"`
	Title     string                  `json:"title"`
	Message   string                  `json:"message"`
	Data      map[string]interface{}  `json:"data,omitempty"`
	IsRead    bool                    `json:"is_read"`
	CreatedAt time.Time               `json:"created_at"`
}

// NotificationListResponse represents a paginated list of notifications
type NotificationListResponse struct {
	Notifications []NotificationResponse `json:"notifications"`
	Pagination    PaginationMeta         `json:"pagination"`
	UnreadCount   int64                  `json:"unread_count"`
}

// NotificationFilter represents filter options for listing notifications
type NotificationFilter struct {
	Type     *models.NotificationType `query:"type"`
	IsRead   *bool                    `query:"is_read"`
	Page     int                      `query:"page"`
	PageSize int                      `query:"page_size"`
}

// MarkAsReadRequest represents the request to mark notifications as read
type MarkAsReadRequest struct {
	NotificationIDs []uint `json:"notification_ids" validate:"required"`
}

// ==================== FCM Token DTOs ====================

// RegisterFCMTokenRequest represents the request to register an FCM token
type RegisterFCMTokenRequest struct {
	Token    string `json:"token" validate:"required"`
	Platform string `json:"platform" validate:"required,oneof=android ios"`
}

// FCMTokenResponse represents an FCM token in responses
type FCMTokenResponse struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Token     string    `json:"token"`
	Platform  string    `json:"platform"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ==================== Queue DTOs ====================

// NotificationQueueItem represents a notification in the queue
// Requirements: 17.1 - THE System SHALL queue the notification in Redis
type NotificationQueueItem struct {
	NotificationID uint                    `json:"notification_id"`
	UserID         uint                    `json:"user_id"`
	Type           models.NotificationType `json:"type"`
	Title          string                  `json:"title"`
	Message        string                  `json:"message"`
	Data           map[string]interface{}  `json:"data,omitempty"`
	RetryCount     int                     `json:"retry_count"`
	CreatedAt      time.Time               `json:"created_at"`
}

// NotificationSummary represents notification summary for a user
type NotificationSummary struct {
	TotalCount  int64 `json:"total_count"`
	UnreadCount int64 `json:"unread_count"`
}
