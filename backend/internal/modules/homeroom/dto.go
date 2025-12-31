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
	StudentID uint   `json:"student_id" validate:"required"`
	Content   string `json:"content" validate:"required"`
}

// UpdateNoteRequest represents the request to update a homeroom note
type UpdateNoteRequest struct {
	Content string `json:"content"`
}

// NoteResponse represents a homeroom note in responses
type NoteResponse struct {
	ID          uint      `json:"id"`
	StudentID   uint      `json:"student_id"`
	StudentName string    `json:"student_name,omitempty"`
	StudentNIS  string    `json:"student_nis,omitempty"`
	StudentNISN string    `json:"student_nisn,omitempty"`
	ClassName   string    `json:"class_name,omitempty"`
	TeacherID   uint      `json:"teacher_id"`
	TeacherName string    `json:"teacher_name,omitempty"`
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
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
