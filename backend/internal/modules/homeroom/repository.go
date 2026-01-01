package homeroom

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/school-management/backend/internal/domain/models"
)

var (
	ErrNoteNotFound    = errors.New("homeroom note not found")
	ErrStudentNotFound = errors.New("student not found")
	ErrUserNotFound    = errors.New("user not found")
)

// Repository defines the interface for Homeroom Note data operations
type Repository interface {
	// Note operations
	Create(ctx context.Context, note *models.HomeroomNote) error
	FindByID(ctx context.Context, id uint) (*models.HomeroomNote, error)
	FindByStudent(ctx context.Context, studentID uint) ([]models.HomeroomNote, error)
	FindAll(ctx context.Context, schoolID uint, filter NoteFilter) ([]models.HomeroomNote, int64, error)
	Update(ctx context.Context, note *models.HomeroomNote) error
	Delete(ctx context.Context, id uint) error

	// Student lookup
	FindStudentByID(ctx context.Context, studentID uint) (*models.Student, error)
	FindUserByID(ctx context.Context, userID uint) (*models.User, error)

	// Summary operations
	GetStudentNoteSummary(ctx context.Context, studentID uint) (*StudentNoteSummary, error)
	GetClassNotes(ctx context.Context, classID uint, filter NoteFilter) ([]models.HomeroomNote, int64, error)
}

// repository implements the Repository interface
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new Homeroom repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Create creates a new homeroom note
// Requirements: 11.1 - WHEN a Wali_Kelas creates a note
func (r *repository) Create(ctx context.Context, note *models.HomeroomNote) error {
	return r.db.WithContext(ctx).Create(note).Error
}

// FindByID retrieves a homeroom note by ID
func (r *repository) FindByID(ctx context.Context, id uint) (*models.HomeroomNote, error) {
	var note models.HomeroomNote
	err := r.db.WithContext(ctx).
		Preload("Student").
		Preload("Student.Class").
		Preload("Teacher").
		Where("id = ?", id).
		First(&note).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNoteNotFound
		}
		return nil, err
	}
	return &note, nil
}

// FindByStudent retrieves all homeroom notes for a student
// Requirements: 11.3 - WHEN parents view notes, THE System SHALL display all homeroom notes for their child
// Requirements: 11.5 - THE System SHALL maintain note history with timestamps and author information
func (r *repository) FindByStudent(ctx context.Context, studentID uint) ([]models.HomeroomNote, error) {
	var notes []models.HomeroomNote
	err := r.db.WithContext(ctx).
		Preload("Student").
		Preload("Student.Class").
		Preload("Teacher").
		Where("student_id = ?", studentID).
		Order("created_at DESC").
		Find(&notes).Error
	return notes, err
}


// FindAll retrieves homeroom notes with pagination and filtering
func (r *repository) FindAll(ctx context.Context, schoolID uint, filter NoteFilter) ([]models.HomeroomNote, int64, error) {
	var notes []models.HomeroomNote
	var total int64

	query := r.db.WithContext(ctx).
		Model(&models.HomeroomNote{}).
		Joins("JOIN students ON students.id = homeroom_notes.student_id").
		Where("students.school_id = ?", schoolID)

	// Apply filters
	if filter.StudentID != nil {
		query = query.Where("homeroom_notes.student_id = ?", *filter.StudentID)
	}
	if filter.ClassID != nil {
		query = query.Where("students.class_id = ?", *filter.ClassID)
	}
	if filter.TeacherID != nil {
		query = query.Where("homeroom_notes.teacher_id = ?", *filter.TeacherID)
	}
	if filter.StartDate != nil {
		startDate, err := time.Parse("2006-01-02", *filter.StartDate)
		if err == nil {
			query = query.Where("homeroom_notes.created_at >= ?", startDate)
		}
	}
	if filter.EndDate != nil {
		endDate, err := time.Parse("2006-01-02", *filter.EndDate)
		if err == nil {
			query = query.Where("homeroom_notes.created_at <= ?", endDate.Add(24*time.Hour))
		}
	}

	// Count total
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
	err := r.db.WithContext(ctx).
		Preload("Student").
		Preload("Student.Class").
		Preload("Teacher").
		Joins("JOIN students ON students.id = homeroom_notes.student_id").
		Where("students.school_id = ?", schoolID).
		Order("homeroom_notes.created_at DESC").
		Offset(offset).
		Limit(filter.PageSize).
		Find(&notes).Error

	if err != nil {
		return nil, 0, err
	}

	return notes, total, nil
}

// Update updates a homeroom note
func (r *repository) Update(ctx context.Context, note *models.HomeroomNote) error {
	result := r.db.WithContext(ctx).
		Model(&models.HomeroomNote{}).
		Where("id = ?", note.ID).
		Updates(map[string]interface{}{
			"student_id": note.StudentID,
			"teacher_id": note.TeacherID,
			"content":    note.Content,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNoteNotFound
	}
	return nil
}

// Delete deletes a homeroom note
func (r *repository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.HomeroomNote{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNoteNotFound
	}
	return nil
}

// FindStudentByID retrieves a student by ID
func (r *repository) FindStudentByID(ctx context.Context, studentID uint) (*models.Student, error) {
	var student models.Student
	err := r.db.WithContext(ctx).
		Preload("Class").
		Preload("School").
		Where("id = ?", studentID).
		First(&student).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrStudentNotFound
		}
		return nil, err
	}
	return &student, nil
}

// FindUserByID retrieves a user by ID
func (r *repository) FindUserByID(ctx context.Context, userID uint) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).
		Where("id = ?", userID).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// GetStudentNoteSummary retrieves note summary for a student
func (r *repository) GetStudentNoteSummary(ctx context.Context, studentID uint) (*StudentNoteSummary, error) {
	student, err := r.FindStudentByID(ctx, studentID)
	if err != nil {
		return nil, err
	}

	var result struct {
		TotalNotes int
		LastNoteAt *time.Time
	}

	err = r.db.WithContext(ctx).
		Model(&models.HomeroomNote{}).
		Where("student_id = ?", studentID).
		Select("COUNT(*) as total_notes, MAX(created_at) as last_note_at").
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	className := ""
	if student.Class.ID != 0 {
		className = student.Class.Name
	}

	return &StudentNoteSummary{
		StudentID:   studentID,
		StudentName: student.Name,
		ClassName:   className,
		TotalNotes:  result.TotalNotes,
		LastNoteAt:  result.LastNoteAt,
	}, nil
}

// GetClassNotes retrieves homeroom notes for all students in a class
func (r *repository) GetClassNotes(ctx context.Context, classID uint, filter NoteFilter) ([]models.HomeroomNote, int64, error) {
	var notes []models.HomeroomNote
	var total int64

	query := r.db.WithContext(ctx).
		Model(&models.HomeroomNote{}).
		Joins("JOIN students ON students.id = homeroom_notes.student_id").
		Where("students.class_id = ?", classID)

	// Apply date filters
	if filter.StartDate != nil {
		startDate, err := time.Parse("2006-01-02", *filter.StartDate)
		if err == nil {
			query = query.Where("homeroom_notes.created_at >= ?", startDate)
		}
	}
	if filter.EndDate != nil {
		endDate, err := time.Parse("2006-01-02", *filter.EndDate)
		if err == nil {
			query = query.Where("homeroom_notes.created_at <= ?", endDate.Add(24*time.Hour))
		}
	}

	// Count total
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
	err := r.db.WithContext(ctx).
		Preload("Student").
		Preload("Student.Class").
		Preload("Teacher").
		Joins("JOIN students ON students.id = homeroom_notes.student_id").
		Where("students.class_id = ?", classID).
		Order("homeroom_notes.created_at DESC").
		Offset(offset).
		Limit(filter.PageSize).
		Find(&notes).Error

	if err != nil {
		return nil, 0, err
	}

	return notes, total, nil
}
