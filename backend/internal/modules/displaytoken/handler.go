package displaytoken

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Handler handles HTTP requests for display token management
type Handler struct {
	service Service
}

// NewHandler creates a new display token handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers display token routes
// Requirements: 6.1 - Display token management endpoints
func (h *Handler) RegisterRoutes(router fiber.Router) {
	tokens := router.Group("/display-tokens")
	tokens.Get("", h.GetAllTokens)
	tokens.Post("", h.CreateToken)
	tokens.Get("/:id", h.GetTokenByID)
	tokens.Put("/:id", h.UpdateToken)
	tokens.Delete("/:id", h.DeleteToken)
	tokens.Post("/:id/revoke", h.RevokeToken)
	tokens.Post("/:id/regenerate", h.RegenerateToken)
}

// RegisterRoutesWithoutGroup registers display token routes without creating a sub-group
func (h *Handler) RegisterRoutesWithoutGroup(router fiber.Router) {
	router.Get("", h.GetAllTokens)
	router.Post("", h.CreateToken)
	router.Get("/:id", h.GetTokenByID)
	router.Put("/:id", h.UpdateToken)
	router.Delete("/:id", h.DeleteToken)
	router.Post("/:id/revoke", h.RevokeToken)
	router.Post("/:id/regenerate", h.RegenerateToken)
}

// getBaseURL extracts the base URL from the request
func (h *Handler) getBaseURL(c *fiber.Ctx) string {
	scheme := "http"
	if c.Protocol() == "https" {
		scheme = "https"
	}
	return scheme + "://" + c.Hostname()
}

// CreateToken handles creating a new display token
// @Summary Create display token
// @Description Create a new display token for public display access
// @Tags Display Tokens
// @Accept json
// @Produce json
// @Param request body CreateDisplayTokenRequest true "Token data"
// @Success 201 {object} DisplayTokenWithSecretResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/display-tokens [post]
func (h *Handler) CreateToken(c *fiber.Ctx) error {
	schoolID, ok := c.Locals("school_id").(uint)
	if !ok {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTHZ_TENANT_REQUIRED",
				"message": "Konteks sekolah diperlukan",
			},
		})
	}

	var req CreateDisplayTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Format data tidak valid",
			},
		})
	}

	baseURL := h.getBaseURL(c)
	response, err := h.service.CreateToken(c.Context(), schoolID, req, baseURL)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    response,
		"message": "Display token berhasil dibuat. Simpan token ini karena tidak akan ditampilkan lagi.",
	})
}

// GetAllTokens handles listing all display tokens
// @Summary List display tokens
// @Description Get all display tokens for the school
// @Tags Display Tokens
// @Produce json
// @Success 200 {object} DisplayTokenListResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/display-tokens [get]
func (h *Handler) GetAllTokens(c *fiber.Ctx) error {
	schoolID, ok := c.Locals("school_id").(uint)
	if !ok {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTHZ_TENANT_REQUIRED",
				"message": "Konteks sekolah diperlukan",
			},
		})
	}

	baseURL := h.getBaseURL(c)
	response, err := h.service.GetAllTokens(c.Context(), schoolID, baseURL)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GetTokenByID handles getting a single display token
// @Summary Get display token by ID
// @Description Get detailed information about a specific display token
// @Tags Display Tokens
// @Produce json
// @Param id path int true "Token ID"
// @Success 200 {object} DisplayTokenResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/display-tokens/{id} [get]
func (h *Handler) GetTokenByID(c *fiber.Ctx) error {
	schoolID, ok := c.Locals("school_id").(uint)
	if !ok {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTHZ_TENANT_REQUIRED",
				"message": "Konteks sekolah diperlukan",
			},
		})
	}

	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "ID token tidak valid",
			},
		})
	}

	baseURL := h.getBaseURL(c)
	response, err := h.service.GetTokenByID(c.Context(), schoolID, uint(id), baseURL)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// UpdateToken handles updating a display token
// @Summary Update display token
// @Description Update an existing display token
// @Tags Display Tokens
// @Accept json
// @Produce json
// @Param id path int true "Token ID"
// @Param request body UpdateDisplayTokenRequest true "Token data"
// @Success 200 {object} DisplayTokenResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/display-tokens/{id} [put]
func (h *Handler) UpdateToken(c *fiber.Ctx) error {
	schoolID, ok := c.Locals("school_id").(uint)
	if !ok {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTHZ_TENANT_REQUIRED",
				"message": "Konteks sekolah diperlukan",
			},
		})
	}

	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "ID token tidak valid",
			},
		})
	}

	var req UpdateDisplayTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Format data tidak valid",
			},
		})
	}

	response, err := h.service.UpdateToken(c.Context(), schoolID, uint(id), req)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
		"message": "Display token berhasil diperbarui",
	})
}

// DeleteToken handles deleting a display token
// @Summary Delete display token
// @Description Delete a display token permanently
// @Tags Display Tokens
// @Produce json
// @Param id path int true "Token ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/display-tokens/{id} [delete]
func (h *Handler) DeleteToken(c *fiber.Ctx) error {
	schoolID, ok := c.Locals("school_id").(uint)
	if !ok {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTHZ_TENANT_REQUIRED",
				"message": "Konteks sekolah diperlukan",
			},
		})
	}

	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "ID token tidak valid",
			},
		})
	}

	if err := h.service.DeleteToken(c.Context(), schoolID, uint(id)); err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Display token berhasil dihapus",
	})
}

// RevokeToken handles revoking a display token
// @Summary Revoke display token
// @Description Revoke a display token to immediately invalidate access
// @Tags Display Tokens
// @Produce json
// @Param id path int true "Token ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/display-tokens/{id}/revoke [post]
func (h *Handler) RevokeToken(c *fiber.Ctx) error {
	schoolID, ok := c.Locals("school_id").(uint)
	if !ok {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTHZ_TENANT_REQUIRED",
				"message": "Konteks sekolah diperlukan",
			},
		})
	}

	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "ID token tidak valid",
			},
		})
	}

	if err := h.service.RevokeToken(c.Context(), schoolID, uint(id)); err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Display token berhasil dicabut",
	})
}

// RegenerateToken handles regenerating a display token
// @Summary Regenerate display token
// @Description Regenerate a display token with a new value (invalidates the old token)
// @Tags Display Tokens
// @Produce json
// @Param id path int true "Token ID"
// @Success 200 {object} DisplayTokenWithSecretResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/display-tokens/{id}/regenerate [post]
func (h *Handler) RegenerateToken(c *fiber.Ctx) error {
	schoolID, ok := c.Locals("school_id").(uint)
	if !ok {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTHZ_TENANT_REQUIRED",
				"message": "Konteks sekolah diperlukan",
			},
		})
	}

	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "ID token tidak valid",
			},
		})
	}

	baseURL := h.getBaseURL(c)
	response, err := h.service.RegenerateToken(c.Context(), schoolID, uint(id), baseURL)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
		"message": "Display token berhasil di-regenerate. Simpan token baru ini karena tidak akan ditampilkan lagi.",
	})
}

// handleError handles service errors and returns appropriate HTTP responses
func (h *Handler) handleError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, ErrTokenNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "TOKEN_NOT_FOUND",
				"message": "Display token tidak ditemukan",
			},
		})
	case errors.Is(err, ErrTokenInvalid):
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "TOKEN_INVALID",
				"message": "Token tidak valid atau tidak ditemukan",
			},
		})
	case errors.Is(err, ErrTokenRevoked):
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "TOKEN_REVOKED",
				"message": "Token telah dicabut",
			},
		})
	case errors.Is(err, ErrTokenExpired):
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "TOKEN_EXPIRED",
				"message": "Token telah kedaluwarsa",
			},
		})
	case errors.Is(err, ErrNameRequired):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "Nama display wajib diisi",
			},
		})
	case errors.Is(err, ErrNameTooLong):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_LENGTH",
				"message": "Nama display maksimal 100 karakter",
			},
		})
	case errors.Is(err, ErrTokenGeneration):
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "TOKEN_GENERATION_FAILED",
				"message": "Gagal membuat token",
			},
		})
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "ERROR",
				"message": err.Error(),
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
