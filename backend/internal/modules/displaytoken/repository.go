package displaytoken

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/school-management/backend/internal/domain/models"
)

var (
	ErrTokenNotFound = errors.New("display token tidak ditemukan")
	ErrTokenInvalid  = errors.New("token tidak valid atau tidak ditemukan")
	ErrTokenRevoked  = errors.New("token telah dicabut")
	ErrTokenExpired  = errors.New("token telah kedaluwarsa")
)

// Repository defines the interface for display token data operations
type Repository interface {
	// CRUD operations
	Create(ctx context.Context, token *models.DisplayToken) error
	FindByID(ctx context.Context, schoolID, id uint) (*models.DisplayToken, error)
	FindAll(ctx context.Context, schoolID uint) ([]models.DisplayToken, error)
	Update(ctx context.Context, token *models.DisplayToken) error
	Delete(ctx context.Context, schoolID, id uint) error

	// Token validation operations
	// Requirements: 5.10, 5.13 - Token validation for access control
	FindByToken(ctx context.Context, token string) (*models.DisplayToken, error)

	// Requirements: 6.7 - Track last accessed timestamp
	UpdateLastAccessed(ctx context.Context, tokenID uint) error

	// Count tokens for a school
	CountBySchool(ctx context.Context, schoolID uint) (int64, error)
}

// repository implements the Repository interface
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new display token repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Create creates a new display token
// Requirements: 5.1, 6.2 - Token creation
func (r *repository) Create(ctx context.Context, token *models.DisplayToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

// FindByID retrieves a display token by ID for a specific school
func (r *repository) FindByID(ctx context.Context, schoolID, id uint) (*models.DisplayToken, error) {
	var token models.DisplayToken
	err := r.db.WithContext(ctx).
		Where("id = ? AND school_id = ?", id, schoolID).
		First(&token).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTokenNotFound
		}
		return nil, err
	}

	return &token, nil
}

// FindAll retrieves all display tokens for a school
// Requirements: 6.1 - List all tokens for their school
func (r *repository) FindAll(ctx context.Context, schoolID uint) ([]models.DisplayToken, error) {
	var tokens []models.DisplayToken
	err := r.db.WithContext(ctx).
		Where("school_id = ?", schoolID).
		Order("created_at DESC").
		Find(&tokens).Error

	return tokens, err
}

// Update updates an existing display token
func (r *repository) Update(ctx context.Context, token *models.DisplayToken) error {
	result := r.db.WithContext(ctx).
		Model(&models.DisplayToken{}).
		Where("id = ? AND school_id = ?", token.ID, token.SchoolID).
		Updates(map[string]interface{}{
			"name":       token.Name,
			"token":      token.Token,
			"is_active":  token.IsActive,
			"expires_at": token.ExpiresAt,
		})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrTokenNotFound
	}
	return nil
}

// Delete deletes a display token
// Requirements: 6.6 - Permanently remove token from the system
func (r *repository) Delete(ctx context.Context, schoolID, id uint) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND school_id = ?", id, schoolID).
		Delete(&models.DisplayToken{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrTokenNotFound
	}
	return nil
}

// FindByToken retrieves a display token by its token value
// Requirements: 5.10, 5.13 - Token validation for access control
func (r *repository) FindByToken(ctx context.Context, token string) (*models.DisplayToken, error) {
	var displayToken models.DisplayToken
	err := r.db.WithContext(ctx).
		Preload("School").
		Where("token = ?", token).
		First(&displayToken).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTokenInvalid
		}
		return nil, err
	}

	return &displayToken, nil
}

// UpdateLastAccessed updates the last accessed timestamp for a token
// Requirements: 6.7 - Track last accessed timestamp for each display token
func (r *repository) UpdateLastAccessed(ctx context.Context, tokenID uint) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&models.DisplayToken{}).
		Where("id = ?", tokenID).
		Update("last_accessed_at", now).Error
}

// CountBySchool counts the number of display tokens for a school
func (r *repository) CountBySchool(ctx context.Context, schoolID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.DisplayToken{}).
		Where("school_id = ?", schoolID).
		Count(&count).Error

	return count, err
}
