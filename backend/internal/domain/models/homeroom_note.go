package models

import (
	"errors"
	"strings"
	"time"
)

// HomeroomNote represents homeroom teacher notes
type HomeroomNote struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	StudentID uint      `gorm:"index;not null" json:"student_id"`
	TeacherID uint      `gorm:"not null" json:"teacher_id"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	Student Student `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	Teacher User    `gorm:"foreignKey:TeacherID" json:"teacher,omitempty"`
}

// TableName specifies the table name for HomeroomNote
func (HomeroomNote) TableName() string {
	return "homeroom_notes"
}

// Validate validates the homeroom note data
// Requirements: 11.1 - Homeroom note SHALL require content and student_id
func (h *HomeroomNote) Validate() error {
	if h.StudentID == 0 {
		return errors.New("student_id is required")
	}
	if strings.TrimSpace(h.Content) == "" {
		return errors.New("content is required")
	}
	if h.TeacherID == 0 {
		return errors.New("teacher_id is required")
	}
	return nil
}
