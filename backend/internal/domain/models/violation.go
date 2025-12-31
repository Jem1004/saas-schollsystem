package models

import (
	"errors"
	"strings"
	"time"
)

// ViolationLevel represents the severity of a violation
type ViolationLevel string

const (
	ViolationLevelRingan ViolationLevel = "ringan"
	ViolationLevelSedang ViolationLevel = "sedang"
	ViolationLevelBerat  ViolationLevel = "berat"
)

// IsValid checks if the violation level is valid
func (v ViolationLevel) IsValid() bool {
	switch v {
	case ViolationLevelRingan, ViolationLevelSedang, ViolationLevelBerat:
		return true
	}
	return false
}

// Violation represents student violation record
type Violation struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	StudentID   uint           `gorm:"index;not null" json:"student_id"`
	Category    string         `gorm:"type:varchar(100);not null" json:"category"`
	Level       ViolationLevel `gorm:"type:varchar(20);not null" json:"level"`
	Description string         `gorm:"type:text;not null" json:"description"`
	CreatedBy   uint           `gorm:"not null" json:"created_by"`
	CreatedAt   time.Time      `json:"created_at"`

	// Relations
	Student Student `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	Creator User    `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

// TableName specifies the table name for Violation
func (Violation) TableName() string {
	return "violations"
}

// Validate validates the violation data
// Requirements: 6.1 - Violation SHALL require category, level, description, and student identifier
func (v *Violation) Validate() error {
	if v.StudentID == 0 {
		return errors.New("student_id is required")
	}
	if strings.TrimSpace(v.Category) == "" {
		return errors.New("category is required")
	}
	if !v.Level.IsValid() {
		return errors.New("level must be one of: ringan, sedang, berat")
	}
	if strings.TrimSpace(v.Description) == "" {
		return errors.New("description is required")
	}
	if v.CreatedBy == 0 {
		return errors.New("created_by is required")
	}
	return nil
}
