package models

import (
	"time"
)

// ClassCounselor represents the many-to-many relationship between classes and BK teachers (counselors)
// This allows multiple BK teachers to be assigned to different classes
type ClassCounselor struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ClassID     uint      `gorm:"index;not null" json:"class_id"`
	CounselorID uint      `gorm:"index;not null" json:"counselor_id"` // User ID with role guru_bk
	SchoolID    uint      `gorm:"index;not null" json:"school_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relations
	Class     Class `gorm:"foreignKey:ClassID" json:"class,omitempty"`
	Counselor User  `gorm:"foreignKey:CounselorID" json:"counselor,omitempty"`
}

// TableName specifies the table name for ClassCounselor
func (ClassCounselor) TableName() string {
	return "class_counselors"
}
