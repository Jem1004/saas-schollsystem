package models

import (
	"errors"
	"strings"
	"time"
)

// Valid timezone constants for Indonesia
const (
	TimezoneWIB  = "Asia/Jakarta"   // WIB (UTC+7) - Western Indonesia
	TimezoneWITA = "Asia/Makassar"  // WITA (UTC+8) - Central Indonesia
	TimezoneWIT  = "Asia/Jayapura"  // WIT (UTC+9) - Eastern Indonesia
)

// School represents a tenant in the multi-tenant system
type School struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	Address   string    `gorm:"type:text" json:"address"`
	Phone     string    `gorm:"type:varchar(20)" json:"phone"`
	Email     string    `gorm:"type:varchar(255)" json:"email"`
	Timezone  string    `gorm:"type:varchar(50);default:'Asia/Makassar'" json:"timezone"`
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
	// Validate timezone if provided
	if s.Timezone != "" && !IsValidTimezone(s.Timezone) {
		return errors.New("invalid timezone, must be one of: Asia/Jakarta (WIB), Asia/Makassar (WITA), Asia/Jayapura (WIT)")
	}
	return nil
}

// IsValidTimezone checks if the timezone is valid for Indonesia
func IsValidTimezone(tz string) bool {
	validTimezones := []string{TimezoneWIB, TimezoneWITA, TimezoneWIT}
	for _, valid := range validTimezones {
		if tz == valid {
			return true
		}
	}
	return false
}

// GetLocation returns the time.Location for this school's timezone
func (s *School) GetLocation() *time.Location {
	loc, err := time.LoadLocation(s.Timezone)
	if err != nil {
		// Fallback to WITA if timezone is invalid
		loc, _ = time.LoadLocation(TimezoneWITA)
	}
	return loc
}

// GetCurrentTime returns the current time in the school's timezone
func (s *School) GetCurrentTime() time.Time {
	return time.Now().In(s.GetLocation())
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
