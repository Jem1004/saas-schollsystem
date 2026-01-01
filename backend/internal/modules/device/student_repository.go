package device

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/school-management/backend/internal/domain/models"
)

var (
	ErrStudentNotFound = errors.New("siswa tidak ditemukan")
)

// studentRepository implements StudentRepository interface for pairing service
type studentRepository struct {
	db *gorm.DB
}

// NewStudentRepository creates a new student repository for device module
func NewStudentRepository(db *gorm.DB) StudentRepository {
	return &studentRepository{db: db}
}

// FindByID retrieves a student by ID
func (r *studentRepository) FindByID(ctx context.Context, id uint) (*models.Student, error) {
	var student models.Student
	err := r.db.WithContext(ctx).
		Preload("Class").
		Where("id = ?", id).
		First(&student).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrStudentNotFound
		}
		return nil, err
	}

	return &student, nil
}

// FindByRFIDCode retrieves a student by RFID code
func (r *studentRepository) FindByRFIDCode(ctx context.Context, rfidCode string) (*models.Student, error) {
	var student models.Student
	err := r.db.WithContext(ctx).
		Where("rf_id_code = ? AND rf_id_code != ''", rfidCode).
		First(&student).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrStudentNotFound
		}
		return nil, err
	}

	return &student, nil
}

// UpdateRFIDCode updates a student's RFID code
func (r *studentRepository) UpdateRFIDCode(ctx context.Context, studentID uint, rfidCode string) error {
	result := r.db.WithContext(ctx).
		Model(&models.Student{}).
		Where("id = ?", studentID).
		Update("rf_id_code", rfidCode)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrStudentNotFound
	}
	return nil
}

// ClearRFIDCode clears a student's RFID code
func (r *studentRepository) ClearRFIDCode(ctx context.Context, studentID uint) error {
	result := r.db.WithContext(ctx).
		Model(&models.Student{}).
		Where("id = ?", studentID).
		Update("rf_id_code", "")

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrStudentNotFound
	}
	return nil
}
