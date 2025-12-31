package models

import (
	"errors"
	"strings"
	"time"
)

// Class represents a school class
type Class struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	SchoolID          uint      `gorm:"index;not null" json:"school_id"`
	Name              string    `gorm:"type:varchar(50);not null" json:"name"`
	Grade             int       `gorm:"not null" json:"grade"` // e.g., 7, 8, 9 for SMP
	Year              string    `gorm:"type:varchar(10);not null" json:"year"` // e.g., "2024/2025"
	HomeroomTeacherID *uint     `gorm:"index" json:"homeroom_teacher_id"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`

	// Relations
	School          School    `gorm:"foreignKey:SchoolID" json:"school,omitempty"`
	Students        []Student `gorm:"foreignKey:ClassID" json:"students,omitempty"`
	HomeroomTeacher *User     `gorm:"foreignKey:HomeroomTeacherID" json:"homeroom_teacher,omitempty"`
}

// TableName specifies the table name for Class
func (Class) TableName() string {
	return "classes"
}

// Validate validates the class data
// Requirements: 3.1 - Class creation SHALL associate it with the school tenant
func (c *Class) Validate() error {
	if c.SchoolID == 0 {
		return errors.New("school_id is required")
	}
	if strings.TrimSpace(c.Name) == "" {
		return errors.New("name is required")
	}
	if c.Grade <= 0 {
		return errors.New("grade must be greater than 0")
	}
	if strings.TrimSpace(c.Year) == "" {
		return errors.New("year is required")
	}
	return nil
}

// AssignHomeroomTeacher assigns a homeroom teacher to the class
// Requirements: 4.3 - Wali_Kelas role assignment SHALL require class assignment
func (c *Class) AssignHomeroomTeacher(teacherID uint) {
	c.HomeroomTeacherID = &teacherID
}

// RemoveHomeroomTeacher removes the homeroom teacher from the class
func (c *Class) RemoveHomeroomTeacher() {
	c.HomeroomTeacherID = nil
}

// HasHomeroomTeacher checks if the class has a homeroom teacher assigned
func (c *Class) HasHomeroomTeacher() bool {
	return c.HomeroomTeacherID != nil
}
