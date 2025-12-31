package notification

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/school-management/backend/internal/domain/models"
	"github.com/school-management/backend/internal/shared/redis"
)

var (
	ErrUserIDRequired    = errors.New("user_id is required")
	ErrTypeRequired      = errors.New("notification type is required")
	ErrTitleRequired     = errors.New("title is required")
	ErrMessageRequired   = errors.New("message is required")
	ErrTokenRequired     = errors.New("token is required")
	ErrPlatformRequired  = errors.New("platform is required")
	ErrInvalidPlatform   = errors.New("platform must be android or ios")
	ErrNotificationIDsRequired = errors.New("notification_ids is required")
)

// Service defines the interface for notification business logic
// Requirements: 17.3, 17.4 - Notification CRUD and mark as read
type Service interface {
	// Notification operations
	CreateNotification(ctx context.Context, req CreateNotificationRequest) (*NotificationResponse, error)
	GetNotificationByID(ctx context.Context, id uint) (*NotificationResponse, error)
	GetUserNotifications(ctx context.Context, userID uint, filter NotificationFilter) (*NotificationListResponse, error)
	MarkAsRead(ctx context.Context, notificationID uint) error
	MarkMultipleAsRead(ctx context.Context, ids []uint) error
	MarkAllAsRead(ctx context.Context, userID uint) error
	DeleteNotification(ctx context.Context, id uint) error
	GetUnreadCount(ctx context.Context, userID uint) (int64, error)
	GetNotificationSummary(ctx context.Context, userID uint) (*NotificationSummary, error)

	// Queue operations
	QueueNotification(ctx context.Context, notification *NotificationQueueItem) error

	// FCM Token operations
	RegisterFCMToken(ctx context.Context, userID uint, req RegisterFCMTokenRequest) (*FCMTokenResponse, error)
	GetUserFCMTokens(ctx context.Context, userID uint) ([]FCMTokenResponse, error)
	DeactivateFCMToken(ctx context.Context, token string) error
	DeactivateAllUserTokens(ctx context.Context, userID uint) error

	// Send notification (creates notification and queues for FCM)
	SendNotification(ctx context.Context, userID uint, notifType models.NotificationType, title, message string, data map[string]interface{}) error
}

// service implements the Service interface
type service struct {
	repo        Repository
	redisClient *redis.Client
}

// NewService creates a new notification service
func NewService(repo Repository, redisClient *redis.Client) Service {
	return &service{
		repo:        repo,
		redisClient: redisClient,
	}
}

// ==================== Notification Service ====================

// CreateNotification creates a new notification
// Requirements: 17.3 - THE System SHALL store notification history in database
func (s *service) CreateNotification(ctx context.Context, req CreateNotificationRequest) (*NotificationResponse, error) {
	// Validate required fields
	if req.UserID == 0 {
		return nil, ErrUserIDRequired
	}
	if !req.Type.IsValid() {
		return nil, ErrTypeRequired
	}
	if req.Title == "" {
		return nil, ErrTitleRequired
	}
	if req.Message == "" {
		return nil, ErrMessageRequired
	}

	// Verify user exists
	_, err := s.repo.FindUserByID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	notification := &models.Notification{
		UserID:  req.UserID,
		Type:    req.Type,
		Title:   req.Title,
		Message: req.Message,
		IsRead:  false,
	}

	// Set additional data if provided
	if req.Data != nil {
		if err := notification.SetData(req.Data); err != nil {
			return nil, err
		}
	}

	if err := s.repo.Create(ctx, notification); err != nil {
		return nil, err
	}

	return toNotificationResponse(notification), nil
}

// GetNotificationByID retrieves a notification by ID
func (s *service) GetNotificationByID(ctx context.Context, id uint) (*NotificationResponse, error) {
	notification, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toNotificationResponse(notification), nil
}


// GetUserNotifications retrieves notifications for a user with pagination
// Requirements: 17.4 - WHEN user views notifications, THE System SHALL display all notifications with read/unread status
func (s *service) GetUserNotifications(ctx context.Context, userID uint, filter NotificationFilter) (*NotificationListResponse, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}

	notifications, total, err := s.repo.FindByUserID(ctx, userID, filter)
	if err != nil {
		return nil, err
	}

	// Get unread count
	unreadCount, err := s.repo.GetUnreadCount(ctx, userID)
	if err != nil {
		return nil, err
	}

	responses := make([]NotificationResponse, len(notifications))
	for i, n := range notifications {
		responses[i] = *toNotificationResponse(&n)
	}

	totalPages := int(total) / filter.PageSize
	if int(total)%filter.PageSize > 0 {
		totalPages++
	}

	return &NotificationListResponse{
		Notifications: responses,
		Pagination: PaginationMeta{
			Page:       filter.Page,
			PageSize:   filter.PageSize,
			Total:      total,
			TotalPages: totalPages,
		},
		UnreadCount: unreadCount,
	}, nil
}

// MarkAsRead marks a notification as read
// Requirements: 17.4 - THE System SHALL display all notifications with read/unread status
func (s *service) MarkAsRead(ctx context.Context, notificationID uint) error {
	return s.repo.MarkAsRead(ctx, notificationID)
}

// MarkMultipleAsRead marks multiple notifications as read
func (s *service) MarkMultipleAsRead(ctx context.Context, ids []uint) error {
	if len(ids) == 0 {
		return ErrNotificationIDsRequired
	}
	return s.repo.MarkMultipleAsRead(ctx, ids)
}

// MarkAllAsRead marks all notifications for a user as read
func (s *service) MarkAllAsRead(ctx context.Context, userID uint) error {
	return s.repo.MarkAllAsRead(ctx, userID)
}

// DeleteNotification deletes a notification
func (s *service) DeleteNotification(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

// GetUnreadCount returns the count of unread notifications for a user
func (s *service) GetUnreadCount(ctx context.Context, userID uint) (int64, error) {
	return s.repo.GetUnreadCount(ctx, userID)
}

// GetNotificationSummary returns notification summary for a user
func (s *service) GetNotificationSummary(ctx context.Context, userID uint) (*NotificationSummary, error) {
	filter := NotificationFilter{Page: 1, PageSize: 1}
	_, total, err := s.repo.FindByUserID(ctx, userID, filter)
	if err != nil {
		return nil, err
	}

	unreadCount, err := s.repo.GetUnreadCount(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &NotificationSummary{
		TotalCount:  total,
		UnreadCount: unreadCount,
	}, nil
}

// ==================== Queue Operations ====================

// QueueNotification adds a notification to the Redis queue for processing
// Requirements: 17.1 - WHEN an event triggers notification, THE System SHALL queue the notification in Redis
func (s *service) QueueNotification(ctx context.Context, item *NotificationQueueItem) error {
	if s.redisClient == nil {
		return nil // Skip queueing if Redis is not available
	}
	return s.redisClient.Enqueue(ctx, redis.NotificationQueue, item)
}

// ==================== FCM Token Operations ====================

// RegisterFCMToken registers or updates an FCM token for a user
func (s *service) RegisterFCMToken(ctx context.Context, userID uint, req RegisterFCMTokenRequest) (*FCMTokenResponse, error) {
	// Validate required fields
	if req.Token == "" {
		return nil, ErrTokenRequired
	}
	if req.Platform == "" {
		return nil, ErrPlatformRequired
	}
	if req.Platform != "android" && req.Platform != "ios" {
		return nil, ErrInvalidPlatform
	}

	// Check if token already exists
	existingToken, err := s.repo.FindFCMTokenByToken(ctx, req.Token)
	if err == nil && existingToken != nil {
		// Token exists, update it
		existingToken.UserID = userID
		existingToken.Platform = req.Platform
		existingToken.IsActive = true
		existingToken.UpdatedAt = time.Now()

		if err := s.repo.UpdateFCMToken(ctx, existingToken); err != nil {
			return nil, err
		}
		return toFCMTokenResponse(existingToken), nil
	}

	// Create new token
	token := &models.FCMToken{
		UserID:   userID,
		Token:    req.Token,
		Platform: req.Platform,
		IsActive: true,
	}

	if err := s.repo.CreateFCMToken(ctx, token); err != nil {
		return nil, err
	}

	return toFCMTokenResponse(token), nil
}

// GetUserFCMTokens retrieves all FCM tokens for a user
func (s *service) GetUserFCMTokens(ctx context.Context, userID uint) ([]FCMTokenResponse, error) {
	tokens, err := s.repo.FindFCMTokenByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	responses := make([]FCMTokenResponse, len(tokens))
	for i, t := range tokens {
		responses[i] = *toFCMTokenResponse(&t)
	}
	return responses, nil
}

// DeactivateFCMToken deactivates an FCM token
func (s *service) DeactivateFCMToken(ctx context.Context, token string) error {
	return s.repo.DeactivateFCMToken(ctx, token)
}

// DeactivateAllUserTokens deactivates all FCM tokens for a user
func (s *service) DeactivateAllUserTokens(ctx context.Context, userID uint) error {
	return s.repo.DeactivateAllUserTokens(ctx, userID)
}

// ==================== Send Notification ====================

// SendNotification creates a notification and queues it for FCM delivery
// Requirements: 17.1, 17.2 - Queue notification and send via FCM
func (s *service) SendNotification(ctx context.Context, userID uint, notifType models.NotificationType, title, message string, data map[string]interface{}) error {
	// Create notification in database
	req := CreateNotificationRequest{
		UserID:  userID,
		Type:    notifType,
		Title:   title,
		Message: message,
		Data:    data,
	}

	notification, err := s.CreateNotification(ctx, req)
	if err != nil {
		return err
	}

	// Queue for FCM delivery
	queueItem := &NotificationQueueItem{
		NotificationID: notification.ID,
		UserID:         userID,
		Type:           notifType,
		Title:          title,
		Message:        message,
		Data:           data,
		RetryCount:     0,
		CreatedAt:      time.Now(),
	}

	return s.QueueNotification(ctx, queueItem)
}

// ==================== Response Converters ====================

func toNotificationResponse(n *models.Notification) *NotificationResponse {
	response := &NotificationResponse{
		ID:        n.ID,
		UserID:    n.UserID,
		Type:      n.Type,
		Title:     n.Title,
		Message:   n.Message,
		IsRead:    n.IsRead,
		CreatedAt: n.CreatedAt,
	}

	// Parse data if present
	if n.Data != "" {
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(n.Data), &data); err == nil {
			response.Data = data
		}
	}

	return response
}

func toFCMTokenResponse(t *models.FCMToken) *FCMTokenResponse {
	return &FCMTokenResponse{
		ID:        t.ID,
		UserID:    t.UserID,
		Token:     t.Token,
		Platform:  t.Platform,
		IsActive:  t.IsActive,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}
