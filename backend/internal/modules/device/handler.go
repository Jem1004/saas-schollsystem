package device

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Handler handles HTTP requests for device management
type Handler struct {
	service Service
}

// NewHandler creates a new device handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers device routes (Super Admin only)
// Requirements: 2.1, 2.4, 2.5 - Super Admin manages devices
func (h *Handler) RegisterRoutes(router fiber.Router) {
	devices := router.Group("/devices")
	devices.Post("", h.RegisterDevice)
	devices.Get("", h.GetDevices)
	devices.Get("/:id", h.GetDevice)
	devices.Put("/:id", h.UpdateDevice)
	devices.Post("/:id/revoke", h.RevokeAPIKey)
	devices.Post("/:id/regenerate", h.RegenerateAPIKey)
	devices.Delete("/:id", h.DeleteDevice)
}

// RegisterPublicRoutes registers public device routes (for ESP32 devices)
func (h *Handler) RegisterPublicRoutes(router fiber.Router) {
	router.Post("/devices/validate-key", h.ValidateAPIKey)
}

// RegisterDevice handles device registration
// @Summary Register a new device
// @Description Register a new RFID device (ESP32) and generate an API key
// @Tags Devices
// @Accept json
// @Produce json
// @Param request body RegisterDeviceRequest true "Device data"
// @Success 201 {object} DeviceWithAPIKeyResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/devices [post]
func (h *Handler) RegisterDevice(c *fiber.Ctx) error {
	var req RegisterDeviceRequest
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
	if req.SchoolID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "School ID is required",
			},
		})
	}
	if req.DeviceCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "Device code is required",
			},
		})
	}

	response, err := h.service.RegisterDevice(c.Context(), req)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    response,
		"message": "Device registered successfully. Please save the API key securely - it will not be shown again.",
	})
}


// GetDevices handles listing all devices
// @Summary List all devices
// @Description Get a paginated list of all RFID devices
// @Tags Devices
// @Produce json
// @Param school_id query int false "Filter by school ID"
// @Param device_code query string false "Filter by device code"
// @Param is_active query bool false "Filter by active status"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} DeviceListResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/devices [get]
func (h *Handler) GetDevices(c *fiber.Ctx) error {
	filter := DefaultDeviceFilter()

	// Parse query parameters
	if schoolIDStr := c.Query("school_id"); schoolIDStr != "" {
		if schoolID, err := strconv.ParseUint(schoolIDStr, 10, 32); err == nil {
			id := uint(schoolID)
			filter.SchoolID = &id
		}
	}

	filter.DeviceCode = c.Query("device_code")

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

	response, err := h.service.GetAllDevices(c.Context(), filter)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GetDevice handles getting a single device
// @Summary Get a device by ID
// @Description Get detailed information about a specific device
// @Tags Devices
// @Produce json
// @Param id path int true "Device ID"
// @Success 200 {object} DeviceResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/devices/{id} [get]
func (h *Handler) GetDevice(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Invalid device ID",
			},
		})
	}

	response, err := h.service.GetDeviceByID(c.Context(), uint(id))
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// UpdateDevice handles updating a device
// @Summary Update a device
// @Description Update device information
// @Tags Devices
// @Accept json
// @Produce json
// @Param id path int true "Device ID"
// @Param request body UpdateDeviceRequest true "Device data to update"
// @Success 200 {object} DeviceResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/devices/{id} [put]
func (h *Handler) UpdateDevice(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Invalid device ID",
			},
		})
	}

	var req UpdateDeviceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Invalid request body",
			},
		})
	}

	response, err := h.service.UpdateDevice(c.Context(), uint(id), req)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
		"message": "Device updated successfully",
	})
}

// RevokeAPIKey handles revoking a device's API key
// @Summary Revoke device API key
// @Description Revoke a device's API key, preventing it from sending attendance data
// @Tags Devices
// @Produce json
// @Param id path int true "Device ID"
// @Success 200 {object} RevokeAPIKeyResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/devices/{id}/revoke [post]
func (h *Handler) RevokeAPIKey(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Invalid device ID",
			},
		})
	}

	response, err := h.service.RevokeAPIKey(c.Context(), uint(id))
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// RegenerateAPIKey handles regenerating a device's API key
// @Summary Regenerate device API key
// @Description Generate a new API key for a device
// @Tags Devices
// @Produce json
// @Param id path int true "Device ID"
// @Success 200 {object} RegenerateAPIKeyResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/devices/{id}/regenerate [post]
func (h *Handler) RegenerateAPIKey(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Invalid device ID",
			},
		})
	}

	response, err := h.service.RegenerateAPIKey(c.Context(), uint(id))
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// DeleteDevice handles deleting a device
// @Summary Delete a device
// @Description Permanently delete a device
// @Tags Devices
// @Produce json
// @Param id path int true "Device ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/devices/{id} [delete]
func (h *Handler) DeleteDevice(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Invalid device ID",
			},
		})
	}

	if err := h.service.DeleteDevice(c.Context(), uint(id)); err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Device deleted successfully",
	})
}

// ValidateAPIKey handles API key validation (for ESP32 devices)
// @Summary Validate device API key
// @Description Validate an API key and return device information
// @Tags Devices
// @Accept json
// @Produce json
// @Param request body map[string]string true "API key"
// @Success 200 {object} APIKeyValidationResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/devices/validate-key [post]
func (h *Handler) ValidateAPIKey(c *fiber.Ctx) error {
	var req struct {
		APIKey string `json:"api_key"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Invalid request body",
			},
		})
	}

	response, err := h.service.ValidateAPIKey(c.Context(), req.APIKey)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": response.Valid,
		"data":    response,
	})
}

// handleError handles service errors and returns appropriate HTTP responses
func (h *Handler) handleError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, ErrDeviceNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "NOT_FOUND_DEVICE",
				"message": "Device not found",
			},
		})
	case errors.Is(err, ErrDuplicateDeviceCode):
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_DUPLICATE_ENTRY",
				"message": "A device with this code already exists",
			},
		})
	case errors.Is(err, ErrSchoolIDRequired):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "School ID is required",
			},
		})
	case errors.Is(err, ErrDeviceCodeRequired):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "Device code is required",
			},
		})
	case errors.Is(err, ErrDeviceInactive):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_STATE",
				"message": "Device is already inactive",
			},
		})
	case errors.Is(err, ErrDeviceActive):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_STATE",
				"message": "Device is already active",
			},
		})
	case errors.Is(err, ErrAPIKeyGeneration):
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to generate API key",
			},
		})
	case errors.Is(err, ErrInvalidAPIKey):
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTH_INVALID_API_KEY",
				"message": "Invalid API key",
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
