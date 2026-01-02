package importmodule

import (
	"bytes"
	"context"
	"mime/multipart"
	"regexp"
	"strings"

	"github.com/school-management/backend/internal/domain/models"
	"github.com/xuri/excelize/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Default password for imported parent accounts
const DefaultPassword = "password123"

// Service defines the interface for import business logic
// Requirements: 1.1, 1.2, 2.1
type Service interface {
	// Template generation
	GenerateStudentTemplate() ([]byte, error)
	GenerateParentTemplate() ([]byte, error)

	// Import operations
	ImportStudents(ctx context.Context, schoolID uint, file multipart.File, fileSize int64) (*ImportResult, error)
	ImportParents(ctx context.Context, schoolID uint, file multipart.File, fileSize int64) (*ImportResult, error)
}

// service implements the Service interface
type service struct {
	db           *gorm.DB
	parser       ExcelParser
	classMatcher ClassMatcher
}

// NewService creates a new import service
func NewService(db *gorm.DB) Service {
	return &service{
		db:           db,
		parser:       NewExcelParser(),
		classMatcher: NewClassMatcher(db),
	}
}

// GenerateStudentTemplate generates an Excel template for student import
// Requirements: 1.1, 1.3, 1.4
func (s *service) GenerateStudentTemplate() ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "Sheet1"

	// Set header row
	// Requirements: 1.1 - Columns: NIS, NISN, Nama, Kelas
	headers := StudentTemplateColumns
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
	}

	// Set header style (bold)
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#E0E0E0"}, Pattern: 1},
	})
	f.SetCellStyle(sheetName, "A1", "D1", headerStyle)

	// Add example data row
	// Requirements: 1.4 - Include example data row
	exampleData := []string{"12345", "1234567890", "Nama Siswa", "X-IPA-1"}
	for i, data := range exampleData {
		cell, _ := excelize.CoordinatesToCellName(i+1, 2)
		f.SetCellValue(sheetName, cell, data)
	}

	// Set column widths
	f.SetColWidth(sheetName, "A", "A", 15)
	f.SetColWidth(sheetName, "B", "B", 15)
	f.SetColWidth(sheetName, "C", "C", 30)
	f.SetColWidth(sheetName, "D", "D", 15)

	// Write to buffer
	buf := new(bytes.Buffer)
	if err := f.Write(buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// GenerateParentTemplate generates an Excel template for parent import
// Requirements: 1.2, 1.3, 1.4
func (s *service) GenerateParentTemplate() ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "Sheet1"

	// Set header row
	// Requirements: 1.2 - Columns: Nama, No_HP, Email
	headers := ParentTemplateColumns
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
	}

	// Set header style (bold)
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#E0E0E0"}, Pattern: 1},
	})
	f.SetCellStyle(sheetName, "A1", "C1", headerStyle)

	// Add example data row
	// Requirements: 1.4 - Include example data row
	exampleData := []string{"Nama Orang Tua", "081234567890", "orangtua@email.com"}
	for i, data := range exampleData {
		cell, _ := excelize.CoordinatesToCellName(i+1, 2)
		f.SetCellValue(sheetName, cell, data)
	}

	// Set column widths
	f.SetColWidth(sheetName, "A", "A", 30)
	f.SetColWidth(sheetName, "B", "B", 20)
	f.SetColWidth(sheetName, "C", "C", 30)

	// Write to buffer
	buf := new(bytes.Buffer)
	if err := f.Write(buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// ImportStudents imports students from an Excel file
// Requirements: 3.1-3.11, 5.1-5.6
func (s *service) ImportStudents(ctx context.Context, schoolID uint, file multipart.File, fileSize int64) (*ImportResult, error) {
	// Parse Excel file
	rows, err := s.parser.ParseStudentFile(file, fileSize)
	if err != nil {
		return nil, err
	}

	result := &ImportResult{
		TotalRows: len(rows),
		Errors:    []ImportError{},
		Warnings:  []ImportWarning{},
	}

	// Process within transaction
	// Requirements: 5.6 - Process import within database transaction
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, row := range rows {
			// Validate required fields
			// Requirements: 3.1, 3.2, 3.3
			if err := s.validateStudentRow(row, result); err != nil {
				continue // Skip invalid rows
			}

			// Check for duplicate NISN
			// Requirements: 3.4
			var existingByNISN models.Student
			if err := tx.Where("nisn = ?", row.NISN).First(&existingByNISN).Error; err == nil {
				result.Warnings = append(result.Warnings, ImportWarning{
					Row:     row.RowNumber,
					Field:   "NISN",
					Message: "NISN sudah terdaftar di sistem",
				})
				result.WarningCount++
				continue
			}

			// Check for duplicate NIS within school
			// Requirements: 3.5
			var existingByNIS models.Student
			if err := tx.Where("school_id = ? AND nis = ?", schoolID, row.NIS).First(&existingByNIS).Error; err == nil {
				result.Warnings = append(result.Warnings, ImportWarning{
					Row:     row.RowNumber,
					Field:   "NIS",
					Message: "NIS sudah terdaftar di sekolah ini",
				})
				result.WarningCount++
				continue
			}

			// Match class if provided
			// Requirements: 3.6, 3.7, 3.8
			var classID *uint
			if row.ClassName != "" {
				class, err := s.classMatcher.FindClassByName(ctx, schoolID, row.ClassName)
				if err != nil {
					return err
				}
				if class != nil {
					classID = &class.ID
				} else {
					// Class not found - add warning
					// Requirements: 3.7
					result.Warnings = append(result.Warnings, ImportWarning{
						Row:     row.RowNumber,
						Field:   "Kelas",
						Message: "Kelas '" + row.ClassName + "' tidak ditemukan, siswa akan diimport tanpa kelas",
					})
				}
			}

			// Create student
			// Requirements: 3.9, 3.10, 3.11
			student := &models.Student{
				SchoolID: schoolID,
				ClassID:  classID,
				NIS:      row.NIS,
				NISN:     row.NISN,
				Name:     row.Name,
				IsActive: classID != nil, // IsActive true only if ClassID is set
			}

			if err := tx.Create(student).Error; err != nil {
				result.Errors = append(result.Errors, ImportError{
					Row:     row.RowNumber,
					Field:   "",
					Message: "Gagal menyimpan data: " + err.Error(),
				})
				result.FailedCount++
				continue
			}

			result.SuccessCount++
			if classID == nil {
				result.StudentsWithoutClass++
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Update failed count based on errors
	result.FailedCount = len(result.Errors)

	return result, nil
}

// ImportParents imports parents from an Excel file
// Requirements: 4.1-4.7, 5.1-5.6
func (s *service) ImportParents(ctx context.Context, schoolID uint, file multipart.File, fileSize int64) (*ImportResult, error) {
	// Parse Excel file
	rows, err := s.parser.ParseParentFile(file, fileSize)
	if err != nil {
		return nil, err
	}

	result := &ImportResult{
		TotalRows: len(rows),
		Errors:    []ImportError{},
		Warnings:  []ImportWarning{},
	}

	// Process within transaction
	// Requirements: 5.6 - Process import within database transaction
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, row := range rows {
			// Validate required fields
			// Requirements: 4.1, 4.2
			if err := s.validateParentRow(row, result); err != nil {
				continue // Skip invalid rows
			}

			// Check for duplicate phone (username)
			// Requirements: 4.3
			var existingUser models.User
			if err := tx.Where("username = ?", row.Phone).First(&existingUser).Error; err == nil {
				result.Warnings = append(result.Warnings, ImportWarning{
					Row:     row.RowNumber,
					Field:   "No_HP",
					Message: "Nomor HP sudah terdaftar sebagai username",
				})
				result.WarningCount++
				continue
			}

			// Generate password hash
			// Requirements: 4.5
			passwordHash, err := bcrypt.GenerateFromPassword([]byte(DefaultPassword), bcrypt.DefaultCost)
			if err != nil {
				result.Errors = append(result.Errors, ImportError{
					Row:     row.RowNumber,
					Field:   "",
					Message: "Gagal membuat password: " + err.Error(),
				})
				result.FailedCount++
				continue
			}

			// Create user account
			// Requirements: 4.4
			user := &models.User{
				SchoolID:     &schoolID,
				Role:         models.RoleParent,
				Username:     row.Phone,
				PasswordHash: string(passwordHash),
				Email:        row.Email,
				Name:         row.Name,
				IsActive:     true,
				MustResetPwd: true,
			}

			if err := tx.Create(user).Error; err != nil {
				result.Errors = append(result.Errors, ImportError{
					Row:     row.RowNumber,
					Field:   "",
					Message: "Gagal membuat akun user: " + err.Error(),
				})
				result.FailedCount++
				continue
			}

			// Create parent record
			// Requirements: 4.6, 4.7
			parent := &models.Parent{
				SchoolID: schoolID,
				UserID:   user.ID,
				Name:     row.Name,
				Phone:    row.Phone,
			}

			if err := tx.Create(parent).Error; err != nil {
				result.Errors = append(result.Errors, ImportError{
					Row:     row.RowNumber,
					Field:   "",
					Message: "Gagal menyimpan data orang tua: " + err.Error(),
				})
				result.FailedCount++
				continue
			}

			result.SuccessCount++
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Update failed count based on errors
	result.FailedCount = len(result.Errors)

	return result, nil
}

// validateStudentRow validates a student row and adds errors to result
// Requirements: 3.1, 3.2, 3.3
func (s *service) validateStudentRow(row StudentRow, result *ImportResult) error {
	hasError := false

	// Requirements: 3.1 - NIS is required
	if strings.TrimSpace(row.NIS) == "" {
		result.Errors = append(result.Errors, ImportError{
			Row:     row.RowNumber,
			Field:   "NIS",
			Message: "NIS wajib diisi",
		})
		hasError = true
	}

	// Requirements: 3.2 - NISN is required
	if strings.TrimSpace(row.NISN) == "" {
		result.Errors = append(result.Errors, ImportError{
			Row:     row.RowNumber,
			Field:   "NISN",
			Message: "NISN wajib diisi",
		})
		hasError = true
	}

	// Requirements: 3.3 - Name is required
	if strings.TrimSpace(row.Name) == "" {
		result.Errors = append(result.Errors, ImportError{
			Row:     row.RowNumber,
			Field:   "Nama",
			Message: "Nama wajib diisi",
		})
		hasError = true
	}

	if hasError {
		result.FailedCount++
		return ErrInvalidHeader // Return any error to indicate validation failed
	}

	return nil
}

// validateParentRow validates a parent row and adds errors to result
// Requirements: 4.1, 4.2
func (s *service) validateParentRow(row ParentRow, result *ImportResult) error {
	hasError := false

	// Requirements: 4.1 - Name is required
	if strings.TrimSpace(row.Name) == "" {
		result.Errors = append(result.Errors, ImportError{
			Row:     row.RowNumber,
			Field:   "Nama",
			Message: "Nama wajib diisi",
		})
		hasError = true
	}

	// Phone is required (used as username)
	if strings.TrimSpace(row.Phone) == "" {
		result.Errors = append(result.Errors, ImportError{
			Row:     row.RowNumber,
			Field:   "No_HP",
			Message: "Nomor HP wajib diisi",
		})
		hasError = true
	}

	// Requirements: 4.2 - Email is required and must be valid format
	email := strings.TrimSpace(row.Email)
	if email == "" {
		result.Errors = append(result.Errors, ImportError{
			Row:     row.RowNumber,
			Field:   "Email",
			Message: "Email wajib diisi",
		})
		hasError = true
	} else if !isValidEmail(email) {
		result.Errors = append(result.Errors, ImportError{
			Row:     row.RowNumber,
			Field:   "Email",
			Message: "Format email tidak valid",
		})
		hasError = true
	}

	if hasError {
		result.FailedCount++
		return ErrInvalidHeader // Return any error to indicate validation failed
	}

	return nil
}

// isValidEmail validates email format
func isValidEmail(email string) bool {
	// Simple email regex pattern
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(strings.TrimSpace(email))
}
