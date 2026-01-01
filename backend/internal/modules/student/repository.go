package student

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/school-management/backend/internal/domain/models"
)

var (
	ErrStudentNotFound = errors.New("siswa tidak ditemukan")
	ErrUserNotStudent  = errors.New("user bukan siswa")
)

// Repository defines the interface for student data operations
type Repository interface {
	// Student operations
	FindStudentByUserID(ctx context.Context, userID uint) (*models.Student, error)
	FindStudentByID(ctx context.Context, studentID uint) (*models.Student, error)

	// Attendance operations
	GetStudentAttendance(ctx context.Context, studentID uint, startDate, endDate time.Time, page, pageSize int) ([]models.Attendance, int64, error)
	GetAttendanceSummary(ctx context.Context, studentID uint, startDate, endDate time.Time) (*AttendanceSummaryResponse, error)

	// Grade operations
	GetStudentGrades(ctx context.Context, studentID uint, page, pageSize int) ([]models.Grade, int64, error)
	GetGradeSummary(ctx context.Context, studentID uint) (*GradeSummaryResponse, error)

	// BK operations
	GetStudentViolations(ctx context.Context, studentID uint, limit int) ([]models.Violation, error)
	GetStudentAchievements(ctx context.Context, studentID uint, limit int) ([]models.Achievement, error)
	GetStudentAchievementPoints(ctx context.Context, studentID uint) (int, error)
	GetStudentBKCounts(ctx context.Context, studentID uint) (violations, achievements int, err error)
}

// repository implements the Repository interface
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new student repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// FindStudentByUserID retrieves a student by user ID
// Requirements: 16.1 - Student logs in using NISN and password
func (r *repository) FindStudentByUserID(ctx context.Context, userID uint) (*models.Student, error) {
	var student models.Student
	err := r.db.WithContext(ctx).
		Preload("Class").
		Where("user_id = ?", userID).
		First(&student).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrStudentNotFound
		}
		return nil, err
	}

	return &student, nil
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
// Requirements: 16.3 - Student views attendance history
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
		TotalDays: int(totalDays),
		Present:   int(present),
		Late:      int(late),
		VeryLate:  int(veryLate),
		Absent:    int(absent),
	}, nil
}

// GetStudentGrades retrieves grades for a student
// Requirements: 16.4 - Student views all grades
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
		TotalGrades:  int(result.Count),
		AverageScore: result.Average,
		HighestScore: result.Highest,
		LowestScore:  result.Lowest,
	}, nil
}

// GetStudentViolations retrieves violations for a student
// Requirements: 16.5 - Student views violations (summary only)
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
// Requirements: 16.5 - Student views achievements
func (r *repository) GetStudentAchievements(ctx context.Context, studentID uint, limit int) ([]models.Achievement, error) {
	var achievements []models.Achievement
	err := r.db.WithContext(ctx).
		Where("student_id = ?", studentID).
		Order("created_at DESC").
		Limit(limit).
		Find(&achievements).Error

	return achievements, err
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
func (r *repository) GetStudentBKCounts(ctx context.Context, studentID uint) (violations, achievements int, err error) {
	var v, a int64

	r.db.WithContext(ctx).Model(&models.Violation{}).Where("student_id = ?", studentID).Count(&v)
	r.db.WithContext(ctx).Model(&models.Achievement{}).Where("student_id = ?", studentID).Count(&a)

	return int(v), int(a), nil
}
