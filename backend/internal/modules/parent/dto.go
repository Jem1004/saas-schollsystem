package parent

import "time"

// ==================== Child DTOs ====================

// ChildResponse represents a linked child's basic information
type ChildResponse struct {
	ID        uint   `json:"id"`
	NIS       string `json:"nis"`
	NISN      string `json:"nisn"`
	Name      string `json:"name"`
	ClassName string `json:"class_name"`
	Grade     int    `json:"grade"`
	IsActive  bool   `json:"is_active"`
}

// ChildListResponse represents the list of linked children
type ChildListResponse struct {
	Children []ChildResponse `json:"children"`
}

// ==================== Attendance DTOs ====================

// ChildAttendanceResponse represents attendance data for a child
type ChildAttendanceResponse struct {
	ID           uint       `json:"id"`
	Date         string     `json:"date"`
	CheckInTime  *time.Time `json:"check_in_time"`
	CheckOutTime *time.Time `json:"check_out_time"`
	Status       string     `json:"status"`
	Method       string     `json:"method"`
}

// ChildAttendanceListResponse represents paginated attendance list
type ChildAttendanceListResponse struct {
	Attendances []ChildAttendanceResponse `json:"attendances"`
	Pagination  PaginationMeta            `json:"pagination"`
}

// AttendanceSummaryResponse represents attendance summary for a child
type AttendanceSummaryResponse struct {
	StudentID   uint   `json:"student_id"`
	StudentName string `json:"student_name"`
	TotalDays   int    `json:"total_days"`
	Present     int    `json:"present"`
	Late        int    `json:"late"`
	VeryLate    int    `json:"very_late"`
	Absent      int    `json:"absent"`
}

// ==================== Grade DTOs ====================

// ChildGradeResponse represents a grade for a child
type ChildGradeResponse struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Score       float64   `json:"score"`
	Description string    `json:"description"`
	TeacherName string    `json:"teacher_name"`
	CreatedAt   time.Time `json:"created_at"`
}

// ChildGradeListResponse represents paginated grade list
type ChildGradeListResponse struct {
	Grades     []ChildGradeResponse `json:"grades"`
	Pagination PaginationMeta       `json:"pagination"`
}

// GradeSummaryResponse represents grade summary for a child
type GradeSummaryResponse struct {
	StudentID    uint    `json:"student_id"`
	StudentName  string  `json:"student_name"`
	TotalGrades  int     `json:"total_grades"`
	AverageScore float64 `json:"average_score"`
	HighestScore float64 `json:"highest_score"`
	LowestScore  float64 `json:"lowest_score"`
}

// ==================== Homeroom Note DTOs ====================

// ChildNoteResponse represents a homeroom note for a child
type ChildNoteResponse struct {
	ID          uint      `json:"id"`
	Content     string    `json:"content"`
	TeacherName string    `json:"teacher_name"`
	CreatedAt   time.Time `json:"created_at"`
}

// ChildNoteListResponse represents paginated note list
type ChildNoteListResponse struct {
	Notes      []ChildNoteResponse `json:"notes"`
	Pagination PaginationMeta      `json:"pagination"`
}

// ==================== BK DTOs ====================

// ChildViolationResponse represents a violation for a child
type ChildViolationResponse struct {
	ID          uint      `json:"id"`
	Category    string    `json:"category"`
	Level       string    `json:"level"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// ChildAchievementResponse represents an achievement for a child
type ChildAchievementResponse struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Point       int       `json:"point"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// ChildPermitResponse represents a permit for a child
type ChildPermitResponse struct {
	ID          uint       `json:"id"`
	Reason      string     `json:"reason"`
	ExitTime    time.Time  `json:"exit_time"`
	ReturnTime  *time.Time `json:"return_time"`
	TeacherName string     `json:"teacher_name"`
	HasReturned bool       `json:"has_returned"`
	CreatedAt   time.Time  `json:"created_at"`
}

// ChildCounselingResponse represents a counseling note summary for a child (parent view)
// Requirements: 9.2, 9.4 - Only parent_summary visible to parents
type ChildCounselingResponse struct {
	ID            uint      `json:"id"`
	ParentSummary string    `json:"parent_summary"`
	CreatedAt     time.Time `json:"created_at"`
}

// ChildBKInfoResponse represents BK information for a child
// Requirements: 14.4 - Parent views violations, achievements, and permits
type ChildBKInfoResponse struct {
	TotalPoints        int                        `json:"total_points"`
	ViolationCount     int                        `json:"violation_count"`
	AchievementCount   int                        `json:"achievement_count"`
	PermitCount        int                        `json:"permit_count"`
	CounselingCount    int                        `json:"counseling_count"`
	RecentViolations   []ChildViolationResponse   `json:"recent_violations"`
	RecentAchievements []ChildAchievementResponse `json:"recent_achievements"`
	RecentPermits      []ChildPermitResponse      `json:"recent_permits"`
	RecentCounseling   []ChildCounselingResponse  `json:"recent_counseling"`
}

// ==================== Dashboard DTOs ====================

// ChildDashboardResponse represents dashboard data for a child
type ChildDashboardResponse struct {
	Student           ChildResponse             `json:"student"`
	AttendanceSummary AttendanceSummaryResponse `json:"attendance_summary"`
	GradeSummary      GradeSummaryResponse      `json:"grade_summary"`
	BKSummary         ChildBKSummaryResponse    `json:"bk_summary"`
	RecentGrades      []ChildGradeResponse      `json:"recent_grades"`
	RecentNotes       []ChildNoteResponse       `json:"recent_notes"`
}

// ChildBKSummaryResponse represents BK summary for dashboard
type ChildBKSummaryResponse struct {
	TotalPoints      int `json:"total_points"`
	ViolationCount   int `json:"violation_count"`
	AchievementCount int `json:"achievement_count"`
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

// NoteFilter represents filter options for notes
type NoteFilter struct {
	Page     int `query:"page"`
	PageSize int `query:"page_size"`
}

// DefaultNoteFilter returns default filter values
func DefaultNoteFilter() NoteFilter {
	return NoteFilter{
		Page:     1,
		PageSize: 20,
	}
}
