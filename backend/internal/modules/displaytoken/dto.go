package displaytoken

import (
	"time"
)

// ==================== Request DTOs ====================

// CreateDisplayTokenRequest represents the request to create a new display token
// Requirements: 6.1, 6.2 - Token creation with name and optional expiration
type CreateDisplayTokenRequest struct {
	Name      string     `json:"name" validate:"required,max=100"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"` // optional expiration date
}

// UpdateDisplayTokenRequest represents the request to update a display token
type UpdateDisplayTokenRequest struct {
	Name      *string    `json:"name,omitempty" validate:"omitempty,max=100"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	IsActive  *bool      `json:"is_active,omitempty"`
}

// ==================== Response DTOs ====================

// DisplayTokenResponse represents a display token in responses (without the actual token value)
// Requirements: 6.1 - List tokens with name, status, and last accessed time
type DisplayTokenResponse struct {
	ID             uint       `json:"id"`
	SchoolID       uint       `json:"school_id"`
	Name           string     `json:"name"`
	IsActive       bool       `json:"is_active"`
	LastAccessedAt *time.Time `json:"last_accessed_at,omitempty"`
	ExpiresAt      *time.Time `json:"expires_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DisplayURL     string     `json:"display_url,omitempty"` // URL for the public display
}

// DisplayTokenWithSecretResponse represents a display token with the actual token value
// Requirements: 6.3 - Show the full token only once (similar to API key)
type DisplayTokenWithSecretResponse struct {
	ID             uint       `json:"id"`
	SchoolID       uint       `json:"school_id"`
	Token          string     `json:"token"` // The actual token value - shown only once
	Name           string     `json:"name"`
	IsActive       bool       `json:"is_active"`
	ExpiresAt      *time.Time `json:"expires_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	DisplayURL     string     `json:"display_url"` // Full URL for the public display
}

// DisplayTokenListResponse represents a list of display tokens
type DisplayTokenListResponse struct {
	Tokens []DisplayTokenResponse `json:"tokens"`
	Total  int                    `json:"total"`
}

// DisplayTokenValidation represents the result of token validation
// Requirements: 5.10, 5.13 - Token validation for access control
type DisplayTokenValidation struct {
	Valid    bool   `json:"valid"`
	SchoolID uint   `json:"school_id,omitempty"`
	TokenID  uint   `json:"token_id,omitempty"`
	Name     string `json:"name,omitempty"`
	Error    string `json:"error,omitempty"`
}
