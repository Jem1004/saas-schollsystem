package auth

import "time"

// LoginRequest represents the login request payload
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse represents the login response payload
type LoginResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresIn    int64        `json:"expires_in"` // seconds until access token expires
	TokenType    string       `json:"token_type"`
	User         UserResponse `json:"user"`
}

// RefreshTokenRequest represents the refresh token request payload
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// RefreshTokenResponse represents the refresh token response payload
type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// ChangePasswordRequest represents the change password request payload
// Requirements: 12.5 - System SHALL enforce password reset on first login
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

// UserResponse represents the user data in responses
type UserResponse struct {
	ID           uint       `json:"id"`
	SchoolID     *uint      `json:"school_id"`
	Role         string     `json:"role"`
	Username     string     `json:"username"`
	Email        string     `json:"email"`
	Name         string     `json:"name"`
	IsActive     bool       `json:"is_active"`
	MustResetPwd bool       `json:"must_reset_pwd"`
	LastLoginAt  *time.Time `json:"last_login_at"`
}

// TokenPair represents a pair of access and refresh tokens
type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
}

// TokenClaims represents the JWT token claims
type TokenClaims struct {
	UserID   uint   `json:"user_id"`
	SchoolID *uint  `json:"school_id"`
	Role     string `json:"role"`
	Username string `json:"username"`
	Type     string `json:"type"` // "access" or "refresh"
}
