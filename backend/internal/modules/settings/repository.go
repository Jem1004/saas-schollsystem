package settings

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/school-management/backend/internal/domain/models"
)

var (
	ErrSettingsNotFound = errors.New("pengaturan tidak ditemukan")
	ErrSchoolNotFound   = errors.New("sekolah tidak ditemukan")
)

// Repository defines the interface for SchoolSettings data operations
type Repository interface {
	// Settings operations
	FindBySchoolID(ctx context.Context, schoolID uint) (*models.SchoolSettings, error)
	Create(ctx context.Context, settings *models.SchoolSettings) error
	Update(ctx context.Context, settings *models.SchoolSettings) error
	Delete(ctx context.Context, schoolID uint) error

	// School lookup
	FindSchoolByID(ctx context.Context, schoolID uint) (*models.School, error)
}

// repository implements the Repository interface
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new Settings repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// FindBySchoolID retrieves settings for a specific school
func (r *repository) FindBySchoolID(ctx context.Context, schoolID uint) (*models.SchoolSettings, error) {
	var settings models.SchoolSettings
	err := r.db.WithContext(ctx).
		Where("school_id = ?", schoolID).
		First(&settings).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSettingsNotFound
		}
		return nil, err
	}
	return &settings, nil
}

// Create creates new school settings
func (r *repository) Create(ctx context.Context, settings *models.SchoolSettings) error {
	return r.db.WithContext(ctx).Create(settings).Error
}

// Update updates existing school settings
func (r *repository) Update(ctx context.Context, settings *models.SchoolSettings) error {
	result := r.db.WithContext(ctx).Save(settings)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrSettingsNotFound
	}
	return nil
}

// Delete deletes school settings
func (r *repository) Delete(ctx context.Context, schoolID uint) error {
	result := r.db.WithContext(ctx).
		Where("school_id = ?", schoolID).
		Delete(&models.SchoolSettings{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// FindSchoolByID retrieves a school by ID
func (r *repository) FindSchoolByID(ctx context.Context, schoolID uint) (*models.School, error) {
	var school models.School
	err := r.db.WithContext(ctx).
		Where("id = ?", schoolID).
		First(&school).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSchoolNotFound
		}
		return nil, err
	}
	return &school, nil
}
