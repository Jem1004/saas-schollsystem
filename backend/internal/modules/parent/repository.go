package parent

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/school-management/backend/internal/domain/models"
)

var (
	ErrParentNotFound  = errors.New("parent not found")
	ErrStudentNotFound = errors.New("student not found")
	ErrNotLinked       = errors.New("student is not linked to this parent")
)

// Repository defines the interface for parent data operations
type Repository interface {
	// Parent operations
	FindParentByUserID(ctx context.Context, userID uint) (*models.Parent, error)
	GetLinkedStudents(ctx context.Context, parentID uint) ([]models.Student, error)
	IsStudentLinked(ctx context.Context, parentID, studentID uint) (bool, error)

	// Student data operations
	FindStudentByID(ctx context.Context, studentID uint) (*models.Student, error)

	// Attendance operations
	GetStudentAttendance(ctx context.Context, studentID uint, startDate, endDate time.Time, page, pageSize int) ([]models.Attendance, int64, error)
	GetAttendanceSummary(ctx context.Context, studentID uint, startDate, endDate time.Time) (*AttendanceSummaryResponse, error)

	// Grade operations
	GetStudentGrades(ctx context.Context, studentID uint, page, pageSize int) ([]models.Grade, int64, error)
	GetGradeSummary(ctx context.Context, studentID uint) (*GradeSummaryResponse, error)

	// Homeroom note operations
	GetStudentNotes(ctx context.Context, studentID uint, page, pageSize int) ([]models.HomeroomNote, int64, error)

	// BK operations
	GetStudentViolations(ctx context.Context, studentID uint, limit int) ([]models.Violation, error)
	GetStudentAchievements(ctx context.Context, studentID uint, limit int) ([]models.Achievement, error)
	GetStudentPermits(ctx context.Context, studentID uint, limit int) ([]models.Permit, error)
	GetStudentCounselingNotes(ctx context.Context, studentID uint, limit int) ([]models.CounselingNote, error)
	GetStudentAchievementPoints(ctx context.Context, studentID uint) (int, error)
	GetStudentBKCounts(ctx context.Context, studentID uint) (violations, achievements, permits, counseling int, err error)
}

// repository implements the Repository interface
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new parent repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// FindParentByUserID retrieves a parent by user ID
func (r *repository) FindParentByUserID(ctx context.Context, userID uint) (*models.Parent, error) {
	var parent models.Parent
	err := r.db.WithContext(ctx).
		Preload("Students").
		Preload("Students.Class").
		Where("user_id = ?", userID).
		First(&parent).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrParentNotFound
		}
		return nil, err
	}

	return &parent, nil
}

// GetLinkedStudents retrieves all students linked to a parent
// Requirements: 12.2 - Authentication succeeds SHALL return access to all linked children's data
func (r *repository) GetLinkedStudents(ctx context.Context, parentID uint) ([]models.Student, error) {
	var parent models.Parent
	err := r.db.WithContext(ctx).
		Preload("Students", func(db *gorm.DB) *gorm.DB {
			return db.Order("name ASC")
		}).
		Preload("Students.Class").
		First(&parent, parentID).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrParentNotFound
		}
		return nil, err
	}

	return parent.Students, nil
}

// IsStudentLinked checks if a student is linked to a parent
func (r *repository) IsStudentLinked(ctx context.Context, parentID, studentID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table("student_parents").
		Where("parent_id = ? AND student_id = ?", parentID, studentID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// FindStudentByID retrieves a student by ID
func (r *repository) FindStudentByID(ctx context.Context, studentID uint) (*models.Student, error) {
	var student models.Student
	err := r.db.WithContext(ctx).
		Preload("Class").
		First(&student, studentID).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrStudentNotFound
		}
		return nil, err
	}

	return &student, nil
}

// GetStudentAttendance retrieves attendance records for a student
// Requirements: 15.1 - Parent opens grades section SHALL display all grades
func (r *repository) GetStudentAttendance(ctx context.Context, studentID uint, startDate, endDate time.Time, page, pageSize int) ([]models.Attendance, int64, error) {
	var attendances []models.Attendance
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Attendance{}).Where("student_id = ?", studentID)

	if !startDate.IsZero() {
		query = query.Where("date >= ?", startDate)
	}
	if !endDate.IsZero() {
		query = query.Where("date <= ?", endDate)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.
		Order("date DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&attendances).Error

	if err != nil {
		return nil, 0, err
	}

	return attendances, total, nil
}

// GetAttendanceSummary retrieves attendance summary for a student
func (r *repository) GetAttendanceSummary(ctx context.Context, studentID uint, startDate, endDate time.Time) (*AttendanceSummaryResponse, error) {
	var student models.Student
	if err := r.db.WithContext(ctx).First(&student, studentID).Error; err != nil {
		return nil, err
	}

	query := r.db.WithContext(ctx).Model(&models.Attendance{}).Where("student_id = ?", studentID)

	if !startDate.IsZero() {
		query = query.Where("date >= ?", startDate)
	}
	if !endDate.IsZero() {
		query = query.Where("date <= ?", endDate)
	}

	var totalDays int64
	query.Count(&totalDays)

	var present, late, veryLate, absent int64
	r.db.WithContext(ctx).Model(&models.Attendance{}).
		Where("student_id = ? AND status = ?", studentID, models.AttendanceStatusOnTime).
		Count(&present)
	r.db.WithContext(ctx).Model(&models.Attendance{}).
		Where("student_id = ? AND status = ?", studentID, models.AttendanceStatusLate).
		Count(&late)
	r.db.WithContext(ctx).Model(&models.Attendance{}).
		Where("student_id = ? AND status = ?", studentID, models.AttendanceStatusVeryLate).
		Count(&veryLate)
	r.db.WithContext(ctx).Model(&models.Attendance{}).
		Where("student_id = ? AND status = ?", studentID, models.AttendanceStatusAbsent).
		Count(&absent)

	return &AttendanceSummaryResponse{
		StudentID:   studentID,
		StudentName: student.Name,
		TotalDays:   int(totalDays),
		Present:     int(present),
		Late:        int(late),
		VeryLate:    int(veryLate),
		Absent:      int(absent),
	}, nil
}

// GetStudentGrades retrieves grades for a student
// Requirements: 15.1 - Parent opens grades section SHALL display all grades sorted by date
func (r *repository) GetStudentGrades(ctx context.Context, studentID uint, page, pageSize int) ([]models.Grade, int64, error) {
	var grades []models.Grade
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Grade{}).Where("student_id = ?", studentID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.
		Preload("Creator").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&grades).Error

	if err != nil {
		return nil, 0, err
	}

	return grades, total, nil
}

// GetGradeSummary retrieves grade summary for a student
func (r *repository) GetGradeSummary(ctx context.Context, studentID uint) (*GradeSummaryResponse, error) {
	var student models.Student
	if err := r.db.WithContext(ctx).First(&student, studentID).Error; err != nil {
		return nil, err
	}

	var result struct {
		Count   int64
		Average float64
		Highest float64
		Lowest  float64
	}

	err := r.db.WithContext(ctx).Model(&models.Grade{}).
		Select("COUNT(*) as count, AVG(score) as average, MAX(score) as highest, MIN(score) as lowest").
		Where("student_id = ?", studentID).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return &GradeSummaryResponse{
		StudentID:    studentID,
		StudentName:  student.Name,
		TotalGrades:  int(result.Count),
		AverageScore: result.Average,
		HighestScore: result.Highest,
		LowestScore:  result.Lowest,
	}, nil
}

// GetStudentNotes retrieves homeroom notes for a student
// Requirements: 15.2 - Parent opens notes section SHALL display all homeroom notes
func (r *repository) GetStudentNotes(ctx context.Context, studentID uint, page, pageSize int) ([]models.HomeroomNote, int64, error) {
	var notes []models.HomeroomNote
	var total int64

	query := r.db.WithContext(ctx).Model(&models.HomeroomNote{}).Where("student_id = ?", studentID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.
		Preload("Teacher").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&notes).Error

	if err != nil {
		return nil, 0, err
	}

	return notes, total, nil
}

// GetStudentViolations retrieves violations for a student
func (r *repository) GetStudentViolations(ctx context.Context, studentID uint, limit int) ([]models.Violation, error) {
	var violations []models.Violation
	err := r.db.WithContext(ctx).
		Where("student_id = ?", studentID).
		Order("created_at DESC").
		Limit(limit).
		Find(&violations).Error

	return violations, err
}

// GetStudentAchievements retrieves achievements for a student
func (r *repository) GetStudentAchievements(ctx context.Context, studentID uint, limit int) ([]models.Achievement, error) {
	var achievements []models.Achievement
	err := r.db.WithContext(ctx).
		Where("student_id = ?", studentID).
		Order("created_at DESC").
		Limit(limit).
		Find(&achievements).Error

	return achievements, err
}

// GetStudentPermits retrieves permits for a student
func (r *repository) GetStudentPermits(ctx context.Context, studentID uint, limit int) ([]models.Permit, error) {
	var permits []models.Permit
	err := r.db.WithContext(ctx).
		Preload("Teacher").
		Where("student_id = ?", studentID).
		Order("created_at DESC").
		Limit(limit).
		Find(&permits).Error

	return permits, err
}

// GetStudentCounselingNotes retrieves counseling notes for a student (parent summary only)
// Requirements: 9.2, 14.5 - Only parent_summary visible to parents
func (r *repository) GetStudentCounselingNotes(ctx context.Context, studentID uint, limit int) ([]models.CounselingNote, error) {
	var notes []models.CounselingNote
	err := r.db.WithContext(ctx).
		Select("id, student_id, parent_summary, created_by, created_at").
		Where("student_id = ? AND parent_summary != ''", studentID).
		Order("created_at DESC").
		Limit(limit).
		Find(&notes).Error

	return notes, err
}

// GetStudentAchievementPoints retrieves total achievement points for a student
func (r *repository) GetStudentAchievementPoints(ctx context.Context, studentID uint) (int, error) {
	var total int
	err := r.db.WithContext(ctx).Model(&models.Achievement{}).
		Select("COALESCE(SUM(point), 0)").
		Where("student_id = ?", studentID).
		Scan(&total).Error

	return total, err
}

// GetStudentBKCounts retrieves BK counts for a student
func (r *repository) GetStudentBKCounts(ctx context.Context, studentID uint) (violations, achievements, permits, counseling int, err error) {
	var v, a, p, c int64

	r.db.WithContext(ctx).Model(&models.Violation{}).Where("student_id = ?", studentID).Count(&v)
	r.db.WithContext(ctx).Model(&models.Achievement{}).Where("student_id = ?", studentID).Count(&a)
	r.db.WithContext(ctx).Model(&models.Permit{}).Where("student_id = ?", studentID).Count(&p)
	r.db.WithContext(ctx).Model(&models.CounselingNote{}).Where("student_id = ? AND parent_summary != ''", studentID).Count(&c)

	return int(v), int(a), int(p), int(c), nil
}
