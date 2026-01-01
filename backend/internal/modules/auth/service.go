package auth

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/school-management/backend/internal/domain/models"
)

var (
	ErrAccountInactive    = errors.New("akun tidak aktif")
	ErrSchoolInactive     = errors.New("sekolah tidak aktif")
	ErrPasswordMismatch   = errors.New("password lama salah")
	ErrSamePassword       = errors.New("password baru harus berbeda dari password lama")
	ErrPasswordTooShort   = errors.New("password minimal 8 karakter")
)

// Service defines the interface for auth business logic
type Service interface {
	Authenticate(ctx context.Context, username, password string) (*LoginResponse, error)
	RefreshAccessToken(ctx context.Context, refreshToken string) (*RefreshTokenResponse, error)
	ChangePassword(ctx context.Context, userID uint, oldPassword, newPassword string) error
	GetUserByID(ctx context.Context, userID uint) (*models.User, error)
}

// service implements the Service interface
type service struct {
	repo       Repository
	jwtManager *JWTManager
}

// NewService creates a new auth service
func NewService(repo Repository, jwtManager *JWTManager) Service {
	return &service{
		repo:       repo,
		jwtManager: jwtManager,
	}
}

// Authenticate authenticates a user and returns tokens
// Requirements: 12.1 - WHEN a parent enters NISN and password, THE System SHALL authenticate and return JWT tokens
func (s *service) Authenticate(ctx context.Context, username, password string) (*LoginResponse, error) {
	// Find user by username
	user, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	// Check if user is active
	// Requirements: 4.4 - Deactivated accounts cannot login
	if !user.IsActive {
		return nil, ErrAccountInactive
	}

	// Check if school is active (for non-super_admin users)
	// Requirements: 1.3 - Deactivated tenant users cannot access the system
	if user.School != nil && !user.School.IsActive {
		return nil, ErrSchoolInactive
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Generate token pair
	tokenClaims := TokenClaims{
		UserID:   user.ID,
		SchoolID: user.SchoolID,
		Role:     string(user.Role),
		Username: user.Username,
	}

	tokenPair, err := s.jwtManager.GenerateTokenPair(tokenClaims)
	if err != nil {
		return nil, err
	}

	// Update last login timestamp
	_ = s.repo.UpdateLastLogin(ctx, user.ID)

	return &LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
		TokenType:    "Bearer",
		User:         toUserResponse(user),
	}, nil
}

// RefreshAccessToken refreshes the access token using a refresh token
// Requirements: 12.3 - Token refresh functionality
func (s *service) RefreshAccessToken(ctx context.Context, refreshToken string) (*RefreshTokenResponse, error) {
	// Validate refresh token
	claims, err := s.jwtManager.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	// Verify user still exists and is active
	user, err := s.repo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	if !user.IsActive {
		return nil, ErrAccountInactive
	}

	// Check if school is still active
	if user.School != nil && !user.School.IsActive {
		return nil, ErrSchoolInactive
	}

	// Generate new token pair
	tokenClaims := TokenClaims{
		UserID:   user.ID,
		SchoolID: user.SchoolID,
		Role:     string(user.Role),
		Username: user.Username,
	}

	tokenPair, err := s.jwtManager.GenerateTokenPair(tokenClaims)
	if err != nil {
		return nil, err
	}

	return &RefreshTokenResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
		TokenType:    "Bearer",
	}, nil
}

// ChangePassword changes the user's password
// Requirements: 12.5 - THE System SHALL enforce password reset on first login
func (s *service) ChangePassword(ctx context.Context, userID uint, oldPassword, newPassword string) error {
	// Validate new password length
	if len(newPassword) < 8 {
		return ErrPasswordTooShort
	}

	// Get user
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return ErrPasswordMismatch
	}

	// Check if new password is different from old
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(newPassword)); err == nil {
		return ErrSamePassword
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update password
	if err := s.repo.UpdatePassword(ctx, userID, string(hashedPassword)); err != nil {
		return err
	}

	// Clear password reset flag
	if err := s.repo.ClearPasswordReset(ctx, userID); err != nil {
		return err
	}

	return nil
}

// GetUserByID retrieves a user by ID
func (s *service) GetUserByID(ctx context.Context, userID uint) (*models.User, error) {
	return s.repo.FindByID(ctx, userID)
}

// toUserResponse converts a User model to UserResponse DTO
func toUserResponse(user *models.User) UserResponse {
	return UserResponse{
		ID:           user.ID,
		SchoolID:     user.SchoolID,
		Role:         string(user.Role),
		Username:     user.Username,
		Email:        user.Email,
		Name:         user.Name,
		IsActive:     user.IsActive,
		MustResetPwd: user.MustResetPwd,
		LastLoginAt:  user.LastLoginAt,
	}
}

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
