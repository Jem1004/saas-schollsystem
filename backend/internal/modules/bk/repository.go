package bk

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/school-management/backend/internal/domain/models"
)

var (
	ErrViolationNotFound      = errors.New("violation not found")
	ErrAchievementNotFound    = errors.New("achievement not found")
	ErrPermitNotFound         = errors.New("permit not found")
	ErrCounselingNoteNotFound = errors.New("counseling note not found")
	ErrStudentNotFound        = errors.New("student not found")
	ErrUserNotFound           = errors.New("user not found")
)

// Repository defines the interface for BK data operations
type Repository interface {
	// Violation operations
	CreateViolation(ctx context.Context, violation *models.Violation) error
	FindViolationByID(ctx context.Context, id uint) (*models.Violation, error)
	FindViolationsByStudent(ctx context.Context, studentID uint) ([]models.Violation, error)
	FindViolations(ctx context.Context, schoolID uint, filter ViolationFilter) ([]models.Violation, int64, error)
	DeleteViolation(ctx context.Context, id uint) error

	// Achievement operations
	CreateAchievement(ctx context.Context, achievement *models.Achievement) error
	FindAchievementByID(ctx context.Context, id uint) (*models.Achievement, error)
	FindAchievementsByStudent(ctx context.Context, studentID uint) ([]models.Achievement, error)
	FindAchievements(ctx context.Context, schoolID uint, filter AchievementFilter) ([]models.Achievement, int64, error)
	GetStudentAchievementPoints(ctx context.Context, studentID uint) (int, error)
	DeleteAchievement(ctx context.Context, id uint) error

	// Permit operations
	CreatePermit(ctx context.Context, permit *models.Permit) error
	FindPermitByID(ctx context.Context, id uint) (*models.Permit, error)
	FindPermitsByStudent(ctx context.Context, studentID uint) ([]models.Permit, error)
	FindPermits(ctx context.Context, schoolID uint, filter PermitFilter) ([]models.Permit, int64, error)
	UpdatePermit(ctx context.Context, permit *models.Permit) error
	DeletePermit(ctx context.Context, id uint) error

	// Counseling Note operations
	CreateCounselingNote(ctx context.Context, note *models.CounselingNote) error
	FindCounselingNoteByID(ctx context.Context, id uint) (*models.CounselingNote, error)
	FindCounselingNotesByStudent(ctx context.Context, studentID uint) ([]models.CounselingNote, error)
	FindCounselingNotes(ctx context.Context, schoolID uint, filter CounselingNoteFilter) ([]models.CounselingNote, int64, error)
	UpdateCounselingNote(ctx context.Context, note *models.CounselingNote) error
	DeleteCounselingNote(ctx context.Context, id uint) error

	// Student lookup
	FindStudentByID(ctx context.Context, studentID uint) (*models.Student, error)
	FindUserByID(ctx context.Context, userID uint) (*models.User, error)

	// Dashboard/Summary operations
	GetViolationCount(ctx context.Context, schoolID uint) (int64, error)
	GetAchievementCount(ctx context.Context, schoolID uint) (int64, error)
	GetPermitCount(ctx context.Context, schoolID uint) (int64, error)
	GetActivePermitCount(ctx context.Context, schoolID uint) (int64, error)
	GetCounselingCount(ctx context.Context, schoolID uint) (int64, error)
	GetStudentsNeedingAttention(ctx context.Context, schoolID uint, limit int) ([]StudentAttentionItem, error)
}

// repository implements the Repository interface
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new BK repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// ==================== Violation Repository ====================

// CreateViolation creates a new violation record
// Requirements: 6.1 - WHEN a Guru_BK records a violation
func (r *repository) CreateViolation(ctx context.Context, violation *models.Violation) error {
	return r.db.WithContext(ctx).Create(violation).Error
}

// FindViolationByID retrieves a violation by ID
func (r *repository) FindViolationByID(ctx context.Context, id uint) (*models.Violation, error) {
	var violation models.Violation
	err := r.db.WithContext(ctx).
		Preload("Student").
		Preload("Student.Class").
		Preload("Creator").
		Where("id = ?", id).
		First(&violation).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrViolationNotFound
		}
		return nil, err
	}
	return &violation, nil
}

// FindViolationsByStudent retrieves all violations for a student
// Requirements: 6.3 - THE System SHALL display all violations for a student sorted by date
func (r *repository) FindViolationsByStudent(ctx context.Context, studentID uint) ([]models.Violation, error) {
	var violations []models.Violation
	err := r.db.WithContext(ctx).
		Preload("Student").
		Preload("Student.Class").
		Preload("Creator").
		Where("student_id = ?", studentID).
		Order("created_at DESC").
		Find(&violations).Error
	return violations, err
}


// FindViolations retrieves violations with pagination and filtering
func (r *repository) FindViolations(ctx context.Context, schoolID uint, filter ViolationFilter) ([]models.Violation, int64, error) {
	var violations []models.Violation
	var total int64

	query := r.db.WithContext(ctx).
		Model(&models.Violation{}).
		Joins("JOIN students ON students.id = violations.student_id").
		Where("students.school_id = ?", schoolID)

	// Apply filters
	if filter.StudentID != nil {
		query = query.Where("violations.student_id = ?", *filter.StudentID)
	}
	if filter.ClassID != nil {
		query = query.Where("students.class_id = ?", *filter.ClassID)
	}
	if filter.Level != nil && *filter.Level != "" {
		query = query.Where("violations.level = ?", *filter.Level)
	}
	if filter.Category != nil && *filter.Category != "" {
		query = query.Where("violations.category ILIKE ?", "%"+*filter.Category+"%")
	}
	if filter.StartDate != nil {
		startDate, err := time.Parse("2006-01-02", *filter.StartDate)
		if err == nil {
			query = query.Where("violations.created_at >= ?", startDate)
		}
	}
	if filter.EndDate != nil {
		endDate, err := time.Parse("2006-01-02", *filter.EndDate)
		if err == nil {
			query = query.Where("violations.created_at <= ?", endDate.Add(24*time.Hour))
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
		Joins("JOIN students ON students.id = violations.student_id").
		Where("students.school_id = ?", schoolID).
		Order("violations.created_at DESC").
		Offset(offset).
		Limit(filter.PageSize).
		Find(&violations).Error

	if err != nil {
		return nil, 0, err
	}

	return violations, total, nil
}

// DeleteViolation deletes a violation record
func (r *repository) DeleteViolation(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.Violation{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrViolationNotFound
	}
	return nil
}

// ==================== Achievement Repository ====================

// CreateAchievement creates a new achievement record
// Requirements: 7.1 - WHEN a Guru_BK records an achievement
func (r *repository) CreateAchievement(ctx context.Context, achievement *models.Achievement) error {
	return r.db.WithContext(ctx).Create(achievement).Error
}

// FindAchievementByID retrieves an achievement by ID
func (r *repository) FindAchievementByID(ctx context.Context, id uint) (*models.Achievement, error) {
	var achievement models.Achievement
	err := r.db.WithContext(ctx).
		Preload("Student").
		Preload("Student.Class").
		Preload("Creator").
		Where("id = ?", id).
		First(&achievement).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAchievementNotFound
		}
		return nil, err
	}
	return &achievement, nil
}

// FindAchievementsByStudent retrieves all achievements for a student
// Requirements: 7.5 - THE System SHALL maintain achievement history per student
func (r *repository) FindAchievementsByStudent(ctx context.Context, studentID uint) ([]models.Achievement, error) {
	var achievements []models.Achievement
	err := r.db.WithContext(ctx).
		Preload("Student").
		Preload("Student.Class").
		Preload("Creator").
		Where("student_id = ?", studentID).
		Order("created_at DESC").
		Find(&achievements).Error
	return achievements, err
}

// FindAchievements retrieves achievements with pagination and filtering
func (r *repository) FindAchievements(ctx context.Context, schoolID uint, filter AchievementFilter) ([]models.Achievement, int64, error) {
	var achievements []models.Achievement
	var total int64

	query := r.db.WithContext(ctx).
		Model(&models.Achievement{}).
		Joins("JOIN students ON students.id = achievements.student_id").
		Where("students.school_id = ?", schoolID)

	// Apply filters
	if filter.StudentID != nil {
		query = query.Where("achievements.student_id = ?", *filter.StudentID)
	}
	if filter.ClassID != nil {
		query = query.Where("students.class_id = ?", *filter.ClassID)
	}
	if filter.StartDate != nil {
		startDate, err := time.Parse("2006-01-02", *filter.StartDate)
		if err == nil {
			query = query.Where("achievements.created_at >= ?", startDate)
		}
	}
	if filter.EndDate != nil {
		endDate, err := time.Parse("2006-01-02", *filter.EndDate)
		if err == nil {
			query = query.Where("achievements.created_at <= ?", endDate.Add(24*time.Hour))
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
		Joins("JOIN students ON students.id = achievements.student_id").
		Where("students.school_id = ?", schoolID).
		Order("achievements.created_at DESC").
		Offset(offset).
		Limit(filter.PageSize).
		Find(&achievements).Error

	if err != nil {
		return nil, 0, err
	}

	return achievements, total, nil
}

// GetStudentAchievementPoints retrieves total achievement points for a student
// Requirements: 7.2, 7.3 - THE System SHALL add points to student's accumulated achievement score
func (r *repository) GetStudentAchievementPoints(ctx context.Context, studentID uint) (int, error) {
	var total int
	err := r.db.WithContext(ctx).
		Model(&models.Achievement{}).
		Where("student_id = ?", studentID).
		Select("COALESCE(SUM(point), 0)").
		Scan(&total).Error
	return total, err
}

// DeleteAchievement deletes an achievement record
func (r *repository) DeleteAchievement(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.Achievement{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrAchievementNotFound
	}
	return nil
}


// ==================== Permit Repository ====================

// CreatePermit creates a new permit record
// Requirements: 8.1 - WHEN a Guru_BK creates an exit permit
func (r *repository) CreatePermit(ctx context.Context, permit *models.Permit) error {
	return r.db.WithContext(ctx).Create(permit).Error
}

// FindPermitByID retrieves a permit by ID
func (r *repository) FindPermitByID(ctx context.Context, id uint) (*models.Permit, error) {
	var permit models.Permit
	err := r.db.WithContext(ctx).
		Preload("Student").
		Preload("Student.Class").
		Preload("Student.School").
		Preload("Teacher").
		Preload("Creator").
		Where("id = ?", id).
		First(&permit).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPermitNotFound
		}
		return nil, err
	}
	return &permit, nil
}

// FindPermitsByStudent retrieves all permits for a student
func (r *repository) FindPermitsByStudent(ctx context.Context, studentID uint) ([]models.Permit, error) {
	var permits []models.Permit
	err := r.db.WithContext(ctx).
		Preload("Student").
		Preload("Student.Class").
		Preload("Teacher").
		Preload("Creator").
		Where("student_id = ?", studentID).
		Order("created_at DESC").
		Find(&permits).Error
	return permits, err
}

// FindPermits retrieves permits with pagination and filtering
func (r *repository) FindPermits(ctx context.Context, schoolID uint, filter PermitFilter) ([]models.Permit, int64, error) {
	var permits []models.Permit
	var total int64

	query := r.db.WithContext(ctx).
		Model(&models.Permit{}).
		Joins("JOIN students ON students.id = permits.student_id").
		Where("students.school_id = ?", schoolID)

	// Apply filters
	if filter.StudentID != nil {
		query = query.Where("permits.student_id = ?", *filter.StudentID)
	}
	if filter.ClassID != nil {
		query = query.Where("students.class_id = ?", *filter.ClassID)
	}
	if filter.TeacherID != nil {
		query = query.Where("permits.responsible_teacher = ?", *filter.TeacherID)
	}
	if filter.HasReturned != nil {
		if *filter.HasReturned {
			query = query.Where("permits.return_time IS NOT NULL")
		} else {
			query = query.Where("permits.return_time IS NULL")
		}
	}
	if filter.StartDate != nil {
		startDate, err := time.Parse("2006-01-02", *filter.StartDate)
		if err == nil {
			query = query.Where("permits.exit_time >= ?", startDate)
		}
	}
	if filter.EndDate != nil {
		endDate, err := time.Parse("2006-01-02", *filter.EndDate)
		if err == nil {
			query = query.Where("permits.exit_time <= ?", endDate.Add(24*time.Hour))
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
		Preload("Creator").
		Joins("JOIN students ON students.id = permits.student_id").
		Where("students.school_id = ?", schoolID).
		Order("permits.exit_time DESC").
		Offset(offset).
		Limit(filter.PageSize).
		Find(&permits).Error

	if err != nil {
		return nil, 0, err
	}

	return permits, total, nil
}

// UpdatePermit updates a permit record
// Requirements: 8.4 - THE System SHALL allow recording of return time
func (r *repository) UpdatePermit(ctx context.Context, permit *models.Permit) error {
	result := r.db.WithContext(ctx).
		Model(&models.Permit{}).
		Where("id = ?", permit.ID).
		Updates(map[string]interface{}{
			"student_id":          permit.StudentID,
			"reason":              permit.Reason,
			"exit_time":           permit.ExitTime,
			"return_time":         permit.ReturnTime,
			"responsible_teacher": permit.ResponsibleTeacher,
			"document_url":        permit.DocumentURL,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrPermitNotFound
	}
	return nil
}

// DeletePermit deletes a permit record
func (r *repository) DeletePermit(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.Permit{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrPermitNotFound
	}
	return nil
}

// ==================== Counseling Note Repository ====================

// CreateCounselingNote creates a new counseling note
// Requirements: 9.1 - WHEN a Guru_BK creates a counseling note
func (r *repository) CreateCounselingNote(ctx context.Context, note *models.CounselingNote) error {
	return r.db.WithContext(ctx).Create(note).Error
}

// FindCounselingNoteByID retrieves a counseling note by ID
func (r *repository) FindCounselingNoteByID(ctx context.Context, id uint) (*models.CounselingNote, error) {
	var note models.CounselingNote
	err := r.db.WithContext(ctx).
		Preload("Student").
		Preload("Student.Class").
		Preload("Creator").
		Where("id = ?", id).
		First(&note).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCounselingNoteNotFound
		}
		return nil, err
	}
	return &note, nil
}

// FindCounselingNotesByStudent retrieves all counseling notes for a student
// Requirements: 9.5 - THE System SHALL maintain counseling history per student with timestamps
func (r *repository) FindCounselingNotesByStudent(ctx context.Context, studentID uint) ([]models.CounselingNote, error) {
	var notes []models.CounselingNote
	err := r.db.WithContext(ctx).
		Preload("Student").
		Preload("Student.Class").
		Preload("Creator").
		Where("student_id = ?", studentID).
		Order("created_at DESC").
		Find(&notes).Error
	return notes, err
}

// FindCounselingNotes retrieves counseling notes with pagination and filtering
func (r *repository) FindCounselingNotes(ctx context.Context, schoolID uint, filter CounselingNoteFilter) ([]models.CounselingNote, int64, error) {
	var notes []models.CounselingNote
	var total int64

	query := r.db.WithContext(ctx).
		Model(&models.CounselingNote{}).
		Joins("JOIN students ON students.id = counseling_notes.student_id").
		Where("students.school_id = ?", schoolID)

	// Apply filters
	if filter.StudentID != nil {
		query = query.Where("counseling_notes.student_id = ?", *filter.StudentID)
	}
	if filter.ClassID != nil {
		query = query.Where("students.class_id = ?", *filter.ClassID)
	}
	if filter.StartDate != nil {
		startDate, err := time.Parse("2006-01-02", *filter.StartDate)
		if err == nil {
			query = query.Where("counseling_notes.created_at >= ?", startDate)
		}
	}
	if filter.EndDate != nil {
		endDate, err := time.Parse("2006-01-02", *filter.EndDate)
		if err == nil {
			query = query.Where("counseling_notes.created_at <= ?", endDate.Add(24*time.Hour))
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
		Joins("JOIN students ON students.id = counseling_notes.student_id").
		Where("students.school_id = ?", schoolID).
		Order("counseling_notes.created_at DESC").
		Offset(offset).
		Limit(filter.PageSize).
		Find(&notes).Error

	if err != nil {
		return nil, 0, err
	}

	return notes, total, nil
}

// UpdateCounselingNote updates a counseling note
func (r *repository) UpdateCounselingNote(ctx context.Context, note *models.CounselingNote) error {
	result := r.db.WithContext(ctx).
		Model(&models.CounselingNote{}).
		Where("id = ?", note.ID).
		Updates(map[string]interface{}{
			"student_id":     note.StudentID,
			"internal_note":  note.InternalNote,
			"parent_summary": note.ParentSummary,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrCounselingNoteNotFound
	}
	return nil
}

// DeleteCounselingNote deletes a counseling note
func (r *repository) DeleteCounselingNote(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.CounselingNote{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrCounselingNoteNotFound
	}
	return nil
}


// ==================== Student/User Lookup ====================

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

// ==================== Dashboard/Summary Operations ====================

// GetViolationCount returns total violation count for a school
func (r *repository) GetViolationCount(ctx context.Context, schoolID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Violation{}).
		Joins("JOIN students ON students.id = violations.student_id").
		Where("students.school_id = ?", schoolID).
		Count(&count).Error
	return count, err
}

// GetAchievementCount returns total achievement count for a school
func (r *repository) GetAchievementCount(ctx context.Context, schoolID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Achievement{}).
		Joins("JOIN students ON students.id = achievements.student_id").
		Where("students.school_id = ?", schoolID).
		Count(&count).Error
	return count, err
}

// GetPermitCount returns total permit count for a school
func (r *repository) GetPermitCount(ctx context.Context, schoolID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Permit{}).
		Joins("JOIN students ON students.id = permits.student_id").
		Where("students.school_id = ?", schoolID).
		Count(&count).Error
	return count, err
}

// GetActivePermitCount returns count of permits where student hasn't returned
func (r *repository) GetActivePermitCount(ctx context.Context, schoolID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Permit{}).
		Joins("JOIN students ON students.id = permits.student_id").
		Where("students.school_id = ? AND permits.return_time IS NULL", schoolID).
		Count(&count).Error
	return count, err
}

// GetCounselingCount returns total counseling note count for a school
func (r *repository) GetCounselingCount(ctx context.Context, schoolID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.CounselingNote{}).
		Joins("JOIN students ON students.id = counseling_notes.student_id").
		Where("students.school_id = ?", schoolID).
		Count(&count).Error
	return count, err
}

// GetStudentsNeedingAttention returns students with high violation counts
func (r *repository) GetStudentsNeedingAttention(ctx context.Context, schoolID uint, limit int) ([]StudentAttentionItem, error) {
	type result struct {
		StudentID      uint
		StudentName    string
		ClassName      string
		ViolationCount int
	}

	var results []result
	err := r.db.WithContext(ctx).
		Model(&models.Violation{}).
		Select("students.id as student_id, students.name as student_name, classes.name as class_name, COUNT(*) as violation_count").
		Joins("JOIN students ON students.id = violations.student_id").
		Joins("JOIN classes ON classes.id = students.class_id").
		Where("students.school_id = ?", schoolID).
		Group("students.id, students.name, classes.name").
		Having("COUNT(*) >= ?", 3). // Students with 3+ violations
		Order("violation_count DESC").
		Limit(limit).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	items := make([]StudentAttentionItem, len(results))
	for i, r := range results {
		items[i] = StudentAttentionItem{
			StudentID:      r.StudentID,
			StudentName:    r.StudentName,
			ClassName:      r.ClassName,
			ViolationCount: r.ViolationCount,
			Reason:         "Multiple violations recorded",
		}
	}

	return items, nil
}
