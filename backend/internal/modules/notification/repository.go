package notification

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/school-management/backend/internal/domain/models"
)

var (
	ErrNotificationNotFound = errors.New("notifikasi tidak ditemukan")
	ErrFCMTokenNotFound     = errors.New("token FCM tidak ditemukan")
	ErrUserNotFound         = errors.New("user tidak ditemukan")
)

// Repository defines the interface for notification data operations
// Requirements: 17.3 - THE System SHALL store notification history in database with read status
type Repository interface {
	// Notification operations
	Create(ctx context.Context, notification *models.Notification) error
	FindByID(ctx context.Context, id uint) (*models.Notification, error)
	FindByUserID(ctx context.Context, userID uint, filter NotificationFilter) ([]models.Notification, int64, error)
	MarkAsRead(ctx context.Context, id uint) error
	MarkMultipleAsRead(ctx context.Context, ids []uint) error
	MarkAllAsRead(ctx context.Context, userID uint) error
	Delete(ctx context.Context, id uint) error
	GetUnreadCount(ctx context.Context, userID uint) (int64, error)

	// FCM Token operations
	CreateFCMToken(ctx context.Context, token *models.FCMToken) error
	FindFCMTokenByUserID(ctx context.Context, userID uint) ([]models.FCMToken, error)
	FindActiveFCMTokensByUserID(ctx context.Context, userID uint) ([]models.FCMToken, error)
	UpdateFCMToken(ctx context.Context, token *models.FCMToken) error
	DeactivateFCMToken(ctx context.Context, token string) error
	DeactivateAllUserTokens(ctx context.Context, userID uint) error
	FindFCMTokenByToken(ctx context.Context, token string) (*models.FCMToken, error)

	// User lookup
	FindUserByID(ctx context.Context, userID uint) (*models.User, error)
}

// repository implements the Repository interface
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new notification repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// ==================== Notification Repository ====================

// Create creates a new notification record
// Requirements: 17.3 - THE System SHALL store notification history in database
func (r *repository) Create(ctx context.Context, notification *models.Notification) error {
	return r.db.WithContext(ctx).Create(notification).Error
}

// FindByID retrieves a notification by ID
func (r *repository) FindByID(ctx context.Context, id uint) (*models.Notification, error) {
	var notification models.Notification
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&notification).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotificationNotFound
		}
		return nil, err
	}
	return &notification, nil
}


// FindByUserID retrieves notifications for a user with pagination and filtering
// Requirements: 17.4 - WHEN user views notifications, THE System SHALL display all notifications with read/unread status
func (r *repository) FindByUserID(ctx context.Context, userID uint, filter NotificationFilter) ([]models.Notification, int64, error) {
	var notifications []models.Notification
	var total int64

	query := r.db.WithContext(ctx).
		Model(&models.Notification{}).
		Where("user_id = ?", userID)

	// Apply filters
	if filter.Type != nil {
		query = query.Where("type = ?", *filter.Type)
	}
	if filter.IsRead != nil {
		query = query.Where("is_read = ?", *filter.IsRead)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}
	if filter.PageSize > 100 {
		filter.PageSize = 100
	}

	offset := (filter.Page - 1) * filter.PageSize

	// Fetch records ordered by created_at DESC (most recent first)
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(filter.PageSize).
		Find(&notifications).Error

	if err != nil {
		return nil, 0, err
	}

	return notifications, total, nil
}

// MarkAsRead marks a notification as read
// Requirements: 17.4 - THE System SHALL display all notifications with read/unread status
func (r *repository) MarkAsRead(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).
		Model(&models.Notification{}).
		Where("id = ?", id).
		Update("is_read", true)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotificationNotFound
	}
	return nil
}

// MarkMultipleAsRead marks multiple notifications as read
func (r *repository) MarkMultipleAsRead(ctx context.Context, ids []uint) error {
	if len(ids) == 0 {
		return nil
	}

	return r.db.WithContext(ctx).
		Model(&models.Notification{}).
		Where("id IN ?", ids).
		Update("is_read", true).Error
}

// MarkAllAsRead marks all notifications for a user as read
func (r *repository) MarkAllAsRead(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).
		Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Update("is_read", true).Error
}

// Delete deletes a notification record
func (r *repository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.Notification{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotificationNotFound
	}
	return nil
}

// GetUnreadCount returns the count of unread notifications for a user
func (r *repository) GetUnreadCount(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&count).Error
	return count, err
}

// ==================== FCM Token Repository ====================

// CreateFCMToken creates a new FCM token record
func (r *repository) CreateFCMToken(ctx context.Context, token *models.FCMToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

// FindFCMTokenByUserID retrieves all FCM tokens for a user
func (r *repository) FindFCMTokenByUserID(ctx context.Context, userID uint) ([]models.FCMToken, error) {
	var tokens []models.FCMToken
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&tokens).Error
	return tokens, err
}

// FindActiveFCMTokensByUserID retrieves active FCM tokens for a user
func (r *repository) FindActiveFCMTokensByUserID(ctx context.Context, userID uint) ([]models.FCMToken, error) {
	var tokens []models.FCMToken
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_active = ?", userID, true).
		Find(&tokens).Error
	return tokens, err
}

// UpdateFCMToken updates an FCM token record
func (r *repository) UpdateFCMToken(ctx context.Context, token *models.FCMToken) error {
	return r.db.WithContext(ctx).Save(token).Error
}

// DeactivateFCMToken deactivates an FCM token by token string
func (r *repository) DeactivateFCMToken(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).
		Model(&models.FCMToken{}).
		Where("token = ?", token).
		Update("is_active", false).Error
}

// DeactivateAllUserTokens deactivates all FCM tokens for a user
func (r *repository) DeactivateAllUserTokens(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).
		Model(&models.FCMToken{}).
		Where("user_id = ?", userID).
		Update("is_active", false).Error
}

// FindFCMTokenByToken retrieves an FCM token by token string
func (r *repository) FindFCMTokenByToken(ctx context.Context, token string) (*models.FCMToken, error) {
	var fcmToken models.FCMToken
	err := r.db.WithContext(ctx).
		Where("token = ?", token).
		First(&fcmToken).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrFCMTokenNotFound
		}
		return nil, err
	}
	return &fcmToken, nil
}

// ==================== User Lookup ====================

// FindUserByID retrieves a user by ID
func (r *repository) FindUserByID(ctx context.Context, userID uint) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).
		Where("id = ?", userID).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}
