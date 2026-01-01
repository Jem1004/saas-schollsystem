package displaytoken

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/school-management/backend/internal/domain/models"
)

var (
	ErrNameRequired     = errors.New("nama display wajib diisi")
	ErrNameTooLong      = errors.New("nama display maksimal 100 karakter")
	ErrTokenGeneration  = errors.New("gagal membuat token")
	ErrInvalidExpiry    = errors.New("tanggal kedaluwarsa tidak valid")
)

// Service defines the interface for display token business logic
type Service interface {
	// CRUD operations
	CreateToken(ctx context.Context, schoolID uint, req CreateDisplayTokenRequest, baseURL string) (*DisplayTokenWithSecretResponse, error)
	GetAllTokens(ctx context.Context, schoolID uint, baseURL string) (*DisplayTokenListResponse, error)
	GetTokenByID(ctx context.Context, schoolID, id uint, baseURL string) (*DisplayTokenResponse, error)
	UpdateToken(ctx context.Context, schoolID, id uint, req UpdateDisplayTokenRequest) (*DisplayTokenResponse, error)
	DeleteToken(ctx context.Context, schoolID, id uint) error

	// Token operations
	// Requirements: 5.10, 5.13 - Token validation for access control
	ValidateToken(ctx context.Context, token string) (*DisplayTokenValidation, error)

	// Requirements: 6.4 - Revoke token
	RevokeToken(ctx context.Context, schoolID, id uint) error

	// Requirements: 6.5 - Regenerate token
	RegenerateToken(ctx context.Context, schoolID, id uint, baseURL string) (*DisplayTokenWithSecretResponse, error)

	// Requirements: 6.7 - Update last accessed
	UpdateLastAccessed(ctx context.Context, token string) error
}

// service implements the Service interface
type service struct {
	repo Repository
}

// NewService creates a new display token service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// CreateToken creates a new display token with secure generation
// Requirements: 5.1, 6.2 - Token creation with cryptographically secure random token
// Requirements: 6.3 - Show the full token only once
func (s *service) CreateToken(ctx context.Context, schoolID uint, req CreateDisplayTokenRequest, baseURL string) (*DisplayTokenWithSecretResponse, error) {
	// Validate request
	if err := s.validateCreateRequest(req); err != nil {
		return nil, err
	}

	// Generate secure token
	tokenValue, err := models.GenerateToken()
	if err != nil {
		return nil, ErrTokenGeneration
	}

	// Create token model
	token := &models.DisplayToken{
		SchoolID:  schoolID,
		Token:     tokenValue,
		Name:      strings.TrimSpace(req.Name),
		IsActive:  true,
		ExpiresAt: req.ExpiresAt,
	}

	// Validate the model
	if err := token.Validate(); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, token); err != nil {
		return nil, err
	}

	return toDisplayTokenWithSecretResponse(token, baseURL), nil
}

// GetAllTokens retrieves all display tokens for a school
// Requirements: 6.1 - List all tokens with name, status, and last accessed time
func (s *service) GetAllTokens(ctx context.Context, schoolID uint, baseURL string) (*DisplayTokenListResponse, error) {
	tokens, err := s.repo.FindAll(ctx, schoolID)
	if err != nil {
		return nil, err
	}

	responses := make([]DisplayTokenResponse, len(tokens))
	for i, token := range tokens {
		responses[i] = *toDisplayTokenResponse(&token, baseURL)
	}

	return &DisplayTokenListResponse{
		Tokens: responses,
		Total:  len(responses),
	}, nil
}

// GetTokenByID retrieves a display token by ID
func (s *service) GetTokenByID(ctx context.Context, schoolID, id uint, baseURL string) (*DisplayTokenResponse, error) {
	token, err := s.repo.FindByID(ctx, schoolID, id)
	if err != nil {
		return nil, err
	}

	return toDisplayTokenResponse(token, baseURL), nil
}

// UpdateToken updates an existing display token
func (s *service) UpdateToken(ctx context.Context, schoolID, id uint, req UpdateDisplayTokenRequest) (*DisplayTokenResponse, error) {
	// Get existing token
	token, err := s.repo.FindByID(ctx, schoolID, id)
	if err != nil {
		return nil, err
	}

	// Apply updates
	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return nil, ErrNameRequired
		}
		if len(name) > 100 {
			return nil, ErrNameTooLong
		}
		token.Name = name
	}
	if req.ExpiresAt != nil {
		token.ExpiresAt = req.ExpiresAt
	}
	if req.IsActive != nil {
		token.IsActive = *req.IsActive
	}

	// Update in database
	if err := s.repo.Update(ctx, token); err != nil {
		return nil, err
	}

	return toDisplayTokenResponse(token, ""), nil
}

// DeleteToken deletes a display token
// Requirements: 6.6 - Permanently remove token from the system
func (s *service) DeleteToken(ctx context.Context, schoolID, id uint) error {
	// Check if token exists
	_, err := s.repo.FindByID(ctx, schoolID, id)
	if err != nil {
		return err
	}

	return s.repo.Delete(ctx, schoolID, id)
}

// ValidateToken validates a display token for public display access
// Requirements: 5.10, 5.13 - Token validation for access control
func (s *service) ValidateToken(ctx context.Context, token string) (*DisplayTokenValidation, error) {
	displayToken, err := s.repo.FindByToken(ctx, token)
	if err != nil {
		if errors.Is(err, ErrTokenInvalid) {
			return &DisplayTokenValidation{
				Valid: false,
				Error: "Token tidak valid atau tidak ditemukan",
			}, nil
		}
		return nil, err
	}

	// Check if token is active
	if !displayToken.IsActive {
		return &DisplayTokenValidation{
			Valid: false,
			Error: "Token telah dicabut",
		}, nil
	}

	// Check if token is expired
	if displayToken.IsExpired() {
		return &DisplayTokenValidation{
			Valid: false,
			Error: "Token telah kedaluwarsa",
		}, nil
	}

	return &DisplayTokenValidation{
		Valid:    true,
		SchoolID: displayToken.SchoolID,
		TokenID:  displayToken.ID,
		Name:     displayToken.Name,
	}, nil
}

// RevokeToken revokes a display token
// Requirements: 5.10, 6.4 - Revoking a token SHALL immediately invalidate access
func (s *service) RevokeToken(ctx context.Context, schoolID, id uint) error {
	token, err := s.repo.FindByID(ctx, schoolID, id)
	if err != nil {
		return err
	}

	token.Revoke()
	return s.repo.Update(ctx, token)
}

// RegenerateToken regenerates a display token with a new value
// Requirements: 6.5 - Regenerating a token SHALL create a new token and invalidate the old one
func (s *service) RegenerateToken(ctx context.Context, schoolID, id uint, baseURL string) (*DisplayTokenWithSecretResponse, error) {
	token, err := s.repo.FindByID(ctx, schoolID, id)
	if err != nil {
		return nil, err
	}

	// Regenerate the token value
	if err := token.Regenerate(); err != nil {
		return nil, ErrTokenGeneration
	}

	// Update in database
	if err := s.repo.Update(ctx, token); err != nil {
		return nil, err
	}

	return toDisplayTokenWithSecretResponse(token, baseURL), nil
}

// UpdateLastAccessed updates the last accessed timestamp for a token
// Requirements: 6.7 - Track last accessed timestamp for each display token
func (s *service) UpdateLastAccessed(ctx context.Context, token string) error {
	displayToken, err := s.repo.FindByToken(ctx, token)
	if err != nil {
		return err
	}

	return s.repo.UpdateLastAccessed(ctx, displayToken.ID)
}

// validateCreateRequest validates the create display token request
func (s *service) validateCreateRequest(req CreateDisplayTokenRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return ErrNameRequired
	}
	if len(req.Name) > 100 {
		return ErrNameTooLong
	}
	return nil
}

// toDisplayTokenResponse converts a model to response DTO (without token value)
func toDisplayTokenResponse(token *models.DisplayToken, baseURL string) *DisplayTokenResponse {
	response := &DisplayTokenResponse{
		ID:             token.ID,
		SchoolID:       token.SchoolID,
		Name:           token.Name,
		IsActive:       token.IsActive,
		LastAccessedAt: token.LastAccessedAt,
		ExpiresAt:      token.ExpiresAt,
		CreatedAt:      token.CreatedAt,
		UpdatedAt:      token.UpdatedAt,
	}

	// Generate display URL if baseURL is provided
	if baseURL != "" {
		response.DisplayURL = fmt.Sprintf("%s/display/%s", baseURL, token.Token)
	}

	return response
}

// toDisplayTokenWithSecretResponse converts a model to response DTO with token value
// Requirements: 6.3 - Show the full token only once
func toDisplayTokenWithSecretResponse(token *models.DisplayToken, baseURL string) *DisplayTokenWithSecretResponse {
	displayURL := ""
	if baseURL != "" {
		displayURL = fmt.Sprintf("%s/display/%s", baseURL, token.Token)
	}

	return &DisplayTokenWithSecretResponse{
		ID:         token.ID,
		SchoolID:   token.SchoolID,
		Token:      token.Token,
		Name:       token.Name,
		IsActive:   token.IsActive,
		ExpiresAt:  token.ExpiresAt,
		CreatedAt:  token.CreatedAt,
		DisplayURL: displayURL,
	}
}
