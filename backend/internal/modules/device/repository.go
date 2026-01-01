package device

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/school-management/backend/internal/domain/models"
)

var (
	ErrDeviceNotFound      = errors.New("perangkat tidak ditemukan")
	ErrDuplicateDeviceCode = errors.New("perangkat dengan kode ini sudah terdaftar")
	ErrDuplicateAPIKey     = errors.New("API key sudah ada")
	ErrInvalidAPIKey       = errors.New("API key tidak valid")
)

// Repository defines the interface for device data operations
type Repository interface {
	Create(ctx context.Context, device *models.Device) error
	FindAll(ctx context.Context, filter DeviceFilter) ([]models.Device, int64, error)
	FindAllGroupedBySchool(ctx context.Context) ([]SchoolDevicesResponse, error)
	FindByID(ctx context.Context, id uint) (*models.Device, error)
	FindByDeviceCode(ctx context.Context, code string) (*models.Device, error)
	FindByAPIKey(ctx context.Context, apiKey string) (*models.Device, error)
	FindBySchoolID(ctx context.Context, schoolID uint) ([]models.Device, error)
	Update(ctx context.Context, device *models.Device) error
	UpdateLastSeen(ctx context.Context, id uint) error
	Deactivate(ctx context.Context, id uint) error
	Delete(ctx context.Context, id uint) error
}

// repository implements the Repository interface
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new device repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Create creates a new device
// Requirements: 2.1 - WHEN a Super_Admin registers a new device, THE System SHALL generate a unique API key
func (r *repository) Create(ctx context.Context, device *models.Device) error {
	result := r.db.WithContext(ctx).Create(device)
	if result.Error != nil {
		// Check for unique constraint violations
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return ErrDuplicateDeviceCode
		}
		return result.Error
	}
	return nil
}


// FindAll retrieves all devices with pagination and filtering
// Requirements: 2.5 - WHEN a Super_Admin views devices, THE System SHALL display device status, school assignment, and last activity
func (r *repository) FindAll(ctx context.Context, filter DeviceFilter) ([]models.Device, int64, error) {
	var devices []models.Device
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Device{}).Preload("School")

	// Apply filters
	if filter.SchoolID != nil {
		query = query.Where("school_id = ?", *filter.SchoolID)
	}
	if filter.DeviceCode != "" {
		query = query.Where("device_code ILIKE ?", "%"+filter.DeviceCode+"%")
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
		Find(&devices).Error

	if err != nil {
		return nil, 0, err
	}

	return devices, total, nil
}

// FindByID retrieves a device by ID
func (r *repository) FindByID(ctx context.Context, id uint) (*models.Device, error) {
	var device models.Device
	err := r.db.WithContext(ctx).
		Preload("School").
		Where("id = ?", id).
		First(&device).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDeviceNotFound
		}
		return nil, err
	}

	return &device, nil
}

// FindByDeviceCode retrieves a device by device code
func (r *repository) FindByDeviceCode(ctx context.Context, code string) (*models.Device, error) {
	var device models.Device
	err := r.db.WithContext(ctx).
		Preload("School").
		Where("device_code = ?", code).
		First(&device).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDeviceNotFound
		}
		return nil, err
	}

	return &device, nil
}

// FindByAPIKey retrieves a device by API key
// Requirements: 2.2 - WHEN a device sends attendance data, THE System SHALL validate the API key before processing
func (r *repository) FindByAPIKey(ctx context.Context, apiKey string) (*models.Device, error) {
	var device models.Device
	err := r.db.WithContext(ctx).
		Preload("School").
		Where("api_key = ? AND is_active = ?", apiKey, true).
		First(&device).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidAPIKey
		}
		return nil, err
	}

	return &device, nil
}

// FindBySchoolID retrieves all devices for a school
func (r *repository) FindBySchoolID(ctx context.Context, schoolID uint) ([]models.Device, error) {
	var devices []models.Device
	err := r.db.WithContext(ctx).
		Where("school_id = ?", schoolID).
		Order("created_at DESC").
		Find(&devices).Error

	if err != nil {
		return nil, err
	}

	return devices, nil
}

// Update updates a device
func (r *repository) Update(ctx context.Context, device *models.Device) error {
	result := r.db.WithContext(ctx).
		Model(&models.Device{}).
		Where("id = ?", device.ID).
		Updates(map[string]interface{}{
			"school_id":   device.SchoolID,
			"device_code": device.DeviceCode,
			"description": device.Description,
			"is_active":   device.IsActive,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrDeviceNotFound
	}
	return nil
}

// UpdateLastSeen updates the last seen timestamp for a device
func (r *repository) UpdateLastSeen(ctx context.Context, id uint) error {
	now := time.Now()
	result := r.db.WithContext(ctx).
		Model(&models.Device{}).
		Where("id = ?", id).
		Update("last_seen_at", now)

	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Deactivate deactivates a device (revokes API key)
// Requirements: 2.4 - WHEN a Super_Admin revokes a device API key, THE System SHALL immediately invalidate that key
func (r *repository) Deactivate(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).
		Model(&models.Device{}).
		Where("id = ?", id).
		Update("is_active", false)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrDeviceNotFound
	}
	return nil
}

// Delete deletes a device
func (r *repository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).
		Delete(&models.Device{}, id)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrDeviceNotFound
	}
	return nil
}

// FindAllGroupedBySchool retrieves all devices grouped by school
func (r *repository) FindAllGroupedBySchool(ctx context.Context) ([]SchoolDevicesResponse, error) {
	// Get all schools that have devices
	var schools []models.School
	err := r.db.WithContext(ctx).
		Joins("JOIN devices ON devices.school_id = schools.id").
		Group("schools.id").
		Order("schools.name ASC").
		Find(&schools).Error
	if err != nil {
		return nil, err
	}

	// Get devices for each school
	result := make([]SchoolDevicesResponse, 0, len(schools))
	for _, school := range schools {
		var devices []models.Device
		err := r.db.WithContext(ctx).
			Where("school_id = ?", school.ID).
			Order("created_at DESC").
			Find(&devices).Error
		if err != nil {
			return nil, err
		}

		deviceResponses := make([]DeviceResponse, len(devices))
		for i, device := range devices {
			deviceResponses[i] = DeviceResponse{
				ID:          device.ID,
				SchoolID:    device.SchoolID,
				SchoolName:  school.Name,
				DeviceCode:  device.DeviceCode,
				Description: device.Description,
				IsActive:    device.IsActive,
				LastSeenAt:  device.LastSeenAt,
				CreatedAt:   device.CreatedAt,
				UpdatedAt:   device.UpdatedAt,
			}
		}

		result = append(result, SchoolDevicesResponse{
			SchoolID:    school.ID,
			SchoolName:  school.Name,
			IsActive:    school.IsActive,
			DeviceCount: len(devices),
			Devices:     deviceResponses,
		})
	}

	return result, nil
}
