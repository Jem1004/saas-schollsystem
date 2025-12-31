package models

import (
	"errors"
	"strings"
	"time"
)

// CounselingNote represents counseling session notes
type CounselingNote struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	StudentID     uint      `gorm:"index;not null" json:"student_id"`
	InternalNote  string    `gorm:"type:text;not null" json:"-"` // Hidden from JSON by default
	ParentSummary string    `gorm:"type:text" json:"parent_summary"`
	CreatedBy     uint      `gorm:"not null" json:"created_by"`
	CreatedAt     time.Time `json:"created_at"`

	// Relations
	Student Student `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	Creator User    `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

// TableName specifies the table name for CounselingNote
func (CounselingNote) TableName() string {
	return "counseling_notes"
}

// Validate validates the counseling note data
// Requirements: 9.1 - Counseling note SHALL require internal_note and student_id
func (c *CounselingNote) Validate() error {
	if c.StudentID == 0 {
		return errors.New("student_id is required")
	}
	if strings.TrimSpace(c.InternalNote) == "" {
		return errors.New("internal_note is required")
	}
	if c.CreatedBy == 0 {
		return errors.New("created_by is required")
	}
	return nil
}

// CounselingNoteWithInternal is used when internal notes should be visible (for Guru BK)
type CounselingNoteWithInternal struct {
	CounselingNote
	InternalNote string `json:"internal_note"`
}

// ToWithInternal converts CounselingNote to CounselingNoteWithInternal
func (c *CounselingNote) ToWithInternal() CounselingNoteWithInternal {
	return CounselingNoteWithInternal{
		CounselingNote: *c,
		InternalNote:   c.InternalNote,
	}
}
