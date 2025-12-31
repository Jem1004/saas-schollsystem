package tenant

import (
	"context"
	"errors"
	"strings"

	"github.com/school-management/backend/internal/domain/models"
)

var (
	ErrNameRequired     = errors.New("school name is required")
	ErrSchoolInactive   = errors.New("school is already inactive")
	ErrSchoolActive     = errors.New("school is already active")
)

// Service defines the interface for tenant business logic
type Service interface {
	CreateSchool(ctx context.Context, req CreateSchoolRequest) (*SchoolResponse, error)
	GetAllSchools(ctx context.Context, filter SchoolFilter) (*SchoolListResponse, error)
	GetSchoolByID(ctx context.Context, id uint) (*SchoolResponse, error)
	UpdateSchool(ctx context.Context, id uint, req UpdateSchoolRequest) (*SchoolResponse, error)
	DeactivateSchool(ctx context.Context, id uint) (*ActivateDeactivateResponse, error)
	ActivateSchool(ctx context.Context, id uint) (*ActivateDeactivateResponse, error)
}

// service implements the Service interface
type service struct {
	repo Repository
}

// NewService creates a new tenant service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// CreateSchool creates a new school (tenant)
// Requirements: 1.1 - WHEN a Super_Admin creates a new tenant, THE System SHALL generate a unique school_id and isolate all data
func (s *service) CreateSchool(ctx context.Context, req CreateSchoolRequest) (*SchoolResponse, error) {
	// Validate required fields
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, ErrNameRequired
	}

	// Check for duplicate name (optional - can be removed if duplicates are allowed)
	existing, err := s.repo.FindByName(ctx, name)
	if err != nil && !errors.Is(err, ErrSchoolNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, ErrDuplicateSchool
	}

	// Create school
	school := &models.School{
		Name:     name,
		Address:  strings.TrimSpace(req.Address),
		Phone:    strings.TrimSpace(req.Phone),
		Email:    strings.TrimSpace(req.Email),
		IsActive: true,
	}

	if err := s.repo.Create(ctx, school); err != nil {
		return nil, err
	}

	return toSchoolResponse(school, nil), nil
}

// GetAllSchools retrieves all schools with pagination
// Requirements: 1.2 - WHEN a Super_Admin views the tenant list, THE System SHALL display all registered schools with their status
func (s *service) GetAllSchools(ctx context.Context, filter SchoolFilter) (*SchoolListResponse, error) {
	// Set defaults
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}

	schools, total, err := s.repo.FindAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Convert to response
	schoolResponses := make([]SchoolResponse, len(schools))
	for i, school := range schools {
		// Get stats for each school
		stats, _ := s.repo.GetStats(ctx, school.ID)
		schoolResponses[i] = *toSchoolResponse(&school, stats)
	}

	// Calculate total pages
	totalPages := int(total) / filter.PageSize
	if int(total)%filter.PageSize > 0 {
		totalPages++
	}

	return &SchoolListResponse{
		Schools: schoolResponses,
		Pagination: PaginationMeta{
			Page:       filter.Page,
			PageSize:   filter.PageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// GetSchoolByID retrieves a school by ID
func (s *service) GetSchoolByID(ctx context.Context, id uint) (*SchoolResponse, error) {
	school, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Get stats
	stats, _ := s.repo.GetStats(ctx, school.ID)

	return toSchoolResponse(school, stats), nil
}

// UpdateSchool updates a school
func (s *service) UpdateSchool(ctx context.Context, id uint, req UpdateSchoolRequest) (*SchoolResponse, error) {
	// Get existing school
	school, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return nil, ErrNameRequired
		}
		// Check for duplicate name if name is being changed
		if name != school.Name {
			existing, err := s.repo.FindByName(ctx, name)
			if err != nil && !errors.Is(err, ErrSchoolNotFound) {
				return nil, err
			}
			if existing != nil && existing.ID != school.ID {
				return nil, ErrDuplicateSchool
			}
		}
		school.Name = name
	}
	if req.Address != nil {
		school.Address = strings.TrimSpace(*req.Address)
	}
	if req.Phone != nil {
		school.Phone = strings.TrimSpace(*req.Phone)
	}
	if req.Email != nil {
		school.Email = strings.TrimSpace(*req.Email)
	}

	if err := s.repo.Update(ctx, school); err != nil {
		return nil, err
	}

	// Get stats
	stats, _ := s.repo.GetStats(ctx, school.ID)

	return toSchoolResponse(school, stats), nil
}

// DeactivateSchool deactivates a school
// Requirements: 1.3 - WHEN a Super_Admin deactivates a tenant, THE System SHALL prevent all users under that tenant from accessing the system
func (s *service) DeactivateSchool(ctx context.Context, id uint) (*ActivateDeactivateResponse, error) {
	// Get existing school
	school, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check if already inactive
	if !school.IsActive {
		return nil, ErrSchoolInactive
	}

	// Deactivate
	if err := s.repo.Deactivate(ctx, id); err != nil {
		return nil, err
	}

	return &ActivateDeactivateResponse{
		ID:       school.ID,
		Name:     school.Name,
		IsActive: false,
		Message:  "School deactivated successfully. All users under this school will no longer be able to access the system.",
	}, nil
}

// ActivateSchool activates a school
func (s *service) ActivateSchool(ctx context.Context, id uint) (*ActivateDeactivateResponse, error) {
	// Get existing school
	school, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check if already active
	if school.IsActive {
		return nil, ErrSchoolActive
	}

	// Activate
	if err := s.repo.Activate(ctx, id); err != nil {
		return nil, err
	}

	return &ActivateDeactivateResponse{
		ID:       school.ID,
		Name:     school.Name,
		IsActive: true,
		Message:  "School activated successfully. Users under this school can now access the system.",
	}, nil
}

// toSchoolResponse converts a School model to SchoolResponse DTO
func toSchoolResponse(school *models.School, stats *SchoolStats) *SchoolResponse {
	return &SchoolResponse{
		ID:        school.ID,
		Name:      school.Name,
		Address:   school.Address,
		Phone:     school.Phone,
		Email:     school.Email,
		IsActive:  school.IsActive,
		CreatedAt: school.CreatedAt,
		UpdatedAt: school.UpdatedAt,
		Stats:     stats,
	}
}
