package device

import "time"

// RegisterDeviceRequest represents the request to register a new device
// Requirements: 2.1 - WHEN a Super_Admin registers a new device, THE System SHALL generate a unique API key
type RegisterDeviceRequest struct {
	SchoolID    uint   `json:"school_id" validate:"required"`
	DeviceCode  string `json:"device_code" validate:"required"`
	Description string `json:"description"`
}

// UpdateDeviceRequest represents the request to update a device
type UpdateDeviceRequest struct {
	Description *string `json:"description"`
	IsActive    *bool   `json:"is_active"`
}

// DeviceResponse represents the device data in responses
// Requirements: 2.5 - WHEN a Super_Admin views devices, THE System SHALL display device status, school assignment, and last activity
type DeviceResponse struct {
	ID          uint       `json:"id"`
	SchoolID    uint       `json:"school_id"`
	SchoolName  string     `json:"school_name,omitempty"`
	DeviceCode  string     `json:"device_code"`
	Description string     `json:"description"`
	IsActive    bool       `json:"is_active"`
	LastSeenAt  *time.Time `json:"last_seen_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// DeviceWithAPIKeyResponse includes the API key (only returned on creation/regeneration)
type DeviceWithAPIKeyResponse struct {
	DeviceResponse
	APIKey string `json:"api_key"`
}

// DeviceListResponse represents a paginated list of devices
type DeviceListResponse struct {
	Devices    []DeviceResponse `json:"devices"`
	Pagination PaginationMeta   `json:"pagination"`
}

// SchoolDevicesResponse represents devices grouped by school
type SchoolDevicesResponse struct {
	SchoolID    uint             `json:"school_id"`
	SchoolName  string           `json:"school_name"`
	IsActive    bool             `json:"is_active"`
	DeviceCount int              `json:"device_count"`
	Devices     []DeviceResponse `json:"devices"`
}

// GroupedDevicesResponse represents all devices grouped by school
type GroupedDevicesResponse struct {
	Schools []SchoolDevicesResponse `json:"schools"`
	Total   int                     `json:"total"`
}

// PaginationMeta represents pagination metadata
type PaginationMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// DeviceFilter represents filter options for listing devices
type DeviceFilter struct {
	SchoolID   *uint  `query:"school_id"`
	DeviceCode string `query:"device_code"`
	IsActive   *bool  `query:"is_active"`
	Page       int    `query:"page"`
	PageSize   int    `query:"page_size"`
}

// DefaultDeviceFilter returns default filter values
func DefaultDeviceFilter() DeviceFilter {
	return DeviceFilter{
		Page:     1,
		PageSize: 20,
	}
}

// APIKeyValidationResponse represents the response for API key validation
type APIKeyValidationResponse struct {
	Valid    bool   `json:"valid"`
	DeviceID uint   `json:"device_id,omitempty"`
	SchoolID uint   `json:"school_id,omitempty"`
	Message  string `json:"message,omitempty"`
}

// RegenerateAPIKeyResponse represents the response for API key regeneration
type RegenerateAPIKeyResponse struct {
	DeviceID   uint   `json:"device_id"`
	DeviceCode string `json:"device_code"`
	APIKey     string `json:"api_key"`
	Message    string `json:"message"`
}

// RevokeAPIKeyResponse represents the response for API key revocation
type RevokeAPIKeyResponse struct {
	DeviceID   uint   `json:"device_id"`
	DeviceCode string `json:"device_code"`
	Message    string `json:"message"`
}
