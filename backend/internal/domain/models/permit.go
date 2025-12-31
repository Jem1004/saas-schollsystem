package models

import (
	"errors"
	"strings"
	"time"
)

// Permit represents school exit permit
type Permit struct {
	ID                 uint       `gorm:"primaryKey" json:"id"`
	StudentID          uint       `gorm:"index;not null" json:"student_id"`
	Reason             string     `gorm:"type:text;not null" json:"reason"`
	ExitTime           time.Time  `gorm:"not null" json:"exit_time"`
	ReturnTime         *time.Time `json:"return_time"`
	ResponsibleTeacher uint       `gorm:"not null" json:"responsible_teacher"`
	DocumentURL        string     `gorm:"type:varchar(500)" json:"document_url"`
	CreatedBy          uint       `gorm:"not null" json:"created_by"`
	CreatedAt          time.Time  `json:"created_at"`

	// Relations
	Student Student `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	Teacher User    `gorm:"foreignKey:ResponsibleTeacher" json:"teacher,omitempty"`
	Creator User    `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

// TableName specifies the table name for Permit
func (Permit) TableName() string {
	return "permits"
}

// Validate validates the permit data
// Requirements: 8.1 - Permit SHALL require reason, exit time, and responsible teacher
func (p *Permit) Validate() error {
	if p.StudentID == 0 {
		return errors.New("student_id is required")
	}
	if strings.TrimSpace(p.Reason) == "" {
		return errors.New("reason is required")
	}
	if p.ExitTime.IsZero() {
		return errors.New("exit_time is required")
	}
	if p.ResponsibleTeacher == 0 {
		return errors.New("responsible_teacher is required")
	}
	if p.CreatedBy == 0 {
		return errors.New("created_by is required")
	}
	return nil
}

// HasReturned checks if the student has returned
func (p *Permit) HasReturned() bool {
	return p.ReturnTime != nil
}
