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
	schools := router.Group("/schools")
	schools.Post("", h.CreateSchool)
	schools.Get("", h.GetSchools)
	schools.Get("/:id", h.GetSchool)
	schools.Put("/:id", h.UpdateSchool)
	schools.Post("/:id/deactivate", h.DeactivateSchool)
	schools.Post("/:id/activate", h.ActivateSchool)
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
				"message": "Invalid request body",
			},
		})
	}

	// Validate required fields
	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "School name is required",
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
		"message": "School created successfully",
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
				"message": "Invalid school ID",
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
				"message": "Invalid school ID",
			},
		})
	}

	var req UpdateSchoolRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Invalid request body",
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
		"message": "School updated successfully",
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
				"message": "Invalid school ID",
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
				"message": "Invalid school ID",
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

// handleError handles service errors and returns appropriate HTTP responses
func (h *Handler) handleError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, ErrSchoolNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "NOT_FOUND_SCHOOL",
				"message": "School not found",
			},
		})
	case errors.Is(err, ErrDuplicateSchool):
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_DUPLICATE_ENTRY",
				"message": "A school with this name already exists",
			},
		})
	case errors.Is(err, ErrNameRequired):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "School name is required",
			},
		})
	case errors.Is(err, ErrSchoolInactive):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_STATE",
				"message": "School is already inactive",
			},
		})
	case errors.Is(err, ErrSchoolActive):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_STATE",
				"message": "School is already active",
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
