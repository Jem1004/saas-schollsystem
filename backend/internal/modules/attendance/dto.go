package attendance

import (
	"time"

	"github.com/school-management/backend/internal/domain/models"
)

// ==================== Request DTOs ====================

// RFIDAttendanceRequest represents the request from ESP32 device
// Requirements: 5.1 - WHEN a student taps RFID card, THE ESP32 SHALL send student identifier and timestamp
type RFIDAttendanceRequest struct {
	APIKey    string    `json:"api_key" validate:"required"`
	RFIDCode  string    `json:"rfid_code" validate:"required"`
	Timestamp time.Time `json:"timestamp" validate:"required"`
}

// ManualAttendanceRequest represents manual attendance entry
// Requirements: 5.5 - IF RFID system fails, THEN THE System SHALL allow manual attendance entry
type ManualAttendanceRequest struct {
	StudentID    uint       `json:"student_id" validate:"required"`
	Date         string     `json:"date" validate:"required"` // Format: YYYY-MM-DD
	CheckInTime  *string    `json:"check_in_time"`            // Format: HH:MM
	CheckOutTime *string    `json:"check_out_time"`           // Format: HH:MM
}

// BulkManualAttendanceRequest represents bulk manual attendance entry
type BulkManualAttendanceRequest struct {
	Date        string                       `json:"date" validate:"required"` // Format: YYYY-MM-DD
	Attendances []BulkManualAttendanceItem   `json:"attendances" validate:"required,min=1"`
}

// BulkManualAttendanceItem represents a single item in bulk attendance
type BulkManualAttendanceItem struct {
	StudentID    uint    `json:"student_id" validate:"required"`
	CheckInTime  *string `json:"check_in_time"`  // Format: HH:MM
	CheckOutTime *string `json:"check_out_time"` // Format: HH:MM
}

// ==================== Response DTOs ====================

// AttendanceResponse represents attendance data in responses
type AttendanceResponse struct {
	ID           uint                    `json:"id"`
	StudentID    uint                    `json:"student_id"`
	StudentName  string                  `json:"student_name,omitempty"`
	StudentNIS   string                  `json:"student_nis,omitempty"`
	StudentNISN  string                  `json:"student_nisn,omitempty"`
	ClassName    string                  `json:"class_name,omitempty"`
	Date         string                  `json:"date"` // Format: YYYY-MM-DD
	CheckInTime  *string                 `json:"check_in_time,omitempty"`
	CheckOutTime *string                 `json:"check_out_time,omitempty"`
	Status       models.AttendanceStatus `json:"status"`
	Method       models.AttendanceMethod `json:"method"`
	CreatedAt    time.Time               `json:"created_at"`
	UpdatedAt    time.Time               `json:"updated_at"`
}

// AttendanceListResponse represents a paginated list of attendance records
type AttendanceListResponse struct {
	Attendances []AttendanceResponse `json:"attendances"`
	Pagination  PaginationMeta       `json:"pagination"`
}

// PaginationMeta represents pagination metadata
type PaginationMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// RFIDAttendanceResponse represents the response for RFID attendance
type RFIDAttendanceResponse struct {
	Success     bool                    `json:"success"`
	StudentID   uint                    `json:"student_id"`
	StudentName string                  `json:"student_name"`
	Type        string                  `json:"type"` // "check_in" or "check_out"
	Status      models.AttendanceStatus `json:"status,omitempty"`
	Time        time.Time               `json:"time"`
	Message     string                  `json:"message"`
}

// AttendanceSummaryResponse represents attendance summary for a class
// Requirements: 5.4 - WHEN an Admin_Sekolah views attendance dashboard, THE System SHALL display summary statistics
type AttendanceSummaryResponse struct {
	Date       string `json:"date"`
	TotalCount int    `json:"total_count"`
	Present    int    `json:"present"`
	Late       int    `json:"late"`
	VeryLate   int    `json:"very_late"`
	Absent     int    `json:"absent"`
}

// ClassAttendanceSummaryResponse represents attendance summary for a class with student details
type ClassAttendanceSummaryResponse struct {
	ClassID     uint                 `json:"class_id"`
	ClassName   string               `json:"class_name"`
	Date        string               `json:"date"`
	Summary     AttendanceSummaryResponse `json:"summary"`
	Attendances []AttendanceResponse `json:"attendances"`
}

// SchoolAttendanceSummaryResponse represents school-wide attendance summary
type SchoolAttendanceSummaryResponse struct {
	SchoolID    uint                        `json:"school_id"`
	SchoolName  string                      `json:"school_name,omitempty"`
	Date        string                      `json:"date"`
	Summary     AttendanceSummaryResponse   `json:"summary"`
	ByClass     []ClassSummaryItem          `json:"by_class,omitempty"`
}

// ClassSummaryItem represents attendance summary for a single class (for list view)
type ClassSummaryItem struct {
	ClassID       uint   `json:"class_id"`
	ClassName     string `json:"class_name"`
	TotalStudents int    `json:"total_students"`
	Present       int    `json:"present"`
	Late          int    `json:"late"`
	Absent        int    `json:"absent"`
}

// ==================== Filter DTOs ====================

// AttendanceFilter represents filter options for listing attendance
type AttendanceFilter struct {
	StudentID *uint   `query:"student_id"`
	ClassID   *uint   `query:"class_id"`
	StartDate *string `query:"start_date"` // Format: YYYY-MM-DD
	EndDate   *string `query:"end_date"`   // Format: YYYY-MM-DD
	Status    *string `query:"status"`
	Method    *string `query:"method"`
	Page      int     `query:"page"`
	PageSize  int     `query:"page_size"`
}

// DefaultAttendanceFilter returns default filter values
func DefaultAttendanceFilter() AttendanceFilter {
	return AttendanceFilter{
		Page:     1,
		PageSize: 20,
	}
}
