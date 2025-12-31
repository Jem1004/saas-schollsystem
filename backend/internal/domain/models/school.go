package models

import (
	"errors"
	"strings"
	"time"
)

// School represents a tenant in the multi-tenant system
type School struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	Address   string    `gorm:"type:text" json:"address"`
	Phone     string    `gorm:"type:varchar(20)" json:"phone"`
	Email     string    `gorm:"type:varchar(255)" json:"email"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	Classes  []Class   `gorm:"foreignKey:SchoolID" json:"classes,omitempty"`
	Students []Student `gorm:"foreignKey:SchoolID" json:"students,omitempty"`
	Users    []User    `gorm:"foreignKey:SchoolID" json:"users,omitempty"`
	Devices  []Device  `gorm:"foreignKey:SchoolID" json:"devices,omitempty"`
}

// TableName specifies the table name for School
func (School) TableName() string {
	return "schools"
}

// Validate validates the school data
// Requirements: 1.1 - Tenant creation SHALL generate a unique school_id
func (s *School) Validate() error {
	if strings.TrimSpace(s.Name) == "" {
		return errors.New("name is required")
	}
	return nil
}

// Deactivate deactivates the school
// Requirements: 1.3 - Deactivating a tenant SHALL prevent all users from accessing the system
func (s *School) Deactivate() {
	s.IsActive = false
}

// Activate activates the school
func (s *School) Activate() {
	s.IsActive = true
}
