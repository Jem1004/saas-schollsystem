package device

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/school-management/backend/internal/domain/models"
)

var (
	ErrSchoolIDRequired    = errors.New("ID sekolah wajib diisi")
	ErrDeviceCodeRequired  = errors.New("kode perangkat wajib diisi")
	ErrDeviceInactive      = errors.New("perangkat sudah nonaktif")
	ErrDeviceActive        = errors.New("perangkat sudah aktif")
	ErrAPIKeyGeneration    = errors.New("gagal membuat API key")
)

// Service defines the interface for device business logic
type Service interface {
	RegisterDevice(ctx context.Context, req RegisterDeviceRequest) (*DeviceWithAPIKeyResponse, error)
	GetAllDevices(ctx context.Context, filter DeviceFilter) (*DeviceListResponse, error)
	GetDevicesGroupedBySchool(ctx context.Context) (*GroupedDevicesResponse, error)
	GetDeviceByID(ctx context.Context, id uint) (*DeviceResponse, error)
	GetDevicesBySchool(ctx context.Context, schoolID uint) ([]DeviceResponse, error)
	UpdateDevice(ctx context.Context, id uint, req UpdateDeviceRequest) (*DeviceResponse, error)
	ValidateAPIKey(ctx context.Context, apiKey string) (*APIKeyValidationResponse, error)
	RevokeAPIKey(ctx context.Context, id uint) (*RevokeAPIKeyResponse, error)
	RegenerateAPIKey(ctx context.Context, id uint) (*RegenerateAPIKeyResponse, error)
	DeleteDevice(ctx context.Context, id uint) error
}

// service implements the Service interface
type service struct {
	repo Repository
}

// NewService creates a new device service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// RegisterDevice registers a new device with a generated API key
// Requirements: 2.1 - WHEN a Super_Admin registers a new device, THE System SHALL generate a unique API key for that device
func (s *service) RegisterDevice(ctx context.Context, req RegisterDeviceRequest) (*DeviceWithAPIKeyResponse, error) {
	// Validate required fields
	if req.SchoolID == 0 {
		return nil, ErrSchoolIDRequired
	}
	
	deviceCode := strings.TrimSpace(req.DeviceCode)
	if deviceCode == "" {
		return nil, ErrDeviceCodeRequired
	}

	// Check for duplicate device code
	existing, err := s.repo.FindByDeviceCode(ctx, deviceCode)
	if err != nil && !errors.Is(err, ErrDeviceNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, ErrDuplicateDeviceCode
	}

	// Create device
	device := &models.Device{
		SchoolID:    req.SchoolID,
		DeviceCode:  deviceCode,
		Description: strings.TrimSpace(req.Description),
		IsActive:    true,
	}

	// Generate API key
	if err := device.GenerateAPIKey(); err != nil {
		return nil, ErrAPIKeyGeneration
	}

	if err := s.repo.Create(ctx, device); err != nil {
		return nil, err
	}

	// Return response with API key (only shown once)
	return &DeviceWithAPIKeyResponse{
		DeviceResponse: *toDeviceResponse(device),
		APIKey:         device.APIKey,
	}, nil
}


// GetAllDevices retrieves all devices with pagination
// Requirements: 2.5 - WHEN a Super_Admin views devices, THE System SHALL display device status, school assignment, and last activity
func (s *service) GetAllDevices(ctx context.Context, filter DeviceFilter) (*DeviceListResponse, error) {
	// Set defaults
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}

	devices, total, err := s.repo.FindAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Convert to response
	deviceResponses := make([]DeviceResponse, len(devices))
	for i, device := range devices {
		deviceResponses[i] = *toDeviceResponse(&device)
	}

	// Calculate total pages
	totalPages := int(total) / filter.PageSize
	if int(total)%filter.PageSize > 0 {
		totalPages++
	}

	return &DeviceListResponse{
		Devices: deviceResponses,
		Pagination: PaginationMeta{
			Page:       filter.Page,
			PageSize:   filter.PageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// GetDeviceByID retrieves a device by ID
func (s *service) GetDeviceByID(ctx context.Context, id uint) (*DeviceResponse, error) {
	device, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return toDeviceResponse(device), nil
}

// GetDevicesBySchool retrieves all devices for a school
func (s *service) GetDevicesBySchool(ctx context.Context, schoolID uint) ([]DeviceResponse, error) {
	devices, err := s.repo.FindBySchoolID(ctx, schoolID)
	if err != nil {
		return nil, err
	}

	responses := make([]DeviceResponse, len(devices))
	for i, device := range devices {
		responses[i] = *toDeviceResponse(&device)
	}

	return responses, nil
}

// UpdateDevice updates a device
func (s *service) UpdateDevice(ctx context.Context, id uint, req UpdateDeviceRequest) (*DeviceResponse, error) {
	// Get existing device
	device, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Description != nil {
		device.Description = strings.TrimSpace(*req.Description)
	}
	if req.IsActive != nil {
		device.IsActive = *req.IsActive
	}

	if err := s.repo.Update(ctx, device); err != nil {
		return nil, err
	}

	return toDeviceResponse(device), nil
}

// ValidateAPIKey validates an API key and returns device info
// Requirements: 2.2 - WHEN a device sends attendance data, THE System SHALL validate the API key before processing
// Requirements: 2.3 - IF an invalid API key is used, THEN THE System SHALL reject the request and log the attempt
func (s *service) ValidateAPIKey(ctx context.Context, apiKey string) (*APIKeyValidationResponse, error) {
	if apiKey == "" {
		log.Printf("API key validation failed: empty API key")
		return &APIKeyValidationResponse{
			Valid:   false,
			Message: "API key wajib diisi",
		}, nil
	}

	device, err := s.repo.FindByAPIKey(ctx, apiKey)
	if err != nil {
		if errors.Is(err, ErrInvalidAPIKey) {
			// Log invalid API key attempt (Requirements: 2.3)
			log.Printf("API key validation failed: invalid API key attempted")
			return &APIKeyValidationResponse{
				Valid:   false,
				Message: "API key tidak valid",
			}, nil
		}
		return nil, err
	}

	// Update last seen timestamp
	if err := s.repo.UpdateLastSeen(ctx, device.ID); err != nil {
		// Log but don't fail the validation
		log.Printf("Failed to update last seen for device %d: %v", device.ID, err)
	}

	return &APIKeyValidationResponse{
		Valid:    true,
		DeviceID: device.ID,
		SchoolID: device.SchoolID,
		Message:  "API key valid",
	}, nil
}

// RevokeAPIKey revokes a device's API key by deactivating the device
// Requirements: 2.4 - WHEN a Super_Admin revokes a device API key, THE System SHALL immediately invalidate that key
func (s *service) RevokeAPIKey(ctx context.Context, id uint) (*RevokeAPIKeyResponse, error) {
	// Get existing device
	device, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check if already inactive
	if !device.IsActive {
		return nil, ErrDeviceInactive
	}

	// Deactivate device
	if err := s.repo.Deactivate(ctx, id); err != nil {
		return nil, err
	}

	return &RevokeAPIKeyResponse{
		DeviceID:   device.ID,
		DeviceCode: device.DeviceCode,
		Message:    "API key berhasil dicabut. Perangkat tidak dapat lagi mengirim data kehadiran.",
	}, nil
}

// RegenerateAPIKey generates a new API key for a device
func (s *service) RegenerateAPIKey(ctx context.Context, id uint) (*RegenerateAPIKeyResponse, error) {
	// Get existing device
	device, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Generate new API key
	if err := device.GenerateAPIKey(); err != nil {
		return nil, ErrAPIKeyGeneration
	}

	// Reactivate device if it was inactive
	device.IsActive = true

	if err := s.repo.Update(ctx, device); err != nil {
		return nil, err
	}

	return &RegenerateAPIKeyResponse{
		DeviceID:   device.ID,
		DeviceCode: device.DeviceCode,
		APIKey:     device.APIKey,
		Message:    "API key berhasil dibuat ulang. Silakan perbarui perangkat dengan key baru.",
	}, nil
}

// DeleteDevice deletes a device
func (s *service) DeleteDevice(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

// GetDevicesGroupedBySchool retrieves all devices grouped by school
func (s *service) GetDevicesGroupedBySchool(ctx context.Context) (*GroupedDevicesResponse, error) {
	schoolDevices, err := s.repo.FindAllGroupedBySchool(ctx)
	if err != nil {
		return nil, err
	}

	total := 0
	for i := range schoolDevices {
		total += len(schoolDevices[i].Devices)
	}

	return &GroupedDevicesResponse{
		Schools: schoolDevices,
		Total:   total,
	}, nil
}

// toDeviceResponse converts a Device model to DeviceResponse DTO
func toDeviceResponse(device *models.Device) *DeviceResponse {
	response := &DeviceResponse{
		ID:          device.ID,
		SchoolID:    device.SchoolID,
		DeviceCode:  device.DeviceCode,
		Description: device.Description,
		IsActive:    device.IsActive,
		LastSeenAt:  device.LastSeenAt,
		CreatedAt:   device.CreatedAt,
		UpdatedAt:   device.UpdatedAt,
	}

	// Include school name if loaded
	if device.School.ID != 0 {
		response.SchoolName = device.School.Name
	}

	return response
}
