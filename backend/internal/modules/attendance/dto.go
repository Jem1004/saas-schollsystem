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
	ScheduleID   *uint                   `json:"schedule_id,omitempty"`
	ScheduleName string                  `json:"schedule_name,omitempty"` // Requirements: 3.10 - Show which schedule the attendance belongs to
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
// Requirements: 5.1, 5.2 - Include sick and excused counts in the response
type AttendanceSummaryResponse struct {
	Date       string `json:"date"`
	TotalCount int    `json:"total_count"`
	Present    int    `json:"present"`
	Late       int    `json:"late"`
	VeryLate   int    `json:"very_late"`
	Absent     int    `json:"absent"`
	Sick       int    `json:"sick"`
	Excused    int    `json:"excused"`
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
// Requirements: 5.3, 5.4 - Include sick and excused counts per class
type ClassSummaryItem struct {
	ClassID       uint   `json:"class_id"`
	ClassName     string `json:"class_name"`
	TotalStudents int    `json:"total_students"`
	Present       int    `json:"present"`
	Late          int    `json:"late"`
	Absent        int    `json:"absent"`
	Sick          int    `json:"sick"`
	Excused       int    `json:"excused"`
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

// ==================== Export DTOs ====================

// ExportFilter represents filter options for exporting attendance data
// Requirements: 1.2, 1.3 - Allow filtering by date range and class
type ExportFilter struct {
	StartDate string `query:"start_date" validate:"required"` // Format: YYYY-MM-DD
	EndDate   string `query:"end_date" validate:"required"`   // Format: YYYY-MM-DD
	ClassID   *uint  `query:"class_id"`
}

// MonthlyRecapFilter represents filter options for monthly recap
// Requirements: 2.3, 2.4 - Allow filtering by month, year, and class
type MonthlyRecapFilter struct {
	Month   int   `query:"month" validate:"required,min=1,max=12"` // 1-12
	Year    int   `query:"year" validate:"required,min=2000"`      // e.g., 2024
	ClassID *uint `query:"class_id"`
}

// MonthlyRecapResponse represents the monthly recap data
// Requirements: 2.1 - Display summary per student including total days present, late, very late, and absent
type MonthlyRecapResponse struct {
	Month          int                       `json:"month"`
	Year           int                       `json:"year"`
	TotalDays      int                       `json:"total_days"`       // Total school days in the month
	ClassID        *uint                     `json:"class_id,omitempty"`
	ClassName      string                    `json:"class_name,omitempty"`
	StudentRecaps  []StudentRecapSummary     `json:"student_recaps"`
}

// StudentRecapSummary represents attendance summary for a single student
// Requirements: 2.1, 2.2 - Summary per student with attendance percentage
// Requirements: 6.1, 6.2 - Include total_sick and total_excused in each student's recap summary
type StudentRecapSummary struct {
	StudentID         uint    `json:"student_id"`
	StudentNIS        string  `json:"student_nis"`
	StudentNISN       string  `json:"student_nisn"`
	StudentName       string  `json:"student_name"`
	ClassName         string  `json:"class_name"`
	TotalPresent      int     `json:"total_present"`
	TotalLate         int     `json:"total_late"`
	TotalVeryLate     int     `json:"total_very_late"`
	TotalAbsent       int     `json:"total_absent"`
	TotalSick         int     `json:"total_sick"`
	TotalExcused      int     `json:"total_excused"`
	AttendancePercent float64 `json:"attendance_percent"` // (present / total_days) * 100
}

// ExportAttendanceRecord represents a single attendance record for export
// Requirements: 1.4, 1.5 - Include student info and attendance details
type ExportAttendanceRecord struct {
	StudentNIS   string `json:"student_nis"`
	StudentNISN  string `json:"student_nisn"`
	StudentName  string `json:"student_name"`
	ClassName    string `json:"class_name"`
	Date         string `json:"date"`
	CheckInTime  string `json:"check_in_time"`
	CheckOutTime string `json:"check_out_time"`
	Status       string `json:"status"`
	ScheduleName string `json:"schedule_name,omitempty"`
}

// DailyAttendanceDetail represents daily attendance detail for monthly recap export
type DailyAttendanceDetail struct {
	Date         string `json:"date"`
	CheckInTime  string `json:"check_in_time"`
	CheckOutTime string `json:"check_out_time"`
	Status       string `json:"status"`
}
