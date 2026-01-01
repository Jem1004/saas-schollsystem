package grade

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/school-management/backend/internal/domain/models"
)

var (
	ErrGradeNotFound   = errors.New("grade not found")
	ErrStudentNotFound = errors.New("student not found")
	ErrUserNotFound    = errors.New("user not found")
)

// Repository defines the interface for Grade data operations
type Repository interface {
	// Grade operations
	Create(ctx context.Context, grade *models.Grade) error
	FindByID(ctx context.Context, id uint) (*models.Grade, error)
	FindByStudent(ctx context.Context, studentID uint) ([]models.Grade, error)
	FindAll(ctx context.Context, schoolID uint, filter GradeFilter) ([]models.Grade, int64, error)
	Update(ctx context.Context, grade *models.Grade) error
	Delete(ctx context.Context, id uint) error

	// Student lookup
	FindStudentByID(ctx context.Context, studentID uint) (*models.Student, error)
	FindUserByID(ctx context.Context, userID uint) (*models.User, error)

	// Summary operations
	GetStudentGradeSummary(ctx context.Context, studentID uint) (*StudentGradeSummary, error)
	GetClassGrades(ctx context.Context, classID uint, filter GradeFilter) ([]models.Grade, int64, error)
}

// repository implements the Repository interface
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new Grade repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Create creates a new grade record
// Requirements: 10.1 - WHEN a Wali_Kelas inputs a grade
func (r *repository) Create(ctx context.Context, grade *models.Grade) error {
	return r.db.WithContext(ctx).Create(grade).Error
}

// FindByID retrieves a grade by ID
func (r *repository) FindByID(ctx context.Context, id uint) (*models.Grade, error) {
	var grade models.Grade
	err := r.db.WithContext(ctx).
		Preload("Student").
		Preload("Student.Class").
		Preload("Creator").
		Where("id = ?", id).
		First(&grade).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrGradeNotFound
		}
		return nil, err
	}
	return &grade, nil
}


// FindByStudent retrieves all grades for a student sorted by date
// Requirements: 10.4 - THE System SHALL display all grades for their child sorted by date
func (r *repository) FindByStudent(ctx context.Context, studentID uint) ([]models.Grade, error) {
	var grades []models.Grade
	err := r.db.WithContext(ctx).
		Preload("Student").
		Preload("Student.Class").
		Preload("Creator").
		Where("student_id = ?", studentID).
		Order("created_at DESC").
		Find(&grades).Error
	return grades, err
}

// FindAll retrieves grades with pagination and filtering
func (r *repository) FindAll(ctx context.Context, schoolID uint, filter GradeFilter) ([]models.Grade, int64, error) {
	var grades []models.Grade
	var total int64

	query := r.db.WithContext(ctx).
		Model(&models.Grade{}).
		Joins("JOIN students ON students.id = grades.student_id").
		Where("students.school_id = ?", schoolID)

	// Apply filters
	if filter.StudentID != nil {
		query = query.Where("grades.student_id = ?", *filter.StudentID)
	}
	if filter.ClassID != nil {
		query = query.Where("students.class_id = ?", *filter.ClassID)
	}
	if filter.StartDate != nil {
		startDate, err := time.Parse("2006-01-02", *filter.StartDate)
		if err == nil {
			query = query.Where("grades.created_at >= ?", startDate)
		}
	}
	if filter.EndDate != nil {
		endDate, err := time.Parse("2006-01-02", *filter.EndDate)
		if err == nil {
			query = query.Where("grades.created_at <= ?", endDate.Add(24*time.Hour))
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
		Preload("Creator").
		Joins("JOIN students ON students.id = grades.student_id").
		Where("students.school_id = ?", schoolID).
		Order("grades.created_at DESC").
		Offset(offset).
		Limit(filter.PageSize).
		Find(&grades).Error

	if err != nil {
		return nil, 0, err
	}

	return grades, total, nil
}

// Update updates a grade record
func (r *repository) Update(ctx context.Context, grade *models.Grade) error {
	result := r.db.WithContext(ctx).
		Model(&models.Grade{}).
		Where("id = ?", grade.ID).
		Updates(map[string]interface{}{
			"student_id":  grade.StudentID,
			"title":       grade.Title,
			"score":       grade.Score,
			"description": grade.Description,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrGradeNotFound
	}
	return nil
}

// Delete deletes a grade record
func (r *repository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.Grade{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrGradeNotFound
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

// GetStudentGradeSummary retrieves grade summary for a student
func (r *repository) GetStudentGradeSummary(ctx context.Context, studentID uint) (*StudentGradeSummary, error) {
	student, err := r.FindStudentByID(ctx, studentID)
	if err != nil {
		return nil, err
	}

	var result struct {
		TotalGrades  int
		AverageScore float64
	}

	err = r.db.WithContext(ctx).
		Model(&models.Grade{}).
		Where("student_id = ?", studentID).
		Select("COUNT(*) as total_grades, COALESCE(AVG(score), 0) as average_score").
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	className := ""
	if student.Class.ID != 0 {
		className = student.Class.Name
	}

	return &StudentGradeSummary{
		StudentID:    studentID,
		StudentName:  student.Name,
		ClassName:    className,
		TotalGrades:  result.TotalGrades,
		AverageScore: result.AverageScore,
	}, nil
}

// GetClassGrades retrieves grades for all students in a class
func (r *repository) GetClassGrades(ctx context.Context, classID uint, filter GradeFilter) ([]models.Grade, int64, error) {
	var grades []models.Grade
	var total int64

	query := r.db.WithContext(ctx).
		Model(&models.Grade{}).
		Joins("JOIN students ON students.id = grades.student_id").
		Where("students.class_id = ?", classID)

	// Apply date filters
	if filter.StartDate != nil {
		startDate, err := time.Parse("2006-01-02", *filter.StartDate)
		if err == nil {
			query = query.Where("grades.created_at >= ?", startDate)
		}
	}
	if filter.EndDate != nil {
		endDate, err := time.Parse("2006-01-02", *filter.EndDate)
		if err == nil {
			query = query.Where("grades.created_at <= ?", endDate.Add(24*time.Hour))
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
		Preload("Creator").
		Joins("JOIN students ON students.id = grades.student_id").
		Where("students.class_id = ?", classID).
		Order("grades.created_at DESC").
		Offset(offset).
		Limit(filter.PageSize).
		Find(&grades).Error

	if err != nil {
		return nil, 0, err
	}

	return grades, total, nil
}
