package tenant

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/school-management/backend/internal/domain/models"
)

var (
	ErrSchoolNotFound    = errors.New("sekolah tidak ditemukan")
	ErrDuplicateSchool   = errors.New("sekolah dengan nama ini sudah terdaftar")
	ErrDuplicateUsername = errors.New("username sudah digunakan")
	ErrSchoolHasData     = errors.New("sekolah memiliki data terkait")
)

// Repository defines the interface for tenant data operations
type Repository interface {
	Create(ctx context.Context, school *models.School) error
	CreateWithAdmin(ctx context.Context, school *models.School, admin *models.User) error
	FindAll(ctx context.Context, filter SchoolFilter) ([]models.School, int64, error)
	FindByID(ctx context.Context, id uint) (*models.School, error)
	FindByName(ctx context.Context, name string) (*models.School, error)
	Update(ctx context.Context, school *models.School) error
	Deactivate(ctx context.Context, id uint) error
	Activate(ctx context.Context, id uint) error
	Delete(ctx context.Context, id uint) error
	GetStats(ctx context.Context, schoolID uint) (*SchoolStats, error)
	GetAdminUsers(ctx context.Context, schoolID uint) ([]models.User, error)
	UsernameExists(ctx context.Context, username string) (bool, error)
}

// repository implements the Repository interface
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new tenant repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Create creates a new school (tenant)
// Requirements: 1.1 - WHEN a Super_Admin creates a new tenant, THE System SHALL generate a unique school_id
func (r *repository) Create(ctx context.Context, school *models.School) error {
	result := r.db.WithContext(ctx).Create(school)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// CreateWithAdmin creates a new school with an admin user in a transaction
func (r *repository) CreateWithAdmin(ctx context.Context, school *models.School, admin *models.User) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create school first
		if err := tx.Create(school).Error; err != nil {
			return err
		}

		// Set school ID on admin user
		admin.SchoolID = &school.ID

		// Create admin user
		if err := tx.Create(admin).Error; err != nil {
			return err
		}

		return nil
	})
}

// UsernameExists checks if a username already exists
func (r *repository) UsernameExists(ctx context.Context, username string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("username = ?", username).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetAdminUsers retrieves admin users for a school
func (r *repository) GetAdminUsers(ctx context.Context, schoolID uint) ([]models.User, error) {
	var users []models.User
	err := r.db.WithContext(ctx).
		Where("school_id = ? AND role = ?", schoolID, models.RoleAdminSekolah).
		Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// Delete deletes a school and all associated data (cascade delete)
func (r *repository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Check if school exists
		var school models.School
		if err := tx.Where("id = ?", id).First(&school).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrSchoolNotFound
			}
			return err
		}

		// Delete all related data in order (respecting foreign key constraints)

		// 1. Delete notifications for users in this school
		if err := tx.Exec("DELETE FROM notifications WHERE user_id IN (SELECT id FROM users WHERE school_id = ?)", id).Error; err != nil {
			return err
		}

		// 2. Delete FCM tokens for users in this school
		if err := tx.Exec("DELETE FROM fcm_tokens WHERE user_id IN (SELECT id FROM users WHERE school_id = ?)", id).Error; err != nil {
			return err
		}

		// 3. Delete homeroom notes for students in this school
		if err := tx.Exec("DELETE FROM homeroom_notes WHERE student_id IN (SELECT id FROM students WHERE school_id = ?)", id).Error; err != nil {
			return err
		}

		// 4. Delete grades for students in this school
		if err := tx.Exec("DELETE FROM grades WHERE student_id IN (SELECT id FROM students WHERE school_id = ?)", id).Error; err != nil {
			return err
		}

		// 5. Delete BK records (violations, achievements, permits, counseling notes)
		if err := tx.Exec("DELETE FROM violations WHERE student_id IN (SELECT id FROM students WHERE school_id = ?)", id).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM achievements WHERE student_id IN (SELECT id FROM students WHERE school_id = ?)", id).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM permits WHERE student_id IN (SELECT id FROM students WHERE school_id = ?)", id).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM counseling_notes WHERE student_id IN (SELECT id FROM students WHERE school_id = ?)", id).Error; err != nil {
			return err
		}

		// 6. Delete attendance records for students in this school
		if err := tx.Exec("DELETE FROM attendances WHERE student_id IN (SELECT id FROM students WHERE school_id = ?)", id).Error; err != nil {
			return err
		}

		// 7. Delete student-parent relationships
		if err := tx.Exec("DELETE FROM student_parents WHERE student_id IN (SELECT id FROM students WHERE school_id = ?)", id).Error; err != nil {
			return err
		}

		// 8. Delete parents for this school
		if err := tx.Where("school_id = ?", id).Delete(&models.Parent{}).Error; err != nil {
			return err
		}

		// 9. Delete students
		if err := tx.Where("school_id = ?", id).Delete(&models.Student{}).Error; err != nil {
			return err
		}

		// 10. Delete classes
		if err := tx.Where("school_id = ?", id).Delete(&models.Class{}).Error; err != nil {
			return err
		}

		// 11. Delete devices
		if err := tx.Where("school_id = ?", id).Delete(&models.Device{}).Error; err != nil {
			return err
		}

		// 12. Delete school settings
		if err := tx.Where("school_id = ?", id).Delete(&models.SchoolSettings{}).Error; err != nil {
			return err
		}

		// 13. Delete users
		if err := tx.Where("school_id = ?", id).Delete(&models.User{}).Error; err != nil {
			return err
		}

		// 14. Finally delete the school
		if err := tx.Delete(&school).Error; err != nil {
			return err
		}

		return nil
	})
}

// FindAll retrieves all schools with pagination and filtering
// Requirements: 1.2 - WHEN a Super_Admin views the tenant list, THE System SHALL display all registered schools
func (r *repository) FindAll(ctx context.Context, filter SchoolFilter) ([]models.School, int64, error) {
	var schools []models.School
	var total int64

	query := r.db.WithContext(ctx).Model(&models.School{})

	// Apply filters
	if filter.Name != "" {
		query = query.Where("name ILIKE ?", "%"+filter.Name+"%")
	}
	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}

	// Count total records
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
	err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(filter.PageSize).
		Find(&schools).Error

	if err != nil {
		return nil, 0, err
	}

	return schools, total, nil
}

// FindByID retrieves a school by ID
func (r *repository) FindByID(ctx context.Context, id uint) (*models.School, error) {
	var school models.School
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&school).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSchoolNotFound
		}
		return nil, err
	}

	return &school, nil
}

// FindByName retrieves a school by name
func (r *repository) FindByName(ctx context.Context, name string) (*models.School, error) {
	var school models.School
	err := r.db.WithContext(ctx).
		Where("name = ?", name).
		First(&school).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSchoolNotFound
		}
		return nil, err
	}

	return &school, nil
}

// Update updates a school
func (r *repository) Update(ctx context.Context, school *models.School) error {
	result := r.db.WithContext(ctx).Save(school)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrSchoolNotFound
	}
	return nil
}

// Deactivate deactivates a school
// Requirements: 1.3 - WHEN a Super_Admin deactivates a tenant, THE System SHALL prevent all users from accessing
func (r *repository) Deactivate(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).
		Model(&models.School{}).
		Where("id = ?", id).
		Update("is_active", false)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrSchoolNotFound
	}
	return nil
}

// Activate activates a school
func (r *repository) Activate(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).
		Model(&models.School{}).
		Where("id = ?", id).
		Update("is_active", true)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrSchoolNotFound
	}
	return nil
}

// GetStats retrieves statistics for a school
func (r *repository) GetStats(ctx context.Context, schoolID uint) (*SchoolStats, error) {
	stats := &SchoolStats{}

	// Count classes
	if err := r.db.WithContext(ctx).
		Model(&models.Class{}).
		Where("school_id = ?", schoolID).
		Count(&stats.TotalClasses).Error; err != nil {
		return nil, err
	}

	// Count students
	if err := r.db.WithContext(ctx).
		Model(&models.Student{}).
		Where("school_id = ?", schoolID).
		Count(&stats.TotalStudents).Error; err != nil {
		return nil, err
	}

	// Count users
	if err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("school_id = ?", schoolID).
		Count(&stats.TotalUsers).Error; err != nil {
		return nil, err
	}

	// Count devices
	if err := r.db.WithContext(ctx).
		Model(&models.Device{}).
		Where("school_id = ?", schoolID).
		Count(&stats.TotalDevices).Error; err != nil {
		return nil, err
	}

	return stats, nil
}
