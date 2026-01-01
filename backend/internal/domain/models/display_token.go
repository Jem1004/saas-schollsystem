package models

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"
	"time"
)

// DisplayToken represents a token for public display access
// Requirements: 5.1, 6.1 - Token-based access for public display screens
type DisplayToken struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	SchoolID       uint       `gorm:"index;not null" json:"school_id"`
	Token          string     `gorm:"type:varchar(64);uniqueIndex;not null" json:"token"`
	Name           string     `gorm:"type:varchar(100)" json:"name"`
	IsActive       bool       `gorm:"default:true" json:"is_active"`
	LastAccessedAt *time.Time `json:"last_accessed_at"`
	ExpiresAt      *time.Time `json:"expires_at"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`

	// Relations
	School School `gorm:"foreignKey:SchoolID;constraint:OnDelete:CASCADE" json:"school,omitempty"`
}

// TableName specifies the table name for DisplayToken
func (DisplayToken) TableName() string {
	return "display_tokens"
}

// Validate validates the display token data
func (t *DisplayToken) Validate() error {
	if t.SchoolID == 0 {
		return errors.New("school_id is required")
	}
	if strings.TrimSpace(t.Token) == "" {
		return errors.New("token is required")
	}
	if len(t.Token) != 64 {
		return errors.New("token must be 64 characters")
	}
	if t.Name != "" && len(t.Name) > 100 {
		return errors.New("name must be at most 100 characters")
	}
	return nil
}

// GenerateToken generates a cryptographically secure random token
// Requirements: 6.2 - Token generation SHALL use cryptographically secure random
func GenerateToken() (string, error) {
	bytes := make([]byte, 32) // 32 bytes = 64 hex characters
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// IsValid checks if the token is valid (active and not expired)
// Requirements: 5.10, 5.13 - Token validation for access control
func (t *DisplayToken) IsValid() bool {
	if !t.IsActive {
		return false
	}
	if t.ExpiresAt != nil && time.Now().After(*t.ExpiresAt) {
		return false
	}
	return true
}

// IsExpired checks if the token has expired
func (t *DisplayToken) IsExpired() bool {
	if t.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*t.ExpiresAt)
}

// Revoke revokes the display token
// Requirements: 5.10 - Revoking a token SHALL immediately invalidate access
func (t *DisplayToken) Revoke() {
	t.IsActive = false
}

// Activate activates the display token
func (t *DisplayToken) Activate() {
	t.IsActive = true
}

// UpdateLastAccessed updates the last accessed timestamp
// Requirements: 6.7 - Track last accessed timestamp for each display token
func (t *DisplayToken) UpdateLastAccessed() {
	now := time.Now()
	t.LastAccessedAt = &now
}

// SetExpiration sets the expiration date for the token
// Requirements: 6.8 - Allow setting optional expiration date
func (t *DisplayToken) SetExpiration(expiresAt time.Time) {
	t.ExpiresAt = &expiresAt
}

// ClearExpiration removes the expiration date
func (t *DisplayToken) ClearExpiration() {
	t.ExpiresAt = nil
}

// Regenerate generates a new token value
// Requirements: 6.5 - Regenerating a token SHALL create a new token and invalidate the old one
func (t *DisplayToken) Regenerate() error {
	newToken, err := GenerateToken()
	if err != nil {
		return err
	}
	t.Token = newToken
	t.IsActive = true
	t.LastAccessedAt = nil
	return nil
}
