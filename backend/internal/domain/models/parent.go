package models

import (
	"errors"
	"strings"
	"time"
)

// Parent represents a parent/guardian
type Parent struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	SchoolID  uint      `gorm:"index;not null" json:"school_id"`
	UserID    uint      `gorm:"uniqueIndex;not null" json:"user_id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	Phone     string    `gorm:"type:varchar(20)" json:"phone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	User     User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Students []Student `gorm:"many2many:student_parents" json:"students,omitempty"`
}

// TableName specifies the table name for Parent
func (Parent) TableName() string {
	return "parents"
}

// Validate validates the parent data
// Requirements: 3.3 - Parent registration SHALL link the parent to one or more students
func (p *Parent) Validate() error {
	if p.SchoolID == 0 {
		return errors.New("school_id is required")
	}
	if p.UserID == 0 {
		return errors.New("user_id is required")
	}
	if strings.TrimSpace(p.Name) == "" {
		return errors.New("name is required")
	}
	return nil
}

// HasLinkedStudents checks if the parent has any linked students
func (p *Parent) HasLinkedStudents() bool {
	return len(p.Students) > 0
}
