package grade

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/school-management/backend/internal/domain/models"
)

var (
	ErrStudentIDRequired   = errors.New("ID siswa wajib diisi")
	ErrTitleRequired       = errors.New("judul wajib diisi")
	ErrScoreInvalid        = errors.New("nilai harus antara 0 dan 100")
	ErrStudentNotInSchool  = errors.New("student does not belong to this school")
	ErrStudentNotInClass   = errors.New("student does not belong to your assigned class")
	ErrNotAuthorized       = errors.New("not authorized to perform this action")
	ErrNoClassAssigned     = errors.New("no class assigned to this teacher")
)

// Service defines the interface for Grade business logic
type Service interface {
	// Grade operations
	CreateGrade(ctx context.Context, schoolID, teacherID uint, req CreateGradeRequest) (*GradeResponse, error)
	GetGradeByID(ctx context.Context, id uint) (*GradeResponse, error)
	GetStudentGrades(ctx context.Context, studentID uint) ([]GradeResponse, error)
	GetGrades(ctx context.Context, schoolID uint, filter GradeFilter) (*GradeListResponse, error)
	UpdateGrade(ctx context.Context, gradeID uint, req UpdateGradeRequest) (*GradeResponse, error)
	DeleteGrade(ctx context.Context, gradeID uint) error

	// Teacher validation
	ValidateTeacherAccess(ctx context.Context, teacherID, studentID uint) error
	GetTeacherClassID(ctx context.Context, teacherID uint) (*uint, error)

	// Summary
	GetStudentGradeSummary(ctx context.Context, studentID uint) (*StudentGradeSummary, error)
	GetClassGrades(ctx context.Context, classID uint, filter GradeFilter) (*GradeListResponse, error)
}

// service implements the Service interface
type service struct {
	repo Repository
	db   *gorm.DB
}

// NewService creates a new Grade service
func NewService(repo Repository, db *gorm.DB) Service {
	return &service{repo: repo, db: db}
}

// CreateGrade creates a new grade record
// Requirements: 10.1 - WHEN a Wali_Kelas inputs a grade, THE System SHALL require title, score
// Requirements: 10.2 - WHEN a grade is saved, THE System SHALL associate it with the student and the Wali_Kelas
// Requirements: 10.5 - THE System SHALL validate that Wali_Kelas can only input grades for students in their assigned class
func (s *service) CreateGrade(ctx context.Context, schoolID, teacherID uint, req CreateGradeRequest) (*GradeResponse, error) {
	// Validate required fields
	if req.StudentID == 0 {
		return nil, ErrStudentIDRequired
	}
	if req.Title == "" {
		return nil, ErrTitleRequired
	}
	if req.Score < 0 || req.Score > 100 {
		return nil, ErrScoreInvalid
	}

	// Verify student belongs to the school
	student, err := s.repo.FindStudentByID(ctx, req.StudentID)
	if err != nil {
		return nil, err
	}
	if student.SchoolID != schoolID {
		return nil, ErrStudentNotInSchool
	}

	// Validate teacher has access to this student (wali kelas validation)
	if err := s.ValidateTeacherAccess(ctx, teacherID, req.StudentID); err != nil {
		return nil, err
	}

	grade := &models.Grade{
		StudentID:   req.StudentID,
		Title:       req.Title,
		Score:       req.Score,
		Description: req.Description,
		CreatedBy:   teacherID,
	}

	if err := s.repo.Create(ctx, grade); err != nil {
		return nil, err
	}

	// Reload with relations
	grade, err = s.repo.FindByID(ctx, grade.ID)
	if err != nil {
		return nil, err
	}

	// TODO: Trigger notification to parent (async)
	// Requirements: 10.3 - WHEN a grade is recorded, THE System SHALL optionally trigger notification to the parent

	return toGradeResponse(grade), nil
}


// GetGradeByID retrieves a grade by ID
func (s *service) GetGradeByID(ctx context.Context, id uint) (*GradeResponse, error) {
	grade, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toGradeResponse(grade), nil
}

// GetStudentGrades retrieves all grades for a student
// Requirements: 10.4 - WHEN parents view grades, THE System SHALL display all grades for their child sorted by date
func (s *service) GetStudentGrades(ctx context.Context, studentID uint) ([]GradeResponse, error) {
	grades, err := s.repo.FindByStudent(ctx, studentID)
	if err != nil {
		return nil, err
	}

	responses := make([]GradeResponse, len(grades))
	for i, g := range grades {
		responses[i] = *toGradeResponse(&g)
	}
	return responses, nil
}

// GetGrades retrieves grades with pagination and filtering
func (s *service) GetGrades(ctx context.Context, schoolID uint, filter GradeFilter) (*GradeListResponse, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}

	grades, total, err := s.repo.FindAll(ctx, schoolID, filter)
	if err != nil {
		return nil, err
	}

	responses := make([]GradeResponse, len(grades))
	for i, g := range grades {
		responses[i] = *toGradeResponse(&g)
	}

	totalPages := int(total) / filter.PageSize
	if int(total)%filter.PageSize > 0 {
		totalPages++
	}

	return &GradeListResponse{
		Grades: responses,
		Pagination: PaginationMeta{
			Page:       filter.Page,
			PageSize:   filter.PageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// UpdateGrade updates a grade record
func (s *service) UpdateGrade(ctx context.Context, gradeID uint, req UpdateGradeRequest) (*GradeResponse, error) {
	grade, err := s.repo.FindByID(ctx, gradeID)
	if err != nil {
		return nil, err
	}

	if req.Title != "" {
		grade.Title = req.Title
	}
	if req.Score >= 0 && req.Score <= 100 {
		grade.Score = req.Score
	}
	if req.Description != "" {
		grade.Description = req.Description
	}

	if err := s.repo.Update(ctx, grade); err != nil {
		return nil, err
	}

	return toGradeResponse(grade), nil
}

// DeleteGrade deletes a grade record
func (s *service) DeleteGrade(ctx context.Context, gradeID uint) error {
	return s.repo.Delete(ctx, gradeID)
}

// ValidateTeacherAccess validates that a teacher can access a student's grades
// Requirements: 10.5 - THE System SHALL validate that Wali_Kelas can only input grades for students in their assigned class
func (s *service) ValidateTeacherAccess(ctx context.Context, teacherID, studentID uint) error {
	// Get teacher's assigned class
	classID, err := s.GetTeacherClassID(ctx, teacherID)
	if err != nil {
		return err
	}
	if classID == nil {
		return ErrNoClassAssigned
	}

	// Get student's class
	student, err := s.repo.FindStudentByID(ctx, studentID)
	if err != nil {
		return err
	}

	// Check if student is in teacher's class (handle nullable ClassID)
	if student.ClassID == nil || *student.ClassID != *classID {
		return ErrStudentNotInClass
	}

	return nil
}

// GetTeacherClassID returns the class ID assigned to a wali kelas
func (s *service) GetTeacherClassID(ctx context.Context, teacherID uint) (*uint, error) {
	var class models.Class
	err := s.db.WithContext(ctx).
		Where("homeroom_teacher_id = ?", teacherID).
		First(&class).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &class.ID, nil
}

// GetStudentGradeSummary retrieves grade summary for a student
func (s *service) GetStudentGradeSummary(ctx context.Context, studentID uint) (*StudentGradeSummary, error) {
	return s.repo.GetStudentGradeSummary(ctx, studentID)
}

// GetClassGrades retrieves grades for all students in a class
func (s *service) GetClassGrades(ctx context.Context, classID uint, filter GradeFilter) (*GradeListResponse, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}

	grades, total, err := s.repo.GetClassGrades(ctx, classID, filter)
	if err != nil {
		return nil, err
	}

	responses := make([]GradeResponse, len(grades))
	for i, g := range grades {
		responses[i] = *toGradeResponse(&g)
	}

	totalPages := int(total) / filter.PageSize
	if int(total)%filter.PageSize > 0 {
		totalPages++
	}

	return &GradeListResponse{
		Grades: responses,
		Pagination: PaginationMeta{
			Page:       filter.Page,
			PageSize:   filter.PageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// ==================== Response Converters ====================

func toGradeResponse(g *models.Grade) *GradeResponse {
	response := &GradeResponse{
		ID:          g.ID,
		StudentID:   g.StudentID,
		Title:       g.Title,
		Score:       g.Score,
		Description: g.Description,
		CreatedBy:   g.CreatedBy,
		CreatedAt:   g.CreatedAt,
		UpdatedAt:   g.UpdatedAt,
	}

	if g.Student.ID != 0 {
		response.StudentName = g.Student.Name
		response.StudentNIS = g.Student.NIS
		response.StudentNISN = g.Student.NISN
		// Handle nullable Class pointer
		if g.Student.Class != nil && g.Student.Class.ID != 0 {
			response.ClassName = g.Student.Class.Name
		}
	}

	if g.Creator.ID != 0 {
		response.CreatorName = g.Creator.Username
	}

	return response
}
