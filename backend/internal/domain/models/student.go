package models

import (
	"errors"
	"strings"
	"time"
)

// Student represents a student
type Student struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	SchoolID  uint      `gorm:"index;not null" json:"school_id"`
	ClassID   uint      `gorm:"index;not null" json:"class_id"`
	UserID    *uint     `gorm:"uniqueIndex" json:"user_id"` // For student login
	NIS       string    `gorm:"type:varchar(20);not null" json:"nis"`
	NISN      string    `gorm:"type:varchar(20);uniqueIndex;not null" json:"nisn"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	RFIDCode  string    `gorm:"type:varchar(50);index" json:"rfid_code"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	School  School   `gorm:"foreignKey:SchoolID" json:"school,omitempty"`
	Class   Class    `gorm:"foreignKey:ClassID" json:"class,omitempty"`
	User    *User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Parents []Parent `gorm:"many2many:student_parents" json:"parents,omitempty"`
}

// TableName specifies the table name for Student
func (Student) TableName() string {
	return "students"
}

// Validate validates the student data
// Requirements: 3.2 - Student registration SHALL require NIS, NISN, name, and class assignment
func (s *Student) Validate() error {
	if strings.TrimSpace(s.NIS) == "" {
		return errors.New("NIS is required")
	}
	if strings.TrimSpace(s.NISN) == "" {
		return errors.New("NISN is required")
	}
	if strings.TrimSpace(s.Name) == "" {
		return errors.New("name is required")
	}
	if s.ClassID == 0 {
		return errors.New("class_id is required")
	}
	if s.SchoolID == 0 {
		return errors.New("school_id is required")
	}
	return nil
}

// StudentParent represents the many-to-many relationship between students and parents
type StudentParent struct {
	StudentID uint `gorm:"primaryKey"`
	ParentID  uint `gorm:"primaryKey"`
}

// TableName specifies the table name for StudentParent
func (StudentParent) TableName() string {
	return "student_parents"
}
