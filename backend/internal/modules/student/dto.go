package student

import "time"

// ==================== Profile DTOs ====================

// StudentProfileResponse represents the student's profile
type StudentProfileResponse struct {
	ID        uint   `json:"id"`
	NIS       string `json:"nis"`
	NISN      string `json:"nisn"`
	Name      string `json:"name"`
	ClassName string `json:"class_name"`
	Grade     int    `json:"grade"`
	Year      string `json:"year"`
	IsActive  bool   `json:"is_active"`
}

// StudentSummaryResponse represents summary statistics for a student
type StudentSummaryResponse struct {
	AttendanceSummary AttendanceSummaryResponse `json:"attendance_summary"`
	GradeSummary      GradeSummaryResponse      `json:"grade_summary"`
	BKSummary         BKSummaryResponse         `json:"bk_summary"`
}

// ==================== Attendance DTOs ====================

// AttendanceResponse represents an attendance record
type AttendanceResponse struct {
	ID           uint       `json:"id"`
	Date         string     `json:"date"`
	CheckInTime  *time.Time `json:"check_in_time"`
	CheckOutTime *time.Time `json:"check_out_time"`
	Status       string     `json:"status"`
	Method       string     `json:"method"`
}

// AttendanceListResponse represents paginated attendance list
type AttendanceListResponse struct {
	Attendances []AttendanceResponse `json:"attendances"`
	Pagination  PaginationMeta       `json:"pagination"`
}

// AttendanceSummaryResponse represents attendance summary
type AttendanceSummaryResponse struct {
	TotalDays int `json:"total_days"`
	Present   int `json:"present"`
	Late      int `json:"late"`
	VeryLate  int `json:"very_late"`
	Absent    int `json:"absent"`
}

// ==================== Grade DTOs ====================

// GradeResponse represents a grade record
type GradeResponse struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Score       float64   `json:"score"`
	Description string    `json:"description"`
	TeacherName string    `json:"teacher_name"`
	CreatedAt   time.Time `json:"created_at"`
}

// GradeListResponse represents paginated grade list
type GradeListResponse struct {
	Grades     []GradeResponse `json:"grades"`
	Pagination PaginationMeta  `json:"pagination"`
}

// GradeSummaryResponse represents grade summary
type GradeSummaryResponse struct {
	TotalGrades  int     `json:"total_grades"`
	AverageScore float64 `json:"average_score"`
	HighestScore float64 `json:"highest_score"`
	LowestScore  float64 `json:"lowest_score"`
}

// ==================== BK DTOs ====================

// ViolationSummaryResponse represents a violation summary (limited info for students)
// Requirements: 16.5 - Students view summary only
type ViolationSummaryResponse struct {
	ID          uint      `json:"id"`
	Category    string    `json:"category"`
	Level       string    `json:"level"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// AchievementResponse represents an achievement record
type AchievementResponse struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Point       int       `json:"point"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// BKSummaryResponse represents BK summary for students
// Requirements: 16.5 - Students view achievements and violations (summary only)
type BKSummaryResponse struct {
	TotalPoints      int `json:"total_points"`
	ViolationCount   int `json:"violation_count"`
	AchievementCount int `json:"achievement_count"`
}

// BKInfoResponse represents BK information for students
type BKInfoResponse struct {
	TotalPoints        int                        `json:"total_points"`
	ViolationCount     int                        `json:"violation_count"`
	AchievementCount   int                        `json:"achievement_count"`
	RecentViolations   []ViolationSummaryResponse `json:"recent_violations"`
	RecentAchievements []AchievementResponse      `json:"recent_achievements"`
}

// ==================== Dashboard DTOs ====================

// DashboardResponse represents the student dashboard
type DashboardResponse struct {
	Profile            StudentProfileResponse    `json:"profile"`
	Summary            StudentSummaryResponse    `json:"summary"`
	RecentAttendance   []AttendanceResponse      `json:"recent_attendance"`
	RecentGrades       []GradeResponse           `json:"recent_grades"`
	RecentAchievements []AchievementResponse     `json:"recent_achievements"`
}

// ==================== Common DTOs ====================

// PaginationMeta represents pagination metadata
type PaginationMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// AttendanceFilter represents filter options for attendance
type AttendanceFilter struct {
	StartDate string `query:"start_date"`
	EndDate   string `query:"end_date"`
	Page      int    `query:"page"`
	PageSize  int    `query:"page_size"`
}

// DefaultAttendanceFilter returns default filter values
func DefaultAttendanceFilter() AttendanceFilter {
	return AttendanceFilter{
		Page:     1,
		PageSize: 20,
	}
}

// GradeFilter represents filter options for grades
type GradeFilter struct {
	Page     int `query:"page"`
	PageSize int `query:"page_size"`
}

// DefaultGradeFilter returns default filter values
func DefaultGradeFilter() GradeFilter {
	return GradeFilter{
		Page:     1,
		PageSize: 20,
	}
}
