package auth

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/school-management/backend/internal/config"
)

var (
	ErrTokenExpired   = errors.New("token has expired")
	ErrTokenInvalid   = errors.New("token is invalid")
	ErrTokenMalformed = errors.New("token is malformed")
)

// JWTManager handles JWT token operations
type JWTManager struct {
	secretKey            []byte
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
	issuer               string
}

// jwtClaims represents the JWT claims structure
type jwtClaims struct {
	UserID   uint    `json:"user_id"`
	SchoolID *uint   `json:"school_id"`
	Role     string  `json:"role"`
	Username string  `json:"username"`
	Type     string  `json:"type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

// Custom UnmarshalJSON to handle school_id properly
func (c *jwtClaims) UnmarshalJSON(data []byte) error {
	// Use a temporary struct with interface{} for school_id
	type Alias jwtClaims
	aux := &struct {
		SchoolID interface{} `json:"school_id"`
		*Alias
	}{
		Alias: (*Alias)(c),
	}
	
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	
	// Handle school_id conversion
	if aux.SchoolID != nil {
		switch v := aux.SchoolID.(type) {
		case float64:
			schoolID := uint(v)
			c.SchoolID = &schoolID
		case int:
			schoolID := uint(v)
			c.SchoolID = &schoolID
		case int64:
			schoolID := uint(v)
			c.SchoolID = &schoolID
		case uint:
			c.SchoolID = &v
		}
	}
	
	return nil
}

// NewJWTManager creates a new JWT manager
func NewJWTManager(cfg config.JWTConfig) *JWTManager {
	return &JWTManager{
		secretKey:            []byte(cfg.SecretKey),
		accessTokenDuration:  time.Duration(cfg.AccessTokenDuration) * time.Minute,
		refreshTokenDuration: time.Duration(cfg.RefreshTokenDuration) * time.Hour,
		issuer:               cfg.Issuer,
	}
}

// GenerateTokenPair generates both access and refresh tokens
// Requirements: 12.1 - WHEN authentication succeeds, THE System SHALL return JWT tokens
func (m *JWTManager) GenerateTokenPair(claims TokenClaims) (*TokenPair, error) {
	accessToken, err := m.generateToken(claims, "access", m.accessTokenDuration)
	if err != nil {
		return nil, err
	}

	refreshToken, err := m.generateToken(claims, "refresh", m.refreshTokenDuration)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(m.accessTokenDuration.Seconds()),
	}, nil
}

// generateToken generates a single JWT token
func (m *JWTManager) generateToken(claims TokenClaims, tokenType string, duration time.Duration) (string, error) {
	now := time.Now()
	expiresAt := now.Add(duration)

	jwtClaims := jwtClaims{
		UserID:   claims.UserID,
		SchoolID: claims.SchoolID,
		Role:     claims.Role,
		Username: claims.Username,
		Type:     tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			Subject:   claims.Username,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	return token.SignedString(m.secretKey)
}

// ValidateToken validates a JWT token and returns the claims
// Requirements: 12.3 - JWT token validation
func (m *JWTManager) ValidateToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrTokenInvalid
		}
		return m.secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, ErrTokenMalformed
		}
		return nil, ErrTokenInvalid
	}

	claims, ok := token.Claims.(*jwtClaims)
	if !ok || !token.Valid {
		return nil, ErrTokenInvalid
	}

	return &TokenClaims{
		UserID:   claims.UserID,
		SchoolID: claims.SchoolID,
		Role:     claims.Role,
		Username: claims.Username,
		Type:     claims.Type,
	}, nil
}

// ValidateAccessToken validates an access token
func (m *JWTManager) ValidateAccessToken(tokenString string) (*TokenClaims, error) {
	claims, err := m.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.Type != "access" {
		return nil, ErrTokenInvalid
	}

	return claims, nil
}

// ValidateRefreshToken validates a refresh token
func (m *JWTManager) ValidateRefreshToken(tokenString string) (*TokenClaims, error) {
	claims, err := m.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.Type != "refresh" {
		return nil, ErrTokenInvalid
	}

	return claims, nil
}

// GetAccessTokenDuration returns the access token duration in seconds
func (m *JWTManager) GetAccessTokenDuration() int64 {
	return int64(m.accessTokenDuration.Seconds())
}
