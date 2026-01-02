package importmodule

// ImportResult represents the result of an import operation
// Requirements: 5.1, 5.2, 5.3, 5.4, 5.5
type ImportResult struct {
	TotalRows            int             `json:"total_rows"`
	SuccessCount         int             `json:"success_count"`
	FailedCount          int             `json:"failed_count"`
	WarningCount         int             `json:"warning_count"`
	Errors               []ImportError   `json:"errors,omitempty"`
	Warnings             []ImportWarning `json:"warnings,omitempty"`
	StudentsWithoutClass int             `json:"students_without_class,omitempty"`
}

// ImportError represents an error during import
type ImportError struct {
	Row     int    `json:"row"`
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ImportWarning represents a warning during import
type ImportWarning struct {
	Row     int    `json:"row"`
	Field   string `json:"field"`
	Message string `json:"message"`
}

// StudentRow represents a parsed student row from Excel
// Requirements: 1.1 - Template columns: NIS, NISN, Nama, Kelas
type StudentRow struct {
	RowNumber int
	NIS       string
	NISN      string
	Name      string
	ClassName string // Optional - for class matching
}

// ParentRow represents a parsed parent row from Excel
// Requirements: 1.2 - Template columns: Nama, No_HP, Email
type ParentRow struct {
	RowNumber int
	Name      string
	Phone     string
	Email     string
}

// BulkAssignClassRequest represents request to assign class to multiple students
// Requirements: 6.2, 6.3
type BulkAssignClassRequest struct {
	StudentIDs []uint `json:"student_ids" validate:"required,min=1"`
	ClassID    uint   `json:"class_id" validate:"required"`
}

// StudentTemplateColumns defines the expected columns for student import
var StudentTemplateColumns = []string{"NIS", "NISN", "Nama", "Kelas"}

// ParentTemplateColumns defines the expected columns for parent import
var ParentTemplateColumns = []string{"Nama", "No_HP", "Email"}
