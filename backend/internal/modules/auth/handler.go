package auth

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

// Handler handles HTTP requests for authentication
type Handler struct {
	service Service
}

// NewHandler creates a new auth handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers auth routes
func (h *Handler) RegisterRoutes(router fiber.Router) {
	auth := router.Group("/auth")
	auth.Post("/login", h.Login)
	auth.Post("/refresh", h.RefreshToken)
	auth.Post("/logout", h.Logout)
}

// RegisterProtectedRoutes registers routes that require authentication
func (h *Handler) RegisterProtectedRoutes(router fiber.Router) {
	auth := router.Group("/auth")
	auth.Post("/change-password", h.ChangePassword)
	auth.Get("/me", h.GetCurrentUser)
}

// Login handles user login
// @Summary User login
// @Description Authenticate user and return JWT tokens
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/v1/auth/login [post]
func (h *Handler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Invalid request body",
			},
		})
	}

	// Validate required fields
	if req.Username == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "Username and password are required",
			},
		})
	}

	// Authenticate user
	response, err := h.service.Authenticate(c.Context(), req.Username, req.Password)
	if err != nil {
		return h.handleAuthError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// RefreshToken handles token refresh
// @Summary Refresh access token
// @Description Refresh access token using refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "Refresh token"
// @Success 200 {object} RefreshTokenResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/v1/auth/refresh [post]
func (h *Handler) RefreshToken(c *fiber.Ctx) error {
	var req RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Invalid request body",
			},
		})
	}

	if req.RefreshToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "Refresh token is required",
			},
		})
	}

	response, err := h.service.RefreshAccessToken(c.Context(), req.RefreshToken)
	if err != nil {
		return h.handleAuthError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// Logout handles user logout
// @Summary User logout
// @Description Logout user (client should discard tokens)
// @Tags Auth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/auth/logout [post]
func (h *Handler) Logout(c *fiber.Ctx) error {
	// For JWT-based auth, logout is handled client-side by discarding tokens
	// Server-side token revocation can be implemented using Redis blacklist if needed
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Logged out successfully",
	})
}

// ChangePassword handles password change
// @Summary Change password
// @Description Change user password (requires authentication)
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body ChangePasswordRequest true "Password change request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/auth/change-password [post]
func (h *Handler) ChangePassword(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTH_TOKEN_INVALID",
				"message": "Invalid authentication",
			},
		})
	}

	var req ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Invalid request body",
			},
		})
	}

	if req.OldPassword == "" || req.NewPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "Old password and new password are required",
			},
		})
	}

	if len(req.NewPassword) < 8 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "New password must be at least 8 characters",
			},
		})
	}

	err := h.service.ChangePassword(c.Context(), userID, req.OldPassword, req.NewPassword)
	if err != nil {
		return h.handleAuthError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Password changed successfully",
	})
}

// GetCurrentUser returns the current authenticated user
// @Summary Get current user
// @Description Get the currently authenticated user's information
// @Tags Auth
// @Produce json
// @Success 200 {object} UserResponse
// @Failure 401 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/auth/me [get]
func (h *Handler) GetCurrentUser(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTH_TOKEN_INVALID",
				"message": "Invalid authentication",
			},
		})
	}

	user, err := h.service.GetUserByID(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "NOT_FOUND_USER",
				"message": "User not found",
			},
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    toUserResponse(user),
	})
}

// handleAuthError handles authentication errors and returns appropriate responses
func (h *Handler) handleAuthError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, ErrInvalidCredentials):
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTH_INVALID_CREDENTIALS",
				"message": "Invalid username or password",
			},
		})
	case errors.Is(err, ErrAccountInactive):
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTH_ACCOUNT_INACTIVE",
				"message": "Account is inactive",
			},
		})
	case errors.Is(err, ErrSchoolInactive):
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTH_SCHOOL_INACTIVE",
				"message": "School is inactive",
			},
		})
	case errors.Is(err, ErrTokenExpired):
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTH_TOKEN_EXPIRED",
				"message": "Token has expired",
			},
		})
	case errors.Is(err, ErrTokenInvalid), errors.Is(err, ErrTokenMalformed):
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTH_TOKEN_INVALID",
				"message": "Invalid token",
			},
		})
	case errors.Is(err, ErrPasswordMismatch):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTH_PASSWORD_MISMATCH",
				"message": "Old password is incorrect",
			},
		})
	case errors.Is(err, ErrSamePassword):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTH_SAME_PASSWORD",
				"message": "New password must be different from old password",
			},
		})
	case errors.Is(err, ErrPasswordTooShort):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Password must be at least 8 characters",
			},
		})
	case errors.Is(err, ErrUserNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "NOT_FOUND_USER",
				"message": "User not found",
			},
		})
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "INTERNAL_ERROR",
				"message": "An internal error occurred",
			},
		})
	}
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Success bool `json:"success"`
	Error   struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}
