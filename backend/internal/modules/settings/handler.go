package settings

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Handler handles HTTP requests for SchoolSettings management
type Handler struct {
	service Service
}

// NewHandler creates a new Settings handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers Settings routes for Admin Sekolah
func (h *Handler) RegisterRoutes(router fiber.Router) {
	settings := router.Group("/settings")

	// Settings CRUD
	settings.Get("", h.GetSettings)
	settings.Put("", h.UpdateSettings)
	settings.Post("/reset", h.ResetToDefaults)

	// Partial updates
	settings.Put("/attendance", h.UpdateAttendanceSettings)
	settings.Put("/notifications", h.UpdateNotificationSettings)
	settings.Put("/academic", h.UpdateAcademicSettings)

	// Utility endpoints
	settings.Get("/attendance-window", h.GetAttendanceTimeWindow)
}

// ==================== Settings Handlers ====================

// GetSettings handles getting school settings
// @Summary Get school settings
// @Description Get settings for the current school
// @Tags Settings
// @Produce json
// @Success 200 {object} SettingsResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/settings [get]
func (h *Handler) GetSettings(c *fiber.Ctx) error {
	schoolID, ok := c.Locals("school_id").(uint)
	if !ok {
		return h.tenantRequiredError(c)
	}

	response, err := h.service.GetSchoolSettings(c.Context(), schoolID)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// UpdateSettings handles updating all school settings
// @Summary Update school settings
// @Description Update all settings for the current school
// @Tags Settings
// @Accept json
// @Produce json
// @Param request body UpdateSettingsRequest true "Settings data"
// @Success 200 {object} SettingsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/settings [put]
func (h *Handler) UpdateSettings(c *fiber.Ctx) error {
	schoolID, ok := c.Locals("school_id").(uint)
	if !ok {
		return h.tenantRequiredError(c)
	}

	var req UpdateSettingsRequest
	if err := c.BodyParser(&req); err != nil {
		return h.invalidBodyError(c)
	}

	response, err := h.service.UpdateSchoolSettings(c.Context(), schoolID, req)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
		"message": "Settings updated successfully",
	})
}

// ResetToDefaults handles resetting settings to defaults
// @Summary Reset settings to defaults
// @Description Reset all settings to default values
// @Tags Settings
// @Produce json
// @Success 200 {object} SettingsResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/settings/reset [post]
func (h *Handler) ResetToDefaults(c *fiber.Ctx) error {
	schoolID, ok := c.Locals("school_id").(uint)
	if !ok {
		return h.tenantRequiredError(c)
	}

	response, err := h.service.ResetToDefaults(c.Context(), schoolID)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
		"message": "Settings reset to defaults successfully",
	})
}

// UpdateAttendanceSettings handles updating attendance settings only
// @Summary Update attendance settings
// @Description Update only attendance-related settings
// @Tags Settings
// @Accept json
// @Produce json
// @Param request body UpdateAttendanceSettingsRequest true "Attendance settings data"
// @Success 200 {object} SettingsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/settings/attendance [put]
func (h *Handler) UpdateAttendanceSettings(c *fiber.Ctx) error {
	schoolID, ok := c.Locals("school_id").(uint)
	if !ok {
		return h.tenantRequiredError(c)
	}

	var req UpdateAttendanceSettingsRequest
	if err := c.BodyParser(&req); err != nil {
		return h.invalidBodyError(c)
	}

	response, err := h.service.UpdateAttendanceSettings(c.Context(), schoolID, req)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
		"message": "Attendance settings updated successfully",
	})
}


// UpdateNotificationSettings handles updating notification settings only
// @Summary Update notification settings
// @Description Update only notification-related settings
// @Tags Settings
// @Accept json
// @Produce json
// @Param request body UpdateNotificationSettingsRequest true "Notification settings data"
// @Success 200 {object} SettingsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/settings/notifications [put]
func (h *Handler) UpdateNotificationSettings(c *fiber.Ctx) error {
	schoolID, ok := c.Locals("school_id").(uint)
	if !ok {
		return h.tenantRequiredError(c)
	}

	var req UpdateNotificationSettingsRequest
	if err := c.BodyParser(&req); err != nil {
		return h.invalidBodyError(c)
	}

	response, err := h.service.UpdateNotificationSettings(c.Context(), schoolID, req)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
		"message": "Notification settings updated successfully",
	})
}

// UpdateAcademicSettings handles updating academic settings only
// @Summary Update academic settings
// @Description Update only academic-related settings
// @Tags Settings
// @Accept json
// @Produce json
// @Param request body UpdateAcademicSettingsRequest true "Academic settings data"
// @Success 200 {object} SettingsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/settings/academic [put]
func (h *Handler) UpdateAcademicSettings(c *fiber.Ctx) error {
	schoolID, ok := c.Locals("school_id").(uint)
	if !ok {
		return h.tenantRequiredError(c)
	}

	var req UpdateAcademicSettingsRequest
	if err := c.BodyParser(&req); err != nil {
		return h.invalidBodyError(c)
	}

	response, err := h.service.UpdateAcademicSettings(c.Context(), schoolID, req)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
		"message": "Academic settings updated successfully",
	})
}

// GetAttendanceTimeWindow handles getting the attendance time window
// @Summary Get attendance time window
// @Description Get the calculated attendance time window for a specific date
// @Tags Settings
// @Produce json
// @Param date query string false "Date (YYYY-MM-DD), defaults to today"
// @Success 200 {object} AttendanceTimeWindowResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/settings/attendance-window [get]
func (h *Handler) GetAttendanceTimeWindow(c *fiber.Ctx) error {
	schoolID, ok := c.Locals("school_id").(uint)
	if !ok {
		return h.tenantRequiredError(c)
	}

	// Parse date parameter, default to today
	dateStr := c.Query("date")
	var date time.Time
	if dateStr != "" {
		var err error
		date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    "VAL_INVALID_FORMAT",
					"message": "Invalid date format. Use YYYY-MM-DD",
				},
			})
		}
	} else {
		date = time.Now()
	}

	response, err := h.service.GetAttendanceTimeWindow(c.Context(), schoolID, date)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// ==================== Error Handlers ====================

func (h *Handler) tenantRequiredError(c *fiber.Ctx) error {
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
		"success": false,
		"error": fiber.Map{
			"code":    "AUTHZ_TENANT_REQUIRED",
			"message": "School context is required",
		},
	})
}

func (h *Handler) invalidBodyError(c *fiber.Ctx) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"success": false,
		"error": fiber.Map{
			"code":    "VAL_INVALID_FORMAT",
			"message": "Invalid request body",
		},
	})
}

func (h *Handler) handleError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, ErrSettingsNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "NOT_FOUND_SETTINGS",
				"message": "Settings not found",
			},
		})
	case errors.Is(err, ErrSchoolNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "NOT_FOUND_SCHOOL",
				"message": "School not found",
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
	case errors.Is(err, ErrInvalidTimeFormat):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Time must be in HH:MM format",
			},
		})
	case errors.Is(err, ErrInvalidLateThreshold):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_VALUE",
				"message": "Late threshold must be non-negative",
			},
		})
	case errors.Is(err, ErrInvalidVeryLateThreshold):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_VALUE",
				"message": "Very late threshold must be greater than or equal to late threshold",
			},
		})
	case errors.Is(err, ErrInvalidSemester):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_VALUE",
				"message": "Semester must be 1 or 2",
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
