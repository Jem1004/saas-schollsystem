package tenant

import "time"

// CreateSchoolRequest represents the request to create a new school (tenant)
// Requirements: 1.1 - WHEN a Super_Admin creates a new tenant, THE System SHALL generate a unique school_id
type CreateSchoolRequest struct {
	Name    string `json:"name" validate:"required"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
}

// UpdateSchoolRequest represents the request to update a school
type UpdateSchoolRequest struct {
	Name    *string `json:"name"`
	Address *string `json:"address"`
	Phone   *string `json:"phone"`
	Email   *string `json:"email"`
}

// SchoolResponse represents the school data in responses
// Requirements: 1.2 - WHEN a Super_Admin views the tenant list, THE System SHALL display all registered schools
type SchoolResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Stats     *SchoolStats `json:"stats,omitempty"`
}

// SchoolStats represents statistics for a school
type SchoolStats struct {
	TotalClasses  int64 `json:"total_classes"`
	TotalStudents int64 `json:"total_students"`
	TotalUsers    int64 `json:"total_users"`
	TotalDevices  int64 `json:"total_devices"`
}

// SchoolListResponse represents a paginated list of schools
type SchoolListResponse struct {
	Schools    []SchoolResponse `json:"schools"`
	Pagination PaginationMeta   `json:"pagination"`
}

// PaginationMeta represents pagination metadata
type PaginationMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// SchoolFilter represents filter options for listing schools
type SchoolFilter struct {
	Name     string `query:"name"`
	IsActive *bool  `query:"is_active"`
	Page     int    `query:"page"`
	PageSize int    `query:"page_size"`
}

// DefaultSchoolFilter returns default filter values
func DefaultSchoolFilter() SchoolFilter {
	return SchoolFilter{
		Page:     1,
		PageSize: 20,
	}
}

// ActivateDeactivateResponse represents the response for activate/deactivate operations
type ActivateDeactivateResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
	Message  string `json:"message"`
}
