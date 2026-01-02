package importmodule

import (
	"errors"
	"io"
	"mime/multipart"
	"strings"

	"github.com/xuri/excelize/v2"
)

// Parser errors
var (
	ErrInvalidFileFormat = errors.New("format file tidak valid, hanya menerima file .xlsx")
	ErrFileTooLarge      = errors.New("ukuran file melebihi batas maksimum 5MB")
	ErrEmptyFile         = errors.New("file tidak memiliki data")
	ErrInvalidHeader     = errors.New("header kolom tidak sesuai dengan template")
)

// MaxFileSize is the maximum allowed file size (5MB)
const MaxFileSize = 5 * 1024 * 1024

// ExcelParser defines the interface for parsing Excel files
// Requirements: 2.1, 2.3
type ExcelParser interface {
	ParseStudentFile(file multipart.File, fileSize int64) ([]StudentRow, error)
	ParseParentFile(file multipart.File, fileSize int64) ([]ParentRow, error)
}

// excelParser implements ExcelParser interface
type excelParser struct{}

// NewExcelParser creates a new Excel parser
func NewExcelParser() ExcelParser {
	return &excelParser{}
}

// ParseStudentFile parses an Excel file for student import
// Requirements: 2.3 - Read all rows starting from row 2 (after header)
func (p *excelParser) ParseStudentFile(file multipart.File, fileSize int64) ([]StudentRow, error) {
	// Validate file size
	// Requirements: 2.5 - File size must not exceed 5MB
	if fileSize > MaxFileSize {
		return nil, ErrFileTooLarge
	}

	// Read file content
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, ErrInvalidFileFormat
	}

	// Open Excel file
	f, err := excelize.OpenReader(strings.NewReader(string(content)))
	if err != nil {
		return nil, ErrInvalidFileFormat
	}
	defer f.Close()

	// Get the first sheet
	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		return nil, ErrEmptyFile
	}

	// Get all rows
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, ErrInvalidFileFormat
	}

	// Check if file has data (at least header + 1 data row)
	// Requirements: 2.4 - Return error if no data rows
	if len(rows) < 2 {
		return nil, ErrEmptyFile
	}

	// Validate header row
	if !p.validateStudentHeader(rows[0]) {
		return nil, ErrInvalidHeader
	}

	// Parse data rows (starting from row 2, index 1)
	var studentRows []StudentRow
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		
		// Skip completely empty rows
		if p.isEmptyRow(row) {
			continue
		}

		studentRow := StudentRow{
			RowNumber: i + 1, // Excel row number (1-indexed)
		}

		// Parse columns (handle missing columns gracefully)
		if len(row) > 0 {
			studentRow.NIS = strings.TrimSpace(row[0])
		}
		if len(row) > 1 {
			studentRow.NISN = strings.TrimSpace(row[1])
		}
		if len(row) > 2 {
			studentRow.Name = strings.TrimSpace(row[2])
		}
		if len(row) > 3 {
			studentRow.ClassName = strings.TrimSpace(row[3])
		}

		studentRows = append(studentRows, studentRow)
	}

	// Check if we have any data rows after parsing
	if len(studentRows) == 0 {
		return nil, ErrEmptyFile
	}

	return studentRows, nil
}

// ParseParentFile parses an Excel file for parent import
// Requirements: 2.3 - Read all rows starting from row 2 (after header)
func (p *excelParser) ParseParentFile(file multipart.File, fileSize int64) ([]ParentRow, error) {
	// Validate file size
	// Requirements: 2.5 - File size must not exceed 5MB
	if fileSize > MaxFileSize {
		return nil, ErrFileTooLarge
	}

	// Read file content
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, ErrInvalidFileFormat
	}

	// Open Excel file
	f, err := excelize.OpenReader(strings.NewReader(string(content)))
	if err != nil {
		return nil, ErrInvalidFileFormat
	}
	defer f.Close()

	// Get the first sheet
	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		return nil, ErrEmptyFile
	}

	// Get all rows
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, ErrInvalidFileFormat
	}

	// Check if file has data (at least header + 1 data row)
	// Requirements: 2.4 - Return error if no data rows
	if len(rows) < 2 {
		return nil, ErrEmptyFile
	}

	// Validate header row
	if !p.validateParentHeader(rows[0]) {
		return nil, ErrInvalidHeader
	}

	// Parse data rows (starting from row 2, index 1)
	var parentRows []ParentRow
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		
		// Skip completely empty rows
		if p.isEmptyRow(row) {
			continue
		}

		parentRow := ParentRow{
			RowNumber: i + 1, // Excel row number (1-indexed)
		}

		// Parse columns (handle missing columns gracefully)
		if len(row) > 0 {
			parentRow.Name = strings.TrimSpace(row[0])
		}
		if len(row) > 1 {
			parentRow.Phone = strings.TrimSpace(row[1])
		}
		if len(row) > 2 {
			parentRow.Email = strings.TrimSpace(row[2])
		}

		parentRows = append(parentRows, parentRow)
	}

	// Check if we have any data rows after parsing
	if len(parentRows) == 0 {
		return nil, ErrEmptyFile
	}

	return parentRows, nil
}

// validateStudentHeader validates the header row for student import
func (p *excelParser) validateStudentHeader(header []string) bool {
	if len(header) < 3 { // At least NIS, NISN, Nama required
		return false
	}

	// Check required columns (case-insensitive)
	expectedHeaders := []string{"nis", "nisn", "nama"}
	for i, expected := range expectedHeaders {
		if i >= len(header) {
			return false
		}
		if strings.ToLower(strings.TrimSpace(header[i])) != expected {
			return false
		}
	}

	return true
}

// validateParentHeader validates the header row for parent import
func (p *excelParser) validateParentHeader(header []string) bool {
	if len(header) < 2 { // At least Nama, No_HP required
		return false
	}

	// Check required columns (case-insensitive)
	// Allow variations: "No_HP", "No HP", "NoHP", "Phone"
	nameHeader := strings.ToLower(strings.TrimSpace(header[0]))
	phoneHeader := strings.ToLower(strings.TrimSpace(header[1]))

	if nameHeader != "nama" {
		return false
	}

	validPhoneHeaders := []string{"no_hp", "no hp", "nohp", "phone", "telepon", "hp"}
	phoneValid := false
	for _, valid := range validPhoneHeaders {
		if phoneHeader == valid {
			phoneValid = true
			break
		}
	}

	return phoneValid
}

// isEmptyRow checks if a row is completely empty
func (p *excelParser) isEmptyRow(row []string) bool {
	for _, cell := range row {
		if strings.TrimSpace(cell) != "" {
			return false
		}
	}
	return true
}
