package tenant

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Handler handles HTTP requests for tenant (school) management
type Handler struct {
	service Service
}

// NewHandler creates a new tenant handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers tenant routes (Super Admin only)
// Requirements: 1.1, 1.2, 1.3 - Super Admin manages tenants
func (h *Handler) RegisterRoutes(router fiber.Router) {
	// Register routes directly without sub-group to avoid potential routing issues
	router.Post("/schools", h.CreateSchool)
	router.Get("/schools", h.GetSchools)
	// Register more specific routes first to avoid route conflicts
	router.Get("/schools/:id/detail", h.GetSchoolDetail)
	router.Post("/schools/:id/deactivate", h.DeactivateSchool)
	router.Post("/schools/:id/activate", h.ActivateSchool)
	// Then register generic parameter routes
	router.Get("/schools/:id", h.GetSchool)
	router.Put("/schools/:id", h.UpdateSchool)
	router.Delete("/schools/:id", h.DeleteSchool)
}

// RegisterRoutesWithoutGroup registers tenant routes without creating a sub-group
// Use this when the router already has the correct path prefix
func (h *Handler) RegisterRoutesWithoutGroup(router fiber.Router) {
	router.Post("", h.CreateSchool)
	router.Get("", h.GetSchools)
	// Register more specific routes first to avoid route conflicts
	router.Get("/:id/detail", h.GetSchoolDetail)
	router.Post("/:id/deactivate", h.DeactivateSchool)
	router.Post("/:id/activate", h.ActivateSchool)
	// Then register generic parameter routes
	router.Get("/:id", h.GetSchool)
	router.Put("/:id", h.UpdateSchool)
	router.Delete("/:id", h.DeleteSchool)
}

// CreateSchool handles school creation
// @Summary Create a new school (tenant)
// @Description Create a new school in the multi-tenant system
// @Tags Schools
// @Accept json
// @Produce json
// @Param request body CreateSchoolRequest true "School data"
// @Success 201 {object} SchoolResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/schools [post]
func (h *Handler) CreateSchool(c *fiber.Ctx) error {
	var req CreateSchoolRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Format data tidak valid",
			},
		})
	}

	// Validate required fields
	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "Nama sekolah wajib diisi",
			},
		})
	}

	response, err := h.service.CreateSchool(c.Context(), req)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    response,
		"message": "Sekolah berhasil dibuat",
	})
}

// GetSchools handles listing all schools
// @Summary List all schools
// @Description Get a paginated list of all schools (tenants)
// @Tags Schools
// @Produce json
// @Param name query string false "Filter by school name"
// @Param is_active query bool false "Filter by active status"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} SchoolListResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/schools [get]
func (h *Handler) GetSchools(c *fiber.Ctx) error {
	filter := DefaultSchoolFilter()

	// Parse query parameters
	filter.Name = c.Query("name")
	
	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		isActive := isActiveStr == "true"
		filter.IsActive = &isActive
	}

	if page, err := strconv.Atoi(c.Query("page", "1")); err == nil && page > 0 {
		filter.Page = page
	}

	if pageSize, err := strconv.Atoi(c.Query("page_size", "20")); err == nil && pageSize > 0 {
		filter.PageSize = pageSize
	}

	response, err := h.service.GetAllSchools(c.Context(), filter)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GetSchool handles getting a single school
// @Summary Get a school by ID
// @Description Get detailed information about a specific school
// @Tags Schools
// @Produce json
// @Param id path int true "School ID"
// @Success 200 {object} SchoolResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/schools/{id} [get]
func (h *Handler) GetSchool(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "ID sekolah tidak valid",
			},
		})
	}

	response, err := h.service.GetSchoolByID(c.Context(), uint(id))
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GetSchoolDetail handles getting a school with admin info
// @Summary Get school detail with admin info
// @Description Get detailed information about a school including admin users
// @Tags Schools
// @Produce json
// @Param id path int true "School ID"
// @Success 200 {object} SchoolDetailResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/schools/{id}/detail [get]
func (h *Handler) GetSchoolDetail(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "ID sekolah tidak valid",
			},
		})
	}

	response, err := h.service.GetSchoolDetail(c.Context(), uint(id))
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// UpdateSchool handles updating a school
// @Summary Update a school
// @Description Update school information
// @Tags Schools
// @Accept json
// @Produce json
// @Param id path int true "School ID"
// @Param request body UpdateSchoolRequest true "School data to update"
// @Success 200 {object} SchoolResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/schools/{id} [put]
func (h *Handler) UpdateSchool(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "ID sekolah tidak valid",
			},
		})
	}

	var req UpdateSchoolRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Format data tidak valid",
			},
		})
	}

	response, err := h.service.UpdateSchool(c.Context(), uint(id), req)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
		"message": "Sekolah berhasil diperbarui",
	})
}

// DeactivateSchool handles deactivating a school
// @Summary Deactivate a school
// @Description Deactivate a school, preventing all users from accessing the system
// @Tags Schools
// @Produce json
// @Param id path int true "School ID"
// @Success 200 {object} ActivateDeactivateResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/schools/{id}/deactivate [post]
func (h *Handler) DeactivateSchool(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "ID sekolah tidak valid",
			},
		})
	}

	response, err := h.service.DeactivateSchool(c.Context(), uint(id))
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// ActivateSchool handles activating a school
// @Summary Activate a school
// @Description Activate a previously deactivated school
// @Tags Schools
// @Produce json
// @Param id path int true "School ID"
// @Success 200 {object} ActivateDeactivateResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/schools/{id}/activate [post]
func (h *Handler) ActivateSchool(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "ID sekolah tidak valid",
			},
		})
	}

	response, err := h.service.ActivateSchool(c.Context(), uint(id))
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// DeleteSchool handles deleting a school and all associated data
// @Summary Delete a school
// @Description Permanently delete a school and all associated data (cascade delete)
// @Tags Schools
// @Produce json
// @Param id path int true "School ID"
// @Success 200 {object} DeleteSchoolResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/schools/{id} [delete]
func (h *Handler) DeleteSchool(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "ID sekolah tidak valid",
			},
		})
	}

	response, err := h.service.DeleteSchool(c.Context(), uint(id))
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// handleError handles service errors and returns appropriate HTTP responses
func (h *Handler) handleError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, ErrSchoolNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "NOT_FOUND_SCHOOL",
				"message": "Sekolah tidak ditemukan",
			},
		})
	case errors.Is(err, ErrDuplicateSchool):
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_DUPLICATE_ENTRY",
				"message": "Sekolah dengan nama ini sudah terdaftar",
			},
		})
	case errors.Is(err, ErrNameRequired):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "Nama sekolah wajib diisi",
			},
		})
	case errors.Is(err, ErrSchoolInactive):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_STATE",
				"message": "Sekolah sudah nonaktif",
			},
		})
	case errors.Is(err, ErrSchoolActive):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_STATE",
				"message": "Sekolah sudah aktif",
			},
		})
	case errors.Is(err, ErrUsernameExists):
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_DUPLICATE_ENTRY",
				"message": "Username sudah digunakan",
			},
		})
	case errors.Is(err, ErrInvalidUsername):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Username hanya boleh berisi huruf, angka, dan underscore",
			},
		})
	default:
		// Return the actual error message for better debugging
		errMsg := err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "ERROR",
				"message": errMsg,
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
