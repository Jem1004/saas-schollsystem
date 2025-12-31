package models

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"
	"time"
)

// Device represents RFID device (ESP32)
type Device struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	SchoolID    uint       `gorm:"index;not null" json:"school_id"`
	DeviceCode  string     `gorm:"type:varchar(50);uniqueIndex;not null" json:"device_code"`
	APIKey      string     `gorm:"type:varchar(255);uniqueIndex;not null" json:"-"` // Hidden from JSON
	Description string     `gorm:"type:varchar(255)" json:"description"`
	IsActive    bool       `gorm:"default:true" json:"is_active"`
	LastSeenAt  *time.Time `json:"last_seen_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	// Relations
	School School `gorm:"foreignKey:SchoolID" json:"school,omitempty"`
}

// TableName specifies the table name for Device
func (Device) TableName() string {
	return "devices"
}

// Validate validates the device data
// Requirements: 2.1 - Device registration SHALL generate a unique API key
func (d *Device) Validate() error {
	if d.SchoolID == 0 {
		return errors.New("school_id is required")
	}
	if strings.TrimSpace(d.DeviceCode) == "" {
		return errors.New("device_code is required")
	}
	return nil
}

// GenerateAPIKey generates a new unique API key for the device
func (d *Device) GenerateAPIKey() error {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return err
	}
	d.APIKey = hex.EncodeToString(bytes)
	return nil
}

// UpdateLastSeen updates the last seen timestamp
func (d *Device) UpdateLastSeen() {
	now := time.Now()
	d.LastSeenAt = &now
}
