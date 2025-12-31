package models

import (
	"errors"
	"strings"
	"time"
)

// Achievement represents student achievement record
type Achievement struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	StudentID   uint      `gorm:"index;not null" json:"student_id"`
	Title       string    `gorm:"type:varchar(255);not null" json:"title"`
	Point       int       `gorm:"not null" json:"point"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedBy   uint      `gorm:"not null" json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`

	// Relations
	Student Student `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	Creator User    `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

// TableName specifies the table name for Achievement
func (Achievement) TableName() string {
	return "achievements"
}

// Validate validates the achievement data
// Requirements: 7.1 - Achievement SHALL require title, point value, and description
func (a *Achievement) Validate() error {
	if a.StudentID == 0 {
		return errors.New("student_id is required")
	}
	if strings.TrimSpace(a.Title) == "" {
		return errors.New("title is required")
	}
	if a.Point <= 0 {
		return errors.New("point must be greater than 0")
	}
	if a.CreatedBy == 0 {
		return errors.New("created_by is required")
	}
	return nil
}
