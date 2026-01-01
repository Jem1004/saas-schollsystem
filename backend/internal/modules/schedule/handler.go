package schedule

import (
	"errors"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Handler handles HTTP requests for schedule management
type Handler struct {
	service Service
}

// NewHandler creates a new schedule handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers schedule routes
// Requirements: 3.1 - Schedule management endpoints
func (h *Handler) RegisterRoutes(router fiber.Router) {
	schedules := router.Group("/schedules")
	schedules.Get("", h.GetAllSchedules)
	schedules.Post("", h.CreateSchedule)
	schedules.Get("/active", h.GetActiveSchedule)
	schedules.Get("/:id", h.GetScheduleByID)
	schedules.Put("/:id", h.UpdateSchedule)
	schedules.Delete("/:id", h.DeleteSchedule)
	schedules.Post("/:id/default", h.SetDefaultSchedule)
}

// RegisterRoutesWithoutGroup registers schedule routes without creating a sub-group
func (h *Handler) RegisterRoutesWithoutGroup(router fiber.Router) {
	router.Get("", h.GetAllSchedules)
	router.Post("", h.CreateSchedule)
	router.Get("/active", h.GetActiveSchedule)
	router.Get("/:id", h.GetScheduleByID)
	router.Put("/:id", h.UpdateSchedule)
	router.Delete("/:id", h.DeleteSchedule)
	router.Post("/:id/default", h.SetDefaultSchedule)
}

// CreateSchedule handles creating a new schedule
// @Summary Create attendance schedule
// @Description Create a new attendance schedule for the school
// @Tags Schedules
// @Accept json
// @Produce json
// @Param request body CreateScheduleRequest true "Schedule data"
// @Success 201 {object} ScheduleResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/schedules [post]
func (h *Handler) CreateSchedule(c *fiber.Ctx) error {
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

	var req CreateScheduleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Format data tidak valid",
			},
		})
	}

	response, err := h.service.CreateSchedule(c.Context(), schoolID, req)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    response,
		"message": "Jadwal absensi berhasil dibuat",
	})
}

// GetAllSchedules handles listing all schedules
// @Summary List attendance schedules
// @Description Get all attendance schedules for the school
// @Tags Schedules
// @Produce json
// @Success 200 {object} ScheduleListResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/schedules [get]
func (h *Handler) GetAllSchedules(c *fiber.Ctx) error {
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

	response, err := h.service.GetAllSchedules(c.Context(), schoolID)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}


// GetScheduleByID handles getting a single schedule
// @Summary Get schedule by ID
// @Description Get detailed information about a specific schedule
// @Tags Schedules
// @Produce json
// @Param id path int true "Schedule ID"
// @Success 200 {object} ScheduleResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/schedules/{id} [get]
func (h *Handler) GetScheduleByID(c *fiber.Ctx) error {
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
				"message": "ID jadwal tidak valid",
			},
		})
	}

	response, err := h.service.GetScheduleByID(c.Context(), schoolID, uint(id))
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// UpdateSchedule handles updating a schedule
// @Summary Update attendance schedule
// @Description Update an existing attendance schedule
// @Tags Schedules
// @Accept json
// @Produce json
// @Param id path int true "Schedule ID"
// @Param request body UpdateScheduleRequest true "Schedule data"
// @Success 200 {object} ScheduleResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/schedules/{id} [put]
func (h *Handler) UpdateSchedule(c *fiber.Ctx) error {
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
				"message": "ID jadwal tidak valid",
			},
		})
	}

	var req UpdateScheduleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Format data tidak valid",
			},
		})
	}

	response, err := h.service.UpdateSchedule(c.Context(), schoolID, uint(id), req)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
		"message": "Jadwal absensi berhasil diperbarui",
	})
}

// DeleteSchedule handles deleting a schedule
// @Summary Delete attendance schedule
// @Description Delete an attendance schedule
// @Tags Schedules
// @Produce json
// @Param id path int true "Schedule ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/schedules/{id} [delete]
func (h *Handler) DeleteSchedule(c *fiber.Ctx) error {
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
				"message": "ID jadwal tidak valid",
			},
		})
	}

	if err := h.service.DeleteSchedule(c.Context(), schoolID, uint(id)); err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Jadwal absensi berhasil dihapus",
	})
}

// GetActiveSchedule handles getting the currently active schedule
// @Summary Get active schedule
// @Description Get the currently active attendance schedule based on current time
// @Tags Schedules
// @Produce json
// @Success 200 {object} ActiveScheduleResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/schedules/active [get]
func (h *Handler) GetActiveSchedule(c *fiber.Ctx) error {
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

	response, err := h.service.GetActiveSchedule(c.Context(), schoolID, time.Now())
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// SetDefaultSchedule handles setting a schedule as default
// @Summary Set default schedule
// @Description Set a schedule as the default for the school
// @Tags Schedules
// @Produce json
// @Param id path int true "Schedule ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/schedules/{id}/default [post]
func (h *Handler) SetDefaultSchedule(c *fiber.Ctx) error {
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
				"message": "ID jadwal tidak valid",
			},
		})
	}

	if err := h.service.SetDefaultSchedule(c.Context(), schoolID, uint(id)); err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Jadwal berhasil ditetapkan sebagai default",
	})
}

// handleError handles service errors and returns appropriate HTTP responses
func (h *Handler) handleError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, ErrScheduleNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "NOT_FOUND_SCHEDULE",
				"message": "Jadwal absensi tidak ditemukan",
			},
		})
	case errors.Is(err, ErrScheduleLimitExceeded):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "SCHEDULE_LIMIT_EXCEEDED",
				"message": "Batas maksimum jadwal (10) telah tercapai",
			},
		})
	case errors.Is(err, ErrScheduleTimeOverlap):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "SCHEDULE_TIME_OVERLAP",
				"message": "Waktu jadwal bertumpang tindih dengan jadwal lain",
			},
		})
	case errors.Is(err, ErrScheduleInUse):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "SCHEDULE_IN_USE",
				"message": "Jadwal tidak dapat dihapus karena masih digunakan oleh data kehadiran",
			},
		})
	case errors.Is(err, ErrNameRequired):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "Nama jadwal wajib diisi",
			},
		})
	case errors.Is(err, ErrStartTimeRequired):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "Waktu mulai wajib diisi",
			},
		})
	case errors.Is(err, ErrEndTimeRequired):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "Waktu akhir wajib diisi",
			},
		})
	case errors.Is(err, ErrInvalidStartTime), errors.Is(err, ErrInvalidEndTime):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Format waktu tidak valid (gunakan HH:MM)",
			},
		})
	case errors.Is(err, ErrEndTimeBeforeStart):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_TIME",
				"message": "Waktu akhir harus setelah waktu mulai",
			},
		})
	case errors.Is(err, ErrVeryLateThreshold):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_THRESHOLD",
				"message": "Batas sangat terlambat harus lebih besar dari batas terlambat",
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
