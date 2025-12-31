package grade

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

// ==================== Grade DTOs ====================

// CreateGradeRequest represents the request to create a grade
// Requirements: 10.1 - Grade SHALL require title, score, and student_id
type CreateGradeRequest struct {
	StudentID   uint    `json:"student_id" validate:"required"`
	Title       string  `json:"title" validate:"required"`
	Score       float64 `json:"score" validate:"required,min=0,max=100"`
	Description string  `json:"description"`
}

// UpdateGradeRequest represents the request to update a grade
type UpdateGradeRequest struct {
	Title       string  `json:"title"`
	Score       float64 `json:"score" validate:"min=0,max=100"`
	Description string  `json:"description"`
}

// GradeResponse represents a grade in responses
type GradeResponse struct {
	ID          uint      `json:"id"`
	StudentID   uint      `json:"student_id"`
	StudentName string    `json:"student_name,omitempty"`
	StudentNIS  string    `json:"student_nis,omitempty"`
	StudentNISN string    `json:"student_nisn,omitempty"`
	ClassName   string    `json:"class_name,omitempty"`
	Title       string    `json:"title"`
	Score       float64   `json:"score"`
	Description string    `json:"description"`
	CreatedBy   uint      `json:"created_by"`
	CreatorName string    `json:"creator_name,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// GradeListResponse represents a paginated list of grades
type GradeListResponse struct {
	Grades     []GradeResponse `json:"grades"`
	Pagination PaginationMeta  `json:"pagination"`
}

// GradeFilter represents filter options for listing grades
type GradeFilter struct {
	StudentID *uint   `query:"student_id"`
	ClassID   *uint   `query:"class_id"`
	StartDate *string `query:"start_date"`
	EndDate   *string `query:"end_date"`
	Page      int     `query:"page"`
	PageSize  int     `query:"page_size"`
}

// StudentGradeSummary represents grade summary for a student
type StudentGradeSummary struct {
	StudentID    uint    `json:"student_id"`
	StudentName  string  `json:"student_name"`
	ClassName    string  `json:"class_name"`
	TotalGrades  int     `json:"total_grades"`
	AverageScore float64 `json:"average_score"`
}
