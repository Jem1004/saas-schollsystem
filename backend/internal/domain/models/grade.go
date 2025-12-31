package models

import (
	"errors"
	"strings"
	"time"
)

// Grade represents student grade entry
type Grade struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	StudentID   uint      `gorm:"index;not null" json:"student_id"`
	Title       string    `gorm:"type:varchar(255);not null" json:"title"`
	Score       float64   `gorm:"not null" json:"score"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedBy   uint      `gorm:"not null" json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relations
	Student Student `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	Creator User    `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

// TableName specifies the table name for Grade
func (Grade) TableName() string {
	return "grades"
}

// Validate validates the grade data
// Requirements: 10.1 - Grade SHALL require title, score, and student_id
func (g *Grade) Validate() error {
	if g.StudentID == 0 {
		return errors.New("student_id is required")
	}
	if strings.TrimSpace(g.Title) == "" {
		return errors.New("title is required")
	}
	if g.Score < 0 || g.Score > 100 {
		return errors.New("score must be between 0 and 100")
	}
	if g.CreatedBy == 0 {
		return errors.New("created_by is required")
	}
	return nil
}
