package homeroom

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/school-management/backend/internal/domain/models"
)

var (
	ErrStudentIDRequired  = errors.New("student_id is required")
	ErrContentRequired    = errors.New("content is required")
	ErrStudentNotInSchool = errors.New("student does not belong to this school")
	ErrStudentNotInClass  = errors.New("student does not belong to your assigned class")
	ErrNotAuthorized      = errors.New("not authorized to perform this action")
	ErrNoClassAssigned    = errors.New("no class assigned to this teacher")
)

// Service defines the interface for Homeroom Note business logic
type Service interface {
	// Note operations
	CreateNote(ctx context.Context, schoolID, teacherID uint, req CreateNoteRequest) (*NoteResponse, error)
	GetNoteByID(ctx context.Context, id uint) (*NoteResponse, error)
	GetStudentNotes(ctx context.Context, studentID uint) ([]NoteResponse, error)
	GetNotes(ctx context.Context, schoolID uint, filter NoteFilter) (*NoteListResponse, error)
	UpdateNote(ctx context.Context, noteID uint, req UpdateNoteRequest) (*NoteResponse, error)
	DeleteNote(ctx context.Context, noteID uint) error

	// Teacher validation
	ValidateTeacherAccess(ctx context.Context, teacherID, studentID uint) error
	GetTeacherClassID(ctx context.Context, teacherID uint) (*uint, error)

	// Summary
	GetStudentNoteSummary(ctx context.Context, studentID uint) (*StudentNoteSummary, error)
	GetClassNotes(ctx context.Context, classID uint, filter NoteFilter) (*NoteListResponse, error)
}

// service implements the Service interface
type service struct {
	repo Repository
	db   *gorm.DB
}

// NewService creates a new Homeroom service
func NewService(repo Repository, db *gorm.DB) Service {
	return &service{repo: repo, db: db}
}

// CreateNote creates a new homeroom note
// Requirements: 11.1 - WHEN a Wali_Kelas creates a note, THE System SHALL require content and associate it with a student
// Requirements: 11.4 - THE System SHALL validate that Wali_Kelas can only create notes for students in their assigned class
func (s *service) CreateNote(ctx context.Context, schoolID, teacherID uint, req CreateNoteRequest) (*NoteResponse, error) {
	// Validate required fields
	if req.StudentID == 0 {
		return nil, ErrStudentIDRequired
	}
	if req.Content == "" {
		return nil, ErrContentRequired
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

	note := &models.HomeroomNote{
		StudentID: req.StudentID,
		TeacherID: teacherID,
		Content:   req.Content,
	}

	if err := s.repo.Create(ctx, note); err != nil {
		return nil, err
	}

	// Reload with relations
	note, err = s.repo.FindByID(ctx, note.ID)
	if err != nil {
		return nil, err
	}

	// TODO: Trigger notification to parent (async)
	// Requirements: 11.2 - WHEN a note is saved, THE System SHALL trigger notification to the parent

	return toNoteResponse(note), nil
}


// GetNoteByID retrieves a homeroom note by ID
func (s *service) GetNoteByID(ctx context.Context, id uint) (*NoteResponse, error) {
	note, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toNoteResponse(note), nil
}

// GetStudentNotes retrieves all homeroom notes for a student
// Requirements: 11.3 - WHEN parents view notes, THE System SHALL display all homeroom notes for their child
// Requirements: 11.5 - THE System SHALL maintain note history with timestamps and author information
func (s *service) GetStudentNotes(ctx context.Context, studentID uint) ([]NoteResponse, error) {
	notes, err := s.repo.FindByStudent(ctx, studentID)
	if err != nil {
		return nil, err
	}

	responses := make([]NoteResponse, len(notes))
	for i, n := range notes {
		responses[i] = *toNoteResponse(&n)
	}
	return responses, nil
}

// GetNotes retrieves homeroom notes with pagination and filtering
func (s *service) GetNotes(ctx context.Context, schoolID uint, filter NoteFilter) (*NoteListResponse, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}

	notes, total, err := s.repo.FindAll(ctx, schoolID, filter)
	if err != nil {
		return nil, err
	}

	responses := make([]NoteResponse, len(notes))
	for i, n := range notes {
		responses[i] = *toNoteResponse(&n)
	}

	totalPages := int(total) / filter.PageSize
	if int(total)%filter.PageSize > 0 {
		totalPages++
	}

	return &NoteListResponse{
		Notes: responses,
		Pagination: PaginationMeta{
			Page:       filter.Page,
			PageSize:   filter.PageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// UpdateNote updates a homeroom note
func (s *service) UpdateNote(ctx context.Context, noteID uint, req UpdateNoteRequest) (*NoteResponse, error) {
	note, err := s.repo.FindByID(ctx, noteID)
	if err != nil {
		return nil, err
	}

	if req.Content != "" {
		note.Content = req.Content
	}

	if err := s.repo.Update(ctx, note); err != nil {
		return nil, err
	}

	return toNoteResponse(note), nil
}

// DeleteNote deletes a homeroom note
func (s *service) DeleteNote(ctx context.Context, noteID uint) error {
	return s.repo.Delete(ctx, noteID)
}

// ValidateTeacherAccess validates that a teacher can access a student's notes
// Requirements: 11.4 - THE System SHALL validate that Wali_Kelas can only create notes for students in their assigned class
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

	// Check if student is in teacher's class
	if student.ClassID != *classID {
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

// GetStudentNoteSummary retrieves note summary for a student
func (s *service) GetStudentNoteSummary(ctx context.Context, studentID uint) (*StudentNoteSummary, error) {
	return s.repo.GetStudentNoteSummary(ctx, studentID)
}

// GetClassNotes retrieves homeroom notes for all students in a class
func (s *service) GetClassNotes(ctx context.Context, classID uint, filter NoteFilter) (*NoteListResponse, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}

	notes, total, err := s.repo.GetClassNotes(ctx, classID, filter)
	if err != nil {
		return nil, err
	}

	responses := make([]NoteResponse, len(notes))
	for i, n := range notes {
		responses[i] = *toNoteResponse(&n)
	}

	totalPages := int(total) / filter.PageSize
	if int(total)%filter.PageSize > 0 {
		totalPages++
	}

	return &NoteListResponse{
		Notes: responses,
		Pagination: PaginationMeta{
			Page:       filter.Page,
			PageSize:   filter.PageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// ==================== Response Converters ====================

func toNoteResponse(n *models.HomeroomNote) *NoteResponse {
	response := &NoteResponse{
		ID:        n.ID,
		StudentID: n.StudentID,
		TeacherID: n.TeacherID,
		Content:   n.Content,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
	}

	if n.Student.ID != 0 {
		response.StudentName = n.Student.Name
		response.StudentNIS = n.Student.NIS
		response.StudentNISN = n.Student.NISN
		if n.Student.Class.ID != 0 {
			response.ClassName = n.Student.Class.Name
		}
	}

	if n.Teacher.ID != 0 {
		response.TeacherName = n.Teacher.Username
	}

	return response
}
