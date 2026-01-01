package auth

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/school-management/backend/internal/domain/models"
)

var (
	ErrUserNotFound = errors.New("user tidak ditemukan")
	ErrInvalidCredentials = errors.New("username atau password salah")
)

// Repository defines the interface for auth data operations
type Repository interface {
	FindByUsername(ctx context.Context, username string) (*models.User, error)
	FindByID(ctx context.Context, id uint) (*models.User, error)
	UpdatePassword(ctx context.Context, id uint, passwordHash string) error
	UpdateLastLogin(ctx context.Context, id uint) error
	ClearPasswordReset(ctx context.Context, id uint) error
}

// repository implements the Repository interface
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new auth repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// FindByUsername finds a user by username
// Requirements: 12.1 - WHEN a parent enters NISN and password, THE System SHALL authenticate
func (r *repository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).
		Preload("School").
		Where("username = ?", username).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

// FindByID finds a user by ID
func (r *repository) FindByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).
		Preload("School").
		Where("id = ?", id).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

// UpdatePassword updates the user's password hash
// Requirements: 12.5 - System SHALL enforce password reset on first login
func (r *repository) UpdatePassword(ctx context.Context, id uint, passwordHash string) error {
	result := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", id).
		Update("password_hash", passwordHash)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

// UpdateLastLogin updates the user's last login timestamp
func (r *repository) UpdateLastLogin(ctx context.Context, id uint) error {
	now := time.Now()
	result := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", id).
		Update("last_login_at", now)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// ClearPasswordReset clears the must_reset_pwd flag
// Requirements: 12.5 - After password change, clear the reset flag
func (r *repository) ClearPasswordReset(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", id).
		Update("must_reset_pwd", false)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
