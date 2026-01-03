package tenant

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/school-management/backend/internal/domain/models"
)

var (
	ErrNameRequired       = errors.New("nama sekolah wajib diisi")
	ErrSchoolInactive     = errors.New("sekolah sudah nonaktif")
	ErrSchoolActive       = errors.New("sekolah sudah aktif")
	ErrUsernameExists     = errors.New("Username sudah digunakan")
	ErrInvalidUsername    = errors.New("username hanya boleh berisi huruf, angka, dan underscore")
)

// Service defines the interface for tenant business logic
type Service interface {
	CreateSchool(ctx context.Context, req CreateSchoolRequest) (*SchoolWithAdminResponse, error)
	GetAllSchools(ctx context.Context, filter SchoolFilter) (*SchoolListResponse, error)
	GetSchoolByID(ctx context.Context, id uint) (*SchoolResponse, error)
	GetSchoolDetail(ctx context.Context, id uint) (*SchoolDetailResponse, error)
	UpdateSchool(ctx context.Context, id uint, req UpdateSchoolRequest) (*SchoolResponse, error)
	DeactivateSchool(ctx context.Context, id uint) (*ActivateDeactivateResponse, error)
	ActivateSchool(ctx context.Context, id uint) (*ActivateDeactivateResponse, error)
	DeleteSchool(ctx context.Context, id uint) (*DeleteSchoolResponse, error)
}

// service implements the Service interface
type service struct {
	repo Repository
}

// NewService creates a new tenant service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// generatePassword generates a random password
func generatePassword(length int) string {
	bytes := make([]byte, length)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)[:length]
}

// generateUsername generates a username from school name
func generateUsername(schoolName string) string {
	// Remove special characters and convert to lowercase
	reg := regexp.MustCompile("[^a-zA-Z0-9]+")
	clean := reg.ReplaceAllString(schoolName, "")
	clean = strings.ToLower(clean)
	
	// Limit length
	if len(clean) > 15 {
		clean = clean[:15]
	}
	
	return "admin_" + clean
}

// validateUsername validates username format
func validateUsername(username string) bool {
	reg := regexp.MustCompile("^[a-zA-Z0-9_]+$")
	return reg.MatchString(username)
}

// CreateSchool creates a new school (tenant) with an admin user
// Requirements: 1.1 - WHEN a Super_Admin creates a new tenant, THE System SHALL generate a unique school_id and isolate all data
func (s *service) CreateSchool(ctx context.Context, req CreateSchoolRequest) (*SchoolWithAdminResponse, error) {
	// Validate required fields
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, ErrNameRequired
	}

	// Check for duplicate name
	existing, err := s.repo.FindByName(ctx, name)
	if err != nil && !errors.Is(err, ErrSchoolNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, ErrDuplicateSchool
	}

	// Set default timezone if not provided
	timezone := strings.TrimSpace(req.Timezone)
	if timezone == "" {
		timezone = models.TimezoneWITA // Default to WITA
	}
	// Validate timezone
	if !models.IsValidTimezone(timezone) {
		return nil, errors.New("invalid timezone, must be one of: Asia/Jakarta (WIB), Asia/Makassar (WITA), Asia/Jayapura (WIT)")
	}

	// Create school
	school := &models.School{
		Name:     name,
		Address:  strings.TrimSpace(req.Address),
		Phone:    strings.TrimSpace(req.Phone),
		Email:    strings.TrimSpace(req.Email),
		Timezone: timezone,
		IsActive: true,
	}

	// Prepare admin credentials
	adminUsername := strings.TrimSpace(req.AdminUsername)
	adminPassword := strings.TrimSpace(req.AdminPassword)
	adminName := strings.TrimSpace(req.AdminName)
	adminEmail := strings.TrimSpace(req.AdminEmail)

	// Auto-generate username if not provided
	if adminUsername == "" {
		adminUsername = generateUsername(name)
	}

	// Validate username format
	if !validateUsername(adminUsername) {
		return nil, ErrInvalidUsername
	}

	// Check if username already exists
	exists, err := s.repo.UsernameExists(ctx, adminUsername)
	if err != nil {
		return nil, err
	}
	if exists {
		// Try with a suffix
		for i := 1; i <= 10; i++ {
			newUsername := fmt.Sprintf("%s%d", adminUsername, i)
			exists, err = s.repo.UsernameExists(ctx, newUsername)
			if err != nil {
				return nil, err
			}
			if !exists {
				adminUsername = newUsername
				break
			}
		}
		if exists {
			return nil, ErrUsernameExists
		}
	}

	// Auto-generate password if not provided
	if adminPassword == "" {
		adminPassword = generatePassword(12)
	}

	// Set default admin name if not provided
	if adminName == "" {
		adminName = "Admin " + name
	}

	// Set default admin email if not provided
	if adminEmail == "" && req.Email != "" {
		adminEmail = req.Email
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create admin user
	admin := &models.User{
		Role:         models.RoleAdminSekolah,
		Username:     adminUsername,
		PasswordHash: string(hashedPassword),
		Email:        adminEmail,
		Name:         adminName,
		IsActive:     true,
		MustResetPwd: true, // Force password change on first login
	}

	// Create school with admin in transaction
	if err := s.repo.CreateWithAdmin(ctx, school, admin); err != nil {
		return nil, err
	}

	return &SchoolWithAdminResponse{
		SchoolResponse: *toSchoolResponse(school, nil),
		Admin: &AdminCredentials{
			Username: adminUsername,
			Password: adminPassword, // Return plain password only on creation
			Name:     adminName,
			Email:    adminEmail,
			Message:  "Simpan kredensial ini dengan aman. Password hanya ditampilkan sekali dan harus diganti saat login pertama.",
		},
	}, nil
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

// GetSchoolDetail retrieves a school with admin info
func (s *service) GetSchoolDetail(ctx context.Context, id uint) (*SchoolDetailResponse, error) {
	school, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Get stats
	stats, _ := s.repo.GetStats(ctx, school.ID)

	// Get admin users
	admins, _ := s.repo.GetAdminUsers(ctx, school.ID)

	// Convert admins to AdminInfo
	adminInfos := make([]AdminInfo, len(admins))
	for i, admin := range admins {
		adminInfos[i] = AdminInfo{
			ID:        admin.ID,
			Username:  admin.Username,
			Name:      admin.Name,
			Email:     admin.Email,
			IsActive:  admin.IsActive,
			CreatedAt: admin.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	return &SchoolDetailResponse{
		SchoolResponse: *toSchoolResponse(school, stats),
		Admins:         adminInfos,
	}, nil
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
	if req.Timezone != nil {
		timezone := strings.TrimSpace(*req.Timezone)
		if timezone != "" && !models.IsValidTimezone(timezone) {
			return nil, errors.New("invalid timezone, must be one of: Asia/Jakarta (WIB), Asia/Makassar (WITA), Asia/Jayapura (WIT)")
		}
		if timezone != "" {
			school.Timezone = timezone
		}
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

// DeleteSchool deletes a school and all associated data
func (s *service) DeleteSchool(ctx context.Context, id uint) (*DeleteSchoolResponse, error) {
	// Get existing school
	school, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Get stats before deletion for response
	stats, _ := s.repo.GetStats(ctx, school.ID)

	// Delete school and all related data
	if err := s.repo.Delete(ctx, id); err != nil {
		return nil, err
	}

	response := &DeleteSchoolResponse{
		ID:      school.ID,
		Name:    school.Name,
		Message: "Sekolah dan seluruh data terkait berhasil dihapus secara permanen.",
	}

	if stats != nil {
		response.Deleted.Users = stats.TotalUsers
		response.Deleted.Students = stats.TotalStudents
		response.Deleted.Classes = stats.TotalClasses
		response.Deleted.Devices = stats.TotalDevices
	}

	return response, nil
}

// toSchoolResponse converts a School model to SchoolResponse DTO
func toSchoolResponse(school *models.School, stats *SchoolStats) *SchoolResponse {
	return &SchoolResponse{
		ID:        school.ID,
		Name:      school.Name,
		Address:   school.Address,
		Phone:     school.Phone,
		Email:     school.Email,
		Timezone:  school.Timezone,
		IsActive:  school.IsActive,
		CreatedAt: school.CreatedAt,
		UpdatedAt: school.UpdatedAt,
		Stats:     stats,
	}
}
