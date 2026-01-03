package homeroom

import (
	"time"
)

// ==================== Pagination ====================

// PaginationMeta represents pagination metadata
type PaginationMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// ==================== Homeroom Note DTOs ====================

// CreateNoteRequest represents the request to create a homeroom note
// Requirements: 11.1 - Homeroom note SHALL require content and student_id
type CreateNoteRequest struct {
	StudentID uint   `json:"studentId" validate:"required"`
	Content   string `json:"content" validate:"required"`
}

// UpdateNoteRequest represents the request to update a homeroom note
type UpdateNoteRequest struct {
	Content string `json:"content"`
}

// NoteResponse represents a homeroom note in responses
type NoteResponse struct {
	ID          uint      `json:"id"`
	StudentID   uint      `json:"studentId"`
	StudentName string    `json:"studentName,omitempty"`
	StudentNIS  string    `json:"studentNis,omitempty"`
	StudentNISN string    `json:"studentNisn,omitempty"`
	ClassName   string    `json:"className,omitempty"`
	TeacherID   uint      `json:"teacherId"`
	TeacherName string    `json:"teacherName,omitempty"`
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// NoteListResponse represents a paginated list of homeroom notes
type NoteListResponse struct {
	Notes      []NoteResponse `json:"notes"`
	Pagination PaginationMeta `json:"pagination"`
}

// NoteFilter represents filter options for listing homeroom notes
type NoteFilter struct {
	StudentID *uint   `query:"student_id"`
	ClassID   *uint   `query:"class_id"`
	TeacherID *uint   `query:"teacher_id"`
	StartDate *string `query:"start_date"`
	EndDate   *string `query:"end_date"`
	Page      int     `query:"page"`
	PageSize  int     `query:"page_size"`
}

// StudentNoteSummary represents note summary for a student
type StudentNoteSummary struct {
	StudentID   uint   `json:"student_id"`
	StudentName string `json:"student_name"`
	ClassName   string `json:"class_name"`
	TotalNotes  int    `json:"total_notes"`
	LastNoteAt  *time.Time `json:"last_note_at,omitempty"`
}

// ==================== Wali Kelas Dashboard DTOs ====================

// TodayAttendanceStats represents today's attendance statistics for a class
type TodayAttendanceStats struct {
	Present int64 `json:"present"`
	Absent  int64 `json:"absent"`
	Late    int64 `json:"late"`
	Excused int64 `json:"excused"`
}

// HomeroomStatsResponse represents the dashboard statistics for wali kelas
type HomeroomStatsResponse struct {
	ClassID         uint                 `json:"classId"`
	ClassName       string               `json:"className"`
	TotalStudents   int64                `json:"totalStudents"`
	TodayAttendance TodayAttendanceStats `json:"todayAttendance"`
	RecentGrades    []GradeResponse      `json:"recentGrades"`
	RecentNotes     []NoteResponse       `json:"recentNotes"`
}

// GradeResponse represents a grade in responses (for dashboard)
type GradeResponse struct {
	ID          uint      `json:"id"`
	StudentID   uint      `json:"studentId"`
	StudentName string    `json:"studentName"`
	StudentNIS  string    `json:"studentNis"`
	Title       string    `json:"title"`
	Score       float64   `json:"score"`
	Description string    `json:"description,omitempty"`
	CreatedBy   uint      `json:"createdBy"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// GradeListResponse represents a paginated list of grades
type GradeListResponse struct {
	Data       []GradeResponse `json:"data"`
	Pagination PaginationMeta  `json:"pagination"`
}

// CreateGradeRequest represents the request to create a grade
type CreateGradeRequest struct {
	StudentID   uint    `json:"studentId" validate:"required"`
	Title       string  `json:"title" validate:"required"`
	Score       float64 `json:"score" validate:"required,min=0,max=100"`
	Description string  `json:"description,omitempty"`
}

// UpdateGradeRequest represents the request to update a grade
type UpdateGradeRequest struct {
	Title       string  `json:"title,omitempty"`
	Score       float64 `json:"score,omitempty"`
	Description string  `json:"description,omitempty"`
}

// BatchGradeRequest represents the request to create multiple grades
type BatchGradeRequest struct {
	Title       string       `json:"title" validate:"required"`
	Description string       `json:"description,omitempty"`
	Grades      []GradeEntry `json:"grades" validate:"required,min=1"`
}

// GradeEntry represents a single grade entry in batch request
type GradeEntry struct {
	StudentID uint    `json:"studentId" validate:"required"`
	Score     float64 `json:"score" validate:"required,min=0,max=100"`
}

// ClassInfoResponse represents class information for wali kelas
type ClassInfoResponse struct {
	ID                  uint   `json:"id"`
	Name                string `json:"name"`
	Grade               int    `json:"grade"`
	Year                string `json:"year"`
	TotalStudents       int64  `json:"totalStudents"`
	HomeroomTeacherID   uint   `json:"homeroomTeacherId"`
	HomeroomTeacherName string `json:"homeroomTeacherName"`
}

// ClassStudentResponse represents a student in the class
type ClassStudentResponse struct {
	ID       uint   `json:"id"`
	NIS      string `json:"nis"`
	NISN     string `json:"nisn"`
	Name     string `json:"name"`
	IsActive bool   `json:"isActive"`
}

// ClassStudentListResponse represents a paginated list of students
type ClassStudentListResponse struct {
	Data       []ClassStudentResponse `json:"data"`
	Pagination PaginationMeta         `json:"pagination"`
}

// StudentAttendanceResponse represents attendance data for a student
type StudentAttendanceResponse struct {
	ID           uint    `json:"id"`
	StudentID    uint    `json:"studentId"`
	StudentName  string  `json:"studentName"`
	StudentNIS   string  `json:"studentNis"`
	Date         string  `json:"date"`
	CheckInTime  *string `json:"checkInTime"`
	CheckOutTime *string `json:"checkOutTime"`
	Status       string  `json:"status"`
	Method       string  `json:"method"`
	CreatedAt    string  `json:"createdAt"`
	UpdatedAt    string  `json:"updatedAt"`
}

// ClassAttendanceListResponse represents attendance data for a class
type ClassAttendanceListResponse struct {
	Data       []StudentAttendanceResponse `json:"data"`
	Pagination PaginationMeta              `json:"pagination"`
}

// ManualAttendanceRequest represents the request to record manual attendance
type ManualAttendanceRequest struct {
	StudentID    uint   `json:"studentId" validate:"required"`
	ScheduleID   uint   `json:"scheduleId" validate:"required"`
	Date         string `json:"date" validate:"required"`
	Status       string `json:"status" validate:"required"`
	CheckInTime  string `json:"checkInTime,omitempty"`
	CheckOutTime string `json:"checkOutTime,omitempty"`
}

// UpdateAttendanceRequest represents the request to update attendance
type UpdateAttendanceRequest struct {
	Status       string `json:"status,omitempty"`
	CheckInTime  string `json:"checkInTime,omitempty"`
	CheckOutTime string `json:"checkOutTime,omitempty"`
}

// ScheduleResponse represents an attendance schedule for wali kelas
type ScheduleResponse struct {
	ID                uint   `json:"id"`
	Name              string `json:"name"`
	StartTime         string `json:"startTime"`
	EndTime           string `json:"endTime"`
	LateThreshold     int    `json:"lateThreshold"`
	VeryLateThreshold *int   `json:"veryLateThreshold,omitempty"`
	IsDefault         bool   `json:"isDefault"`
}
