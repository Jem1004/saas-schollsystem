package bk

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/school-management/backend/internal/domain/models"
)

var (
	ErrStudentIDRequired          = errors.New("ID siswa wajib diisi")
	ErrCategoryRequired           = errors.New("kategori wajib diisi")
	ErrLevelRequired              = errors.New("tingkat wajib diisi")
	ErrDescriptionRequired        = errors.New("deskripsi wajib diisi")
	ErrTitleRequired              = errors.New("judul wajib diisi")
	ErrPointRequired              = errors.New("poin harus lebih dari 0")
	ErrReasonRequired             = errors.New("alasan wajib diisi")
	ErrExitTimeRequired           = errors.New("waktu keluar wajib diisi")
	ErrResponsibleTeacherRequired = errors.New("guru penanggung jawab wajib diisi")
	ErrInternalNoteRequired       = errors.New("catatan internal wajib diisi")
	ErrReturnTimeRequired         = errors.New("waktu kembali wajib diisi")
	ErrAlreadyReturned            = errors.New("siswa sudah kembali")
	ErrStudentNotInSchool         = errors.New("siswa bukan dari sekolah ini")
	ErrTeacherNotInSchool         = errors.New("guru bukan dari sekolah ini")
	ErrInvalidViolationLevel      = errors.New("tingkat pelanggaran tidak valid")
)

// Service defines the interface for BK business logic
type Service interface {
	// Violation operations
	CreateViolation(ctx context.Context, schoolID, createdBy uint, req CreateViolationRequest) (*ViolationResponse, error)
	GetViolationByID(ctx context.Context, id uint) (*ViolationResponse, error)
	GetStudentViolations(ctx context.Context, studentID uint) ([]ViolationResponse, error)
	GetViolations(ctx context.Context, schoolID uint, filter ViolationFilter) (*ViolationListResponse, error)
	DeleteViolation(ctx context.Context, id uint) error
	GetStudentViolationPoints(ctx context.Context, studentID uint) (int, error)

	// Violation Category operations
	CreateViolationCategory(ctx context.Context, schoolID uint, req CreateViolationCategoryRequest) (*ViolationCategoryResponse, error)
	GetViolationCategories(ctx context.Context, schoolID uint, activeOnly bool) (*ViolationCategoryListResponse, error)
	GetViolationCategoryByID(ctx context.Context, id uint) (*ViolationCategoryResponse, error)
	UpdateViolationCategory(ctx context.Context, id uint, req UpdateViolationCategoryRequest) (*ViolationCategoryResponse, error)
	DeleteViolationCategory(ctx context.Context, id uint) error
	InitializeDefaultCategories(ctx context.Context, schoolID uint) error

	// Achievement operations
	CreateAchievement(ctx context.Context, schoolID, createdBy uint, req CreateAchievementRequest) (*AchievementResponse, error)
	GetAchievementByID(ctx context.Context, id uint) (*AchievementResponse, error)
	GetStudentAchievements(ctx context.Context, studentID uint) ([]AchievementResponse, error)
	GetAchievements(ctx context.Context, schoolID uint, filter AchievementFilter) (*AchievementListResponse, error)
	GetStudentAchievementPoints(ctx context.Context, studentID uint) (*AchievementPointsResponse, error)
	DeleteAchievement(ctx context.Context, id uint) error

	// Permit operations
	CreatePermit(ctx context.Context, schoolID, createdBy uint, req CreatePermitRequest) (*PermitResponse, error)
	GetPermitByID(ctx context.Context, id uint) (*PermitResponse, error)
	GetStudentPermits(ctx context.Context, studentID uint) ([]PermitResponse, error)
	GetPermits(ctx context.Context, schoolID uint, filter PermitFilter) (*PermitListResponse, error)
	RecordReturn(ctx context.Context, permitID uint, req RecordReturnRequest) (*PermitResponse, error)
	GetPermitDocument(ctx context.Context, permitID uint) (*PermitDocumentData, error)
	DeletePermit(ctx context.Context, id uint) error

	// Counseling Note operations
	CreateCounselingNote(ctx context.Context, schoolID, createdBy uint, req CreateCounselingNoteRequest) (*CounselingNoteFullResponse, error)
	GetCounselingNoteByID(ctx context.Context, id uint, includeInternal bool) (interface{}, error)
	GetStudentCounselingNotes(ctx context.Context, studentID uint, includeInternal bool) (interface{}, error)
	GetCounselingNotes(ctx context.Context, schoolID uint, filter CounselingNoteFilter, includeInternal bool) (interface{}, error)
	UpdateCounselingNote(ctx context.Context, id uint, req UpdateCounselingNoteRequest) (*CounselingNoteFullResponse, error)
	DeleteCounselingNote(ctx context.Context, id uint) error

	// Student BK Profile
	GetStudentBKProfile(ctx context.Context, studentID uint, includeInternal bool) (interface{}, error)

	// Dashboard
	GetBKDashboard(ctx context.Context, schoolID uint) (*BKDashboardResponse, error)
}

// service implements the Service interface
type service struct {
	repo Repository
}

// NewService creates a new BK service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// ==================== Violation Service ====================

// CreateViolation creates a new violation record
// Requirements: 6.1 - WHEN a Guru_BK records a violation, THE System SHALL require category, level, description, and student identifier
func (s *service) CreateViolation(ctx context.Context, schoolID, createdBy uint, req CreateViolationRequest) (*ViolationResponse, error) {
	// Validate required fields
	if req.StudentID == 0 {
		return nil, ErrStudentIDRequired
	}
	if req.Category == "" {
		return nil, ErrCategoryRequired
	}
	if !req.Level.IsValid() {
		return nil, ErrInvalidViolationLevel
	}
	if req.Description == "" {
		return nil, ErrDescriptionRequired
	}

	// Verify student belongs to the school
	student, err := s.repo.FindStudentByID(ctx, req.StudentID)
	if err != nil {
		return nil, err
	}
	if student.SchoolID != schoolID {
		return nil, ErrStudentNotInSchool
	}

	// Determine point value
	point := -5 // default
	if req.Point != nil {
		point = *req.Point
	} else if req.CategoryID != nil {
		// Get default point from category
		category, err := s.repo.FindViolationCategoryByID(ctx, *req.CategoryID)
		if err == nil && category != nil {
			point = category.DefaultPoint
		}
	} else {
		// Set default point based on level
		switch req.Level {
		case models.ViolationLevelRingan:
			point = -5
		case models.ViolationLevelSedang:
			point = -15
		case models.ViolationLevelBerat:
			point = -30
		}
	}

	violation := &models.Violation{
		StudentID:   req.StudentID,
		CategoryID:  req.CategoryID,
		Category:    req.Category,
		Level:       req.Level,
		Point:       point,
		Description: req.Description,
		CreatedBy:   createdBy,
	}

	if err := s.repo.CreateViolation(ctx, violation); err != nil {
		return nil, err
	}

	// Reload with relations
	violation, err = s.repo.FindViolationByID(ctx, violation.ID)
	if err != nil {
		return nil, err
	}

	// TODO: Trigger notification to parent (async)
	// Requirements: 6.2 - WHEN a violation is saved, THE System SHALL trigger notification to the parent

	return toViolationResponse(violation), nil
}

// GetViolationByID retrieves a violation by ID
func (s *service) GetViolationByID(ctx context.Context, id uint) (*ViolationResponse, error) {
	violation, err := s.repo.FindViolationByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toViolationResponse(violation), nil
}

// GetStudentViolations retrieves all violations for a student
// Requirements: 6.3 - THE System SHALL display all violations for a student sorted by date
func (s *service) GetStudentViolations(ctx context.Context, studentID uint) ([]ViolationResponse, error) {
	violations, err := s.repo.FindViolationsByStudent(ctx, studentID)
	if err != nil {
		return nil, err
	}

	responses := make([]ViolationResponse, len(violations))
	for i, v := range violations {
		responses[i] = *toViolationResponse(&v)
	}
	return responses, nil
}

// GetViolations retrieves violations with pagination and filtering
func (s *service) GetViolations(ctx context.Context, schoolID uint, filter ViolationFilter) (*ViolationListResponse, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}

	violations, total, err := s.repo.FindViolations(ctx, schoolID, filter)
	if err != nil {
		return nil, err
	}

	responses := make([]ViolationResponse, len(violations))
	for i, v := range violations {
		responses[i] = *toViolationResponse(&v)
	}

	totalPages := int(total) / filter.PageSize
	if int(total)%filter.PageSize > 0 {
		totalPages++
	}

	return &ViolationListResponse{
		Violations: responses,
		Pagination: PaginationMeta{
			Page:       filter.Page,
			PageSize:   filter.PageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// DeleteViolation deletes a violation record
func (s *service) DeleteViolation(ctx context.Context, id uint) error {
	return s.repo.DeleteViolation(ctx, id)
}

// GetStudentViolationPoints retrieves total violation points for a student
func (s *service) GetStudentViolationPoints(ctx context.Context, studentID uint) (int, error) {
	return s.repo.GetStudentViolationPoints(ctx, studentID)
}

// ==================== Violation Category Service ====================

// CreateViolationCategory creates a new violation category
func (s *service) CreateViolationCategory(ctx context.Context, schoolID uint, req CreateViolationCategoryRequest) (*ViolationCategoryResponse, error) {
	if req.Name == "" {
		return nil, errors.New("nama kategori wajib diisi")
	}
	if req.DefaultPoint > 0 {
		return nil, errors.New("poin default harus 0 atau negatif")
	}
	if !req.DefaultLevel.IsValid() {
		return nil, ErrInvalidViolationLevel
	}

	category := &models.ViolationCategory{
		SchoolID:     schoolID,
		Name:         req.Name,
		DefaultPoint: req.DefaultPoint,
		DefaultLevel: req.DefaultLevel,
		Description:  req.Description,
		IsActive:     true,
	}

	if err := s.repo.CreateViolationCategory(ctx, category); err != nil {
		return nil, err
	}

	return toViolationCategoryResponse(category), nil
}

// GetViolationCategories retrieves all violation categories for a school
func (s *service) GetViolationCategories(ctx context.Context, schoolID uint, activeOnly bool) (*ViolationCategoryListResponse, error) {
	categories, err := s.repo.FindViolationCategories(ctx, schoolID, activeOnly)
	if err != nil {
		return nil, err
	}

	responses := make([]ViolationCategoryResponse, len(categories))
	for i, c := range categories {
		responses[i] = *toViolationCategoryResponse(&c)
	}

	return &ViolationCategoryListResponse{Categories: responses}, nil
}

// GetViolationCategoryByID retrieves a violation category by ID
func (s *service) GetViolationCategoryByID(ctx context.Context, id uint) (*ViolationCategoryResponse, error) {
	category, err := s.repo.FindViolationCategoryByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toViolationCategoryResponse(category), nil
}

// UpdateViolationCategory updates a violation category
func (s *service) UpdateViolationCategory(ctx context.Context, id uint, req UpdateViolationCategoryRequest) (*ViolationCategoryResponse, error) {
	category, err := s.repo.FindViolationCategoryByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	if req.DefaultPoint != nil {
		if *req.DefaultPoint > 0 {
			return nil, errors.New("poin default harus 0 atau negatif")
		}
		category.DefaultPoint = *req.DefaultPoint
	}
	if req.DefaultLevel != "" && req.DefaultLevel.IsValid() {
		category.DefaultLevel = req.DefaultLevel
	}
	if req.Description != "" {
		category.Description = req.Description
	}
	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	}

	if err := s.repo.UpdateViolationCategory(ctx, category); err != nil {
		return nil, err
	}

	return toViolationCategoryResponse(category), nil
}

// DeleteViolationCategory deletes a violation category
func (s *service) DeleteViolationCategory(ctx context.Context, id uint) error {
	return s.repo.DeleteViolationCategory(ctx, id)
}

// InitializeDefaultCategories creates default violation categories for a school
func (s *service) InitializeDefaultCategories(ctx context.Context, schoolID uint) error {
	// Check if categories already exist
	existing, err := s.repo.FindViolationCategories(ctx, schoolID, false)
	if err != nil {
		return err
	}
	if len(existing) > 0 {
		return nil // Already initialized
	}

	defaults := models.DefaultViolationCategories(schoolID)
	for _, cat := range defaults {
		if err := s.repo.CreateViolationCategory(ctx, &cat); err != nil {
			return err
		}
	}
	return nil
}


// ==================== Achievement Service ====================

// CreateAchievement creates a new achievement record
// Requirements: 7.1 - WHEN a Guru_BK records an achievement, THE System SHALL require title, point value, and description
func (s *service) CreateAchievement(ctx context.Context, schoolID, createdBy uint, req CreateAchievementRequest) (*AchievementResponse, error) {
	// Validate required fields
	if req.StudentID == 0 {
		return nil, ErrStudentIDRequired
	}
	if req.Title == "" {
		return nil, ErrTitleRequired
	}
	if req.Point <= 0 {
		return nil, ErrPointRequired
	}

	// Verify student belongs to the school
	student, err := s.repo.FindStudentByID(ctx, req.StudentID)
	if err != nil {
		return nil, err
	}
	if student.SchoolID != schoolID {
		return nil, ErrStudentNotInSchool
	}

	achievement := &models.Achievement{
		StudentID:   req.StudentID,
		Title:       req.Title,
		Point:       req.Point,
		Description: req.Description,
		CreatedBy:   createdBy,
	}

	if err := s.repo.CreateAchievement(ctx, achievement); err != nil {
		return nil, err
	}

	// Reload with relations
	achievement, err = s.repo.FindAchievementByID(ctx, achievement.ID)
	if err != nil {
		return nil, err
	}

	// TODO: Trigger notification to parent (async)
	// Requirements: 7.4 - WHEN an achievement is recorded, THE System SHALL trigger notification to the parent

	return toAchievementResponse(achievement), nil
}

// GetAchievementByID retrieves an achievement by ID
func (s *service) GetAchievementByID(ctx context.Context, id uint) (*AchievementResponse, error) {
	achievement, err := s.repo.FindAchievementByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toAchievementResponse(achievement), nil
}

// GetStudentAchievements retrieves all achievements for a student
// Requirements: 7.5 - THE System SHALL maintain achievement history per student
func (s *service) GetStudentAchievements(ctx context.Context, studentID uint) ([]AchievementResponse, error) {
	achievements, err := s.repo.FindAchievementsByStudent(ctx, studentID)
	if err != nil {
		return nil, err
	}

	responses := make([]AchievementResponse, len(achievements))
	for i, a := range achievements {
		responses[i] = *toAchievementResponse(&a)
	}
	return responses, nil
}

// GetAchievements retrieves achievements with pagination and filtering
func (s *service) GetAchievements(ctx context.Context, schoolID uint, filter AchievementFilter) (*AchievementListResponse, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}

	achievements, total, err := s.repo.FindAchievements(ctx, schoolID, filter)
	if err != nil {
		return nil, err
	}

	responses := make([]AchievementResponse, len(achievements))
	for i, a := range achievements {
		responses[i] = *toAchievementResponse(&a)
	}

	totalPages := int(total) / filter.PageSize
	if int(total)%filter.PageSize > 0 {
		totalPages++
	}

	return &AchievementListResponse{
		Achievements: responses,
		Pagination: PaginationMeta{
			Page:       filter.Page,
			PageSize:   filter.PageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// GetStudentAchievementPoints retrieves total achievement points for a student
// Requirements: 7.2, 7.3 - THE System SHALL display total achievement points
func (s *service) GetStudentAchievementPoints(ctx context.Context, studentID uint) (*AchievementPointsResponse, error) {
	student, err := s.repo.FindStudentByID(ctx, studentID)
	if err != nil {
		return nil, err
	}

	points, err := s.repo.GetStudentAchievementPoints(ctx, studentID)
	if err != nil {
		return nil, err
	}

	return &AchievementPointsResponse{
		StudentID:   studentID,
		StudentName: student.Name,
		TotalPoints: points,
	}, nil
}

// DeleteAchievement deletes an achievement record
func (s *service) DeleteAchievement(ctx context.Context, id uint) error {
	return s.repo.DeleteAchievement(ctx, id)
}

// ==================== Permit Service ====================

// CreatePermit creates a new exit permit
// Requirements: 8.1 - WHEN a Guru_BK creates an exit permit, THE System SHALL require reason, exit time, and responsible teacher
func (s *service) CreatePermit(ctx context.Context, schoolID, createdBy uint, req CreatePermitRequest) (*PermitResponse, error) {
	// Validate required fields
	if req.StudentID == 0 {
		return nil, ErrStudentIDRequired
	}
	if req.Reason == "" {
		return nil, ErrReasonRequired
	}
	if req.ExitTime.IsZero() {
		return nil, ErrExitTimeRequired
	}
	if req.ResponsibleTeacher == 0 {
		return nil, ErrResponsibleTeacherRequired
	}

	// Verify student belongs to the school
	student, err := s.repo.FindStudentByID(ctx, req.StudentID)
	if err != nil {
		return nil, err
	}
	if student.SchoolID != schoolID {
		return nil, ErrStudentNotInSchool
	}

	// Verify teacher belongs to the school
	teacher, err := s.repo.FindUserByID(ctx, req.ResponsibleTeacher)
	if err != nil {
		return nil, err
	}
	if teacher.SchoolID != nil && *teacher.SchoolID != schoolID {
		return nil, ErrTeacherNotInSchool
	}

	permit := &models.Permit{
		StudentID:          req.StudentID,
		Reason:             req.Reason,
		ExitTime:           req.ExitTime,
		ResponsibleTeacher: req.ResponsibleTeacher,
		CreatedBy:          createdBy,
	}

	if err := s.repo.CreatePermit(ctx, permit); err != nil {
		return nil, err
	}

	// Reload with relations
	permit, err = s.repo.FindPermitByID(ctx, permit.ID)
	if err != nil {
		return nil, err
	}

	// TODO: Trigger notification to parent (async)
	// Requirements: 8.3 - WHEN a permit is issued, THE System SHALL trigger notification to the parent

	return toPermitResponse(permit), nil
}

// GetPermitByID retrieves a permit by ID
func (s *service) GetPermitByID(ctx context.Context, id uint) (*PermitResponse, error) {
	permit, err := s.repo.FindPermitByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toPermitResponse(permit), nil
}

// GetStudentPermits retrieves all permits for a student
func (s *service) GetStudentPermits(ctx context.Context, studentID uint) ([]PermitResponse, error) {
	permits, err := s.repo.FindPermitsByStudent(ctx, studentID)
	if err != nil {
		return nil, err
	}

	responses := make([]PermitResponse, len(permits))
	for i, p := range permits {
		responses[i] = *toPermitResponse(&p)
	}
	return responses, nil
}

// GetPermits retrieves permits with pagination and filtering
func (s *service) GetPermits(ctx context.Context, schoolID uint, filter PermitFilter) (*PermitListResponse, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}

	permits, total, err := s.repo.FindPermits(ctx, schoolID, filter)
	if err != nil {
		return nil, err
	}

	responses := make([]PermitResponse, len(permits))
	for i, p := range permits {
		responses[i] = *toPermitResponse(&p)
	}

	totalPages := int(total) / filter.PageSize
	if int(total)%filter.PageSize > 0 {
		totalPages++
	}

	return &PermitListResponse{
		Permits: responses,
		Pagination: PaginationMeta{
			Page:       filter.Page,
			PageSize:   filter.PageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// RecordReturn records the return time for a permit
// Requirements: 8.4 - WHEN a student returns, THE System SHALL allow recording of return time
func (s *service) RecordReturn(ctx context.Context, permitID uint, req RecordReturnRequest) (*PermitResponse, error) {
	if req.ReturnTime.IsZero() {
		return nil, ErrReturnTimeRequired
	}

	permit, err := s.repo.FindPermitByID(ctx, permitID)
	if err != nil {
		return nil, err
	}

	if permit.HasReturned() {
		return nil, ErrAlreadyReturned
	}

	permit.ReturnTime = &req.ReturnTime

	if err := s.repo.UpdatePermit(ctx, permit); err != nil {
		return nil, err
	}

	return toPermitResponse(permit), nil
}

// GetPermitDocument generates permit document data
// Requirements: 8.2, 8.5 - THE System SHALL generate a PDF/receipt document with student info, reason, and timestamp
func (s *service) GetPermitDocument(ctx context.Context, permitID uint) (*PermitDocumentData, error) {
	permit, err := s.repo.FindPermitByID(ctx, permitID)
	if err != nil {
		return nil, err
	}

	// Handle nullable Class pointer
	className := ""
	if permit.Student.Class != nil {
		className = permit.Student.Class.Name
	}

	return &PermitDocumentData{
		PermitID:           permit.ID,
		StudentName:        permit.Student.Name,
		StudentNIS:         permit.Student.NIS,
		StudentNISN:        permit.Student.NISN,
		ClassName:          className,
		SchoolName:         permit.Student.School.Name,
		Reason:             permit.Reason,
		ExitTime:           permit.ExitTime,
		ResponsibleTeacher: permit.Teacher.Username,
		GeneratedAt:        time.Now(),
	}, nil
}

// DeletePermit deletes a permit record
func (s *service) DeletePermit(ctx context.Context, id uint) error {
	return s.repo.DeletePermit(ctx, id)
}


// ==================== Counseling Note Service ====================

// CreateCounselingNote creates a new counseling note
// Requirements: 9.1 - WHEN a Guru_BK creates a counseling note, THE System SHALL require internal_note and optional parent_summary
func (s *service) CreateCounselingNote(ctx context.Context, schoolID, createdBy uint, req CreateCounselingNoteRequest) (*CounselingNoteFullResponse, error) {
	// Validate required fields
	if req.StudentID == 0 {
		return nil, ErrStudentIDRequired
	}
	if req.InternalNote == "" {
		return nil, ErrInternalNoteRequired
	}

	// Verify student belongs to the school
	student, err := s.repo.FindStudentByID(ctx, req.StudentID)
	if err != nil {
		return nil, err
	}
	if student.SchoolID != schoolID {
		return nil, ErrStudentNotInSchool
	}

	note := &models.CounselingNote{
		StudentID:     req.StudentID,
		InternalNote:  req.InternalNote,
		ParentSummary: req.ParentSummary,
		CreatedBy:     createdBy,
	}

	if err := s.repo.CreateCounselingNote(ctx, note); err != nil {
		return nil, err
	}

	// Reload with relations
	note, err = s.repo.FindCounselingNoteByID(ctx, note.ID)
	if err != nil {
		return nil, err
	}

	// TODO: Trigger notification to parent if parent_summary is provided (async)
	// Requirements: 9.2 - WHEN parent_summary is provided, THE System SHALL make it visible to parents

	return toCounselingNoteFullResponse(note), nil
}

// GetCounselingNoteByID retrieves a counseling note by ID
// Requirements: 9.3, 9.4 - Internal note accessible only to Guru BK
func (s *service) GetCounselingNoteByID(ctx context.Context, id uint, includeInternal bool) (interface{}, error) {
	note, err := s.repo.FindCounselingNoteByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if includeInternal {
		return toCounselingNoteFullResponse(note), nil
	}
	return toCounselingNoteResponse(note), nil
}

// GetStudentCounselingNotes retrieves all counseling notes for a student
// Requirements: 9.5 - THE System SHALL maintain counseling history per student with timestamps
func (s *service) GetStudentCounselingNotes(ctx context.Context, studentID uint, includeInternal bool) (interface{}, error) {
	notes, err := s.repo.FindCounselingNotesByStudent(ctx, studentID)
	if err != nil {
		return nil, err
	}

	if includeInternal {
		responses := make([]CounselingNoteFullResponse, len(notes))
		for i, n := range notes {
			responses[i] = *toCounselingNoteFullResponse(&n)
		}
		return responses, nil
	}

	responses := make([]CounselingNoteResponse, len(notes))
	for i, n := range notes {
		responses[i] = *toCounselingNoteResponse(&n)
	}
	return responses, nil
}

// GetCounselingNotes retrieves counseling notes with pagination and filtering
func (s *service) GetCounselingNotes(ctx context.Context, schoolID uint, filter CounselingNoteFilter, includeInternal bool) (interface{}, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}

	notes, total, err := s.repo.FindCounselingNotes(ctx, schoolID, filter)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / filter.PageSize
	if int(total)%filter.PageSize > 0 {
		totalPages++
	}

	pagination := PaginationMeta{
		Page:       filter.Page,
		PageSize:   filter.PageSize,
		Total:      total,
		TotalPages: totalPages,
	}

	if includeInternal {
		responses := make([]CounselingNoteFullResponse, len(notes))
		for i, n := range notes {
			responses[i] = *toCounselingNoteFullResponse(&n)
		}
		return &CounselingNoteFullListResponse{
			Notes:      responses,
			Pagination: pagination,
		}, nil
	}

	responses := make([]CounselingNoteResponse, len(notes))
	for i, n := range notes {
		responses[i] = *toCounselingNoteResponse(&n)
	}
	return &CounselingNoteListResponse{
		Notes:      responses,
		Pagination: pagination,
	}, nil
}

// UpdateCounselingNote updates a counseling note
func (s *service) UpdateCounselingNote(ctx context.Context, id uint, req UpdateCounselingNoteRequest) (*CounselingNoteFullResponse, error) {
	note, err := s.repo.FindCounselingNoteByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.InternalNote != "" {
		note.InternalNote = req.InternalNote
	}
	if req.ParentSummary != "" {
		note.ParentSummary = req.ParentSummary
	}

	if err := s.repo.UpdateCounselingNote(ctx, note); err != nil {
		return nil, err
	}

	return toCounselingNoteFullResponse(note), nil
}

// DeleteCounselingNote deletes a counseling note
func (s *service) DeleteCounselingNote(ctx context.Context, id uint) error {
	return s.repo.DeleteCounselingNote(ctx, id)
}

// ==================== Student BK Profile ====================

// GetStudentBKProfile retrieves a student's complete BK profile
// Requirements: 6.3, 7.5, 8.4, 9.5 - Timeline view per student
func (s *service) GetStudentBKProfile(ctx context.Context, studentID uint, includeInternal bool) (interface{}, error) {
	student, err := s.repo.FindStudentByID(ctx, studentID)
	if err != nil {
		return nil, err
	}

	// Get counts
	violations, _ := s.repo.FindViolationsByStudent(ctx, studentID)
	achievements, _ := s.repo.FindAchievementsByStudent(ctx, studentID)
	permits, _ := s.repo.FindPermitsByStudent(ctx, studentID)
	counselingNotes, _ := s.repo.FindCounselingNotesByStudent(ctx, studentID)

	// Get total points
	totalPoints, _ := s.repo.GetStudentAchievementPoints(ctx, studentID)

	// Limit to recent items (5 each)
	recentViolations := make([]ViolationResponse, 0)
	for i, v := range violations {
		if i >= 5 {
			break
		}
		recentViolations = append(recentViolations, *toViolationResponse(&v))
	}

	recentAchievements := make([]AchievementResponse, 0)
	for i, a := range achievements {
		if i >= 5 {
			break
		}
		recentAchievements = append(recentAchievements, *toAchievementResponse(&a))
	}

	recentPermits := make([]PermitResponse, 0)
	for i, p := range permits {
		if i >= 5 {
			break
		}
		recentPermits = append(recentPermits, *toPermitResponse(&p))
	}

	className := ""
	if student.Class.ID != 0 {
		className = student.Class.Name
	}

	if includeInternal {
		recentCounseling := make([]CounselingNoteFullResponse, 0)
		for i, n := range counselingNotes {
			if i >= 5 {
				break
			}
			recentCounseling = append(recentCounseling, *toCounselingNoteFullResponse(&n))
		}

		return &StudentBKProfileFullResponse{
			StudentID:          studentID,
			StudentName:        student.Name,
			StudentNIS:         student.NIS,
			StudentNISN:        student.NISN,
			ClassName:          className,
			TotalPoints:        totalPoints,
			ViolationCount:     len(violations),
			AchievementCount:   len(achievements),
			PermitCount:        len(permits),
			CounselingCount:    len(counselingNotes),
			RecentViolations:   recentViolations,
			RecentAchievements: recentAchievements,
			RecentPermits:      recentPermits,
			RecentCounseling:   recentCounseling,
		}, nil
	}

	recentCounseling := make([]CounselingNoteResponse, 0)
	for i, n := range counselingNotes {
		if i >= 5 {
			break
		}
		recentCounseling = append(recentCounseling, *toCounselingNoteResponse(&n))
	}

	return &StudentBKProfileResponse{
		StudentID:          studentID,
		StudentName:        student.Name,
		StudentNIS:         student.NIS,
		StudentNISN:        student.NISN,
		ClassName:          className,
		TotalPoints:        totalPoints,
		ViolationCount:     len(violations),
		AchievementCount:   len(achievements),
		PermitCount:        len(permits),
		CounselingCount:    len(counselingNotes),
		RecentViolations:   recentViolations,
		RecentAchievements: recentAchievements,
		RecentPermits:      recentPermits,
		RecentCounseling:   recentCounseling,
	}, nil
}

// ==================== Dashboard ====================

// GetBKDashboard retrieves BK dashboard data
// Requirements: 6.1, 7.1 - Overview: recent violations, achievements
func (s *service) GetBKDashboard(ctx context.Context, schoolID uint) (*BKDashboardResponse, error) {
	// Get counts
	violationCount, _ := s.repo.GetViolationCount(ctx, schoolID)
	achievementCount, _ := s.repo.GetAchievementCount(ctx, schoolID)
	permitCount, _ := s.repo.GetPermitCount(ctx, schoolID)
	activePermitCount, _ := s.repo.GetActivePermitCount(ctx, schoolID)
	counselingCount, _ := s.repo.GetCounselingCount(ctx, schoolID)

	// Get recent violations (limit 5)
	violationFilter := ViolationFilter{Page: 1, PageSize: 5}
	violations, _, _ := s.repo.FindViolations(ctx, schoolID, violationFilter)
	recentViolations := make([]ViolationResponse, len(violations))
	for i, v := range violations {
		recentViolations[i] = *toViolationResponse(&v)
	}

	// Get recent achievements (limit 5)
	achievementFilter := AchievementFilter{Page: 1, PageSize: 5}
	achievements, _, _ := s.repo.FindAchievements(ctx, schoolID, achievementFilter)
	recentAchievements := make([]AchievementResponse, len(achievements))
	for i, a := range achievements {
		recentAchievements[i] = *toAchievementResponse(&a)
	}

	// Get students needing attention
	studentsNeedingAttention, _ := s.repo.GetStudentsNeedingAttention(ctx, schoolID, 10)

	return &BKDashboardResponse{
		TotalViolations:          int(violationCount),
		TotalAchievements:        int(achievementCount),
		TotalPermits:             int(permitCount),
		ActivePermits:            int(activePermitCount),
		TotalCounseling:          int(counselingCount),
		RecentViolations:         recentViolations,
		RecentAchievements:       recentAchievements,
		StudentsNeedingAttention: studentsNeedingAttention,
	}, nil
}

// ==================== Response Converters ====================

func toViolationResponse(v *models.Violation) *ViolationResponse {
	response := &ViolationResponse{
		ID:          v.ID,
		StudentID:   v.StudentID,
		CategoryID:  v.CategoryID,
		Category:    v.Category,
		Level:       v.Level,
		Point:       v.Point,
		Description: v.Description,
		CreatedBy:   v.CreatedBy,
		CreatedAt:   v.CreatedAt,
	}

	if v.Student.ID != 0 {
		response.StudentName = v.Student.Name
		response.StudentNIS = v.Student.NIS
		response.StudentNISN = v.Student.NISN
		// Handle nullable Class pointer
		if v.Student.Class != nil && v.Student.Class.ID != 0 {
			response.ClassName = v.Student.Class.Name
		}
	}

	if v.Creator.ID != 0 {
		response.CreatorName = v.Creator.Username
	}

	return response
}

func toViolationCategoryResponse(c *models.ViolationCategory) *ViolationCategoryResponse {
	return &ViolationCategoryResponse{
		ID:           c.ID,
		SchoolID:     c.SchoolID,
		Name:         c.Name,
		DefaultPoint: c.DefaultPoint,
		DefaultLevel: c.DefaultLevel,
		Description:  c.Description,
		IsActive:     c.IsActive,
		CreatedAt:    c.CreatedAt,
		UpdatedAt:    c.UpdatedAt,
	}
}

func toAchievementResponse(a *models.Achievement) *AchievementResponse {
	response := &AchievementResponse{
		ID:          a.ID,
		StudentID:   a.StudentID,
		Title:       a.Title,
		Point:       a.Point,
		Description: a.Description,
		CreatedBy:   a.CreatedBy,
		CreatedAt:   a.CreatedAt,
	}

	if a.Student.ID != 0 {
		response.StudentName = a.Student.Name
		response.StudentNIS = a.Student.NIS
		response.StudentNISN = a.Student.NISN
		// Handle nullable Class pointer
		if a.Student.Class != nil && a.Student.Class.ID != 0 {
			response.ClassName = a.Student.Class.Name
		}
	}

	if a.Creator.ID != 0 {
		response.CreatorName = a.Creator.Username
	}

	return response
}

func toPermitResponse(p *models.Permit) *PermitResponse {
	response := &PermitResponse{
		ID:                 p.ID,
		StudentID:          p.StudentID,
		Reason:             p.Reason,
		ExitTime:           p.ExitTime,
		ReturnTime:         p.ReturnTime,
		ResponsibleTeacher: p.ResponsibleTeacher,
		DocumentURL:        p.DocumentURL,
		CreatedBy:          p.CreatedBy,
		CreatedAt:          p.CreatedAt,
		HasReturned:        p.HasReturned(),
	}

	if p.Student.ID != 0 {
		response.StudentName = p.Student.Name
		response.StudentNIS = p.Student.NIS
		response.StudentNISN = p.Student.NISN
		// Handle nullable Class pointer
		if p.Student.Class != nil && p.Student.Class.ID != 0 {
			response.ClassName = p.Student.Class.Name
		}
	}

	if p.Teacher.ID != 0 {
		response.TeacherName = p.Teacher.Username
	}

	if p.Creator.ID != 0 {
		response.CreatorName = p.Creator.Username
	}

	return response
}

func toCounselingNoteResponse(n *models.CounselingNote) *CounselingNoteResponse {
	response := &CounselingNoteResponse{
		ID:            n.ID,
		StudentID:     n.StudentID,
		ParentSummary: n.ParentSummary,
		CreatedBy:     n.CreatedBy,
		CreatedAt:     n.CreatedAt,
	}

	if n.Student.ID != 0 {
		response.StudentName = n.Student.Name
		response.StudentNIS = n.Student.NIS
		response.StudentNISN = n.Student.NISN
		// Handle nullable Class pointer
		if n.Student.Class != nil && n.Student.Class.ID != 0 {
			response.ClassName = n.Student.Class.Name
		}
	}

	if n.Creator.ID != 0 {
		response.CreatorName = n.Creator.Username
	}

	return response
}

func toCounselingNoteFullResponse(n *models.CounselingNote) *CounselingNoteFullResponse {
	response := &CounselingNoteFullResponse{
		ID:            n.ID,
		StudentID:     n.StudentID,
		InternalNote:  n.InternalNote,
		ParentSummary: n.ParentSummary,
		CreatedBy:     n.CreatedBy,
		CreatedAt:     n.CreatedAt,
	}

	if n.Student.ID != 0 {
		response.StudentName = n.Student.Name
		response.StudentNIS = n.Student.NIS
		response.StudentNISN = n.Student.NISN
		// Handle nullable Class pointer
		if n.Student.Class != nil && n.Student.Class.ID != 0 {
			response.ClassName = n.Student.Class.Name
		}
	}

	if n.Creator.ID != 0 {
		response.CreatorName = n.Creator.Username
	}

	return response
}

// Unused but kept for potential future use
var _ = fmt.Sprintf
