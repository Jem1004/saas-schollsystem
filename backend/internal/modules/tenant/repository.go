package tenant

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/school-management/backend/internal/domain/models"
)

var (
	ErrSchoolNotFound   = errors.New("school not found")
	ErrDuplicateSchool  = errors.New("school with this name already exists")
)

// Repository defines the interface for tenant data operations
type Repository interface {
	Create(ctx context.Context, school *models.School) error
	FindAll(ctx context.Context, filter SchoolFilter) ([]models.School, int64, error)
	FindByID(ctx context.Context, id uint) (*models.School, error)
	FindByName(ctx context.Context, name string) (*models.School, error)
	Update(ctx context.Context, school *models.School) error
	Deactivate(ctx context.Context, id uint) error
	Activate(ctx context.Context, id uint) error
	GetStats(ctx context.Context, schoolID uint) (*SchoolStats, error)
}

// repository implements the Repository interface
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new tenant repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Create creates a new school (tenant)
// Requirements: 1.1 - WHEN a Super_Admin creates a new tenant, THE System SHALL generate a unique school_id
func (r *repository) Create(ctx context.Context, school *models.School) error {
	result := r.db.WithContext(ctx).Create(school)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// FindAll retrieves all schools with pagination and filtering
// Requirements: 1.2 - WHEN a Super_Admin views the tenant list, THE System SHALL display all registered schools
func (r *repository) FindAll(ctx context.Context, filter SchoolFilter) ([]models.School, int64, error) {
	var schools []models.School
	var total int64

	query := r.db.WithContext(ctx).Model(&models.School{})

	// Apply filters
	if filter.Name != "" {
		query = query.Where("name ILIKE ?", "%"+filter.Name+"%")
	}
	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (filter.Page - 1) * filter.PageSize
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}
	if filter.PageSize > 100 {
		filter.PageSize = 100
	}

	// Fetch records
	err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(filter.PageSize).
		Find(&schools).Error

	if err != nil {
		return nil, 0, err
	}

	return schools, total, nil
}

// FindByID retrieves a school by ID
func (r *repository) FindByID(ctx context.Context, id uint) (*models.School, error) {
	var school models.School
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&school).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSchoolNotFound
		}
		return nil, err
	}

	return &school, nil
}

// FindByName retrieves a school by name
func (r *repository) FindByName(ctx context.Context, name string) (*models.School, error) {
	var school models.School
	err := r.db.WithContext(ctx).
		Where("name = ?", name).
		First(&school).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSchoolNotFound
		}
		return nil, err
	}

	return &school, nil
}

// Update updates a school
func (r *repository) Update(ctx context.Context, school *models.School) error {
	result := r.db.WithContext(ctx).Save(school)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrSchoolNotFound
	}
	return nil
}

// Deactivate deactivates a school
// Requirements: 1.3 - WHEN a Super_Admin deactivates a tenant, THE System SHALL prevent all users from accessing
func (r *repository) Deactivate(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).
		Model(&models.School{}).
		Where("id = ?", id).
		Update("is_active", false)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrSchoolNotFound
	}
	return nil
}

// Activate activates a school
func (r *repository) Activate(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).
		Model(&models.School{}).
		Where("id = ?", id).
		Update("is_active", true)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrSchoolNotFound
	}
	return nil
}

// GetStats retrieves statistics for a school
func (r *repository) GetStats(ctx context.Context, schoolID uint) (*SchoolStats, error) {
	stats := &SchoolStats{}

	// Count classes
	if err := r.db.WithContext(ctx).
		Model(&models.Class{}).
		Where("school_id = ?", schoolID).
		Count(&stats.TotalClasses).Error; err != nil {
		return nil, err
	}

	// Count students
	if err := r.db.WithContext(ctx).
		Model(&models.Student{}).
		Where("school_id = ?", schoolID).
		Count(&stats.TotalStudents).Error; err != nil {
		return nil, err
	}

	// Count users
	if err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("school_id = ?", schoolID).
		Count(&stats.TotalUsers).Error; err != nil {
		return nil, err
	}

	// Count devices
	if err := r.db.WithContext(ctx).
		Model(&models.Device{}).
		Where("school_id = ?", schoolID).
		Count(&stats.TotalDevices).Error; err != nil {
		return nil, err
	}

	return stats, nil
}
