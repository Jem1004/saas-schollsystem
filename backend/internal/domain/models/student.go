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
	ClassID   *uint     `gorm:"index" json:"class_id"` // Nullable for import workflow
	UserID    *uint     `gorm:"uniqueIndex" json:"user_id"` // For student login
	NIS       string    `gorm:"type:varchar(20);not null" json:"nis"`
	NISN      string    `gorm:"type:varchar(20);uniqueIndex;not null" json:"nisn"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	RFIDCode  string    `gorm:"type:varchar(50);index" json:"rfid_code"`
	IsActive  bool      `gorm:"default:false" json:"is_active"` // Default false, true only when ClassID is set
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	School  School   `gorm:"foreignKey:SchoolID" json:"school,omitempty"`
	Class   *Class   `gorm:"foreignKey:ClassID" json:"class,omitempty"` // Pointer for nullable ClassID
	User    *User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Parents []Parent `gorm:"many2many:student_parents" json:"parents,omitempty"`
}

// TableName specifies the table name for Student
func (Student) TableName() string {
	return "students"
}

// Validate validates the student data
// Requirements: 8.1 - ClassID is now optional for import workflow
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
	if s.SchoolID == 0 {
		return errors.New("school_id is required")
	}
	// ClassID validation removed - now optional for import workflow
	return nil
}

// CanBeActive checks if student can be set to active
// Requirements: 8.2, 8.3 - Student can only be active if ClassID is set
func (s *Student) CanBeActive() bool {
	return s.ClassID != nil && *s.ClassID > 0
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
