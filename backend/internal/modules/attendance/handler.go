package attendance

import (
	"errors"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/school-management/backend/internal/modules/device"
)

// Handler handles HTTP requests for attendance management
type Handler struct {
	service Service
}

// NewHandler creates a new attendance handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers attendance routes for admin/teachers
func (h *Handler) RegisterRoutes(router fiber.Router) {
	attendance := router.Group("/attendance")
	attendance.Get("", h.GetAttendance)
	attendance.Get("/summary", h.GetSchoolSummary)
	attendance.Get("/class/:classId", h.GetClassAttendance)
	attendance.Get("/student/:studentId", h.GetStudentAttendance)
	attendance.Get("/:id", h.GetAttendanceByID)
	attendance.Post("/manual", h.RecordManualAttendance)
	attendance.Post("/manual/bulk", h.RecordBulkManualAttendance)
	attendance.Delete("/:id", h.DeleteAttendance)
}

// RegisterPublicRoutes registers public routes for ESP32 devices
func (h *Handler) RegisterPublicRoutes(router fiber.Router) {
	router.Post("/attendance/rfid", h.RecordRFIDAttendance)
}

// RecordRFIDAttendance handles RFID attendance recording from ESP32 devices
// @Summary Record RFID attendance
// @Description Record student attendance via RFID card tap from ESP32 device
// @Tags Attendance
// @Accept json
// @Produce json
// @Param request body RFIDAttendanceRequest true "RFID attendance data"
// @Success 200 {object} RFIDAttendanceResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/v1/attendance/rfid [post]
func (h *Handler) RecordRFIDAttendance(c *fiber.Ctx) error {
	var req RFIDAttendanceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Invalid request body",
			},
		})
	}

	response, err := h.service.RecordRFIDAttendance(c.Context(), req)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": response.Success,
		"data":    response,
	})
}

// RecordManualAttendance handles manual attendance recording
// @Summary Record manual attendance
// @Description Record student attendance manually (fallback when RFID fails)
// @Tags Attendance
// @Accept json
// @Produce json
// @Param request body ManualAttendanceRequest true "Manual attendance data"
// @Success 201 {object} AttendanceResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/attendance/manual [post]
func (h *Handler) RecordManualAttendance(c *fiber.Ctx) error {
	// Get school ID from context (set by tenant middleware)
	schoolID, ok := c.Locals("school_id").(uint)
	if !ok {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTHZ_TENANT_REQUIRED",
				"message": "School context is required",
			},
		})
	}

	var req ManualAttendanceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Invalid request body",
			},
		})
	}

	response, err := h.service.RecordManualAttendance(c.Context(), schoolID, req)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    response,
		"message": "Attendance recorded successfully",
	})
}

// RecordBulkManualAttendance handles bulk manual attendance recording
// @Summary Record bulk manual attendance
// @Description Record attendance for multiple students at once
// @Tags Attendance
// @Accept json
// @Produce json
// @Param request body BulkManualAttendanceRequest true "Bulk attendance data"
// @Success 201 {object} []AttendanceResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/attendance/manual/bulk [post]
func (h *Handler) RecordBulkManualAttendance(c *fiber.Ctx) error {
	schoolID, ok := c.Locals("school_id").(uint)
	if !ok {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTHZ_TENANT_REQUIRED",
				"message": "School context is required",
			},
		})
	}

	var req BulkManualAttendanceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Invalid request body",
			},
		})
	}

	responses, err := h.service.RecordBulkManualAttendance(c.Context(), schoolID, req)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    responses,
		"message": "Bulk attendance recorded successfully",
	})
}


// GetAttendance handles listing attendance records
// @Summary List attendance records
// @Description Get a paginated list of attendance records with optional filters
// @Tags Attendance
// @Produce json
// @Param student_id query int false "Filter by student ID"
// @Param class_id query int false "Filter by class ID"
// @Param start_date query string false "Filter by start date (YYYY-MM-DD)"
// @Param end_date query string false "Filter by end date (YYYY-MM-DD)"
// @Param status query string false "Filter by status (on_time, late, very_late, absent)"
// @Param method query string false "Filter by method (rfid, manual)"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} AttendanceListResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/attendance [get]
func (h *Handler) GetAttendance(c *fiber.Ctx) error {
	schoolID, ok := c.Locals("school_id").(uint)
	if !ok {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTHZ_TENANT_REQUIRED",
				"message": "School context is required",
			},
		})
	}

	filter := DefaultAttendanceFilter()

	// Parse query parameters
	if studentIDStr := c.Query("student_id"); studentIDStr != "" {
		if studentID, err := strconv.ParseUint(studentIDStr, 10, 32); err == nil {
			id := uint(studentID)
			filter.StudentID = &id
		}
	}

	if classIDStr := c.Query("class_id"); classIDStr != "" {
		if classID, err := strconv.ParseUint(classIDStr, 10, 32); err == nil {
			id := uint(classID)
			filter.ClassID = &id
		}
	}

	if startDate := c.Query("start_date"); startDate != "" {
		filter.StartDate = &startDate
	}

	if endDate := c.Query("end_date"); endDate != "" {
		filter.EndDate = &endDate
	}

	if status := c.Query("status"); status != "" {
		filter.Status = &status
	}

	if method := c.Query("method"); method != "" {
		filter.Method = &method
	}

	if page, err := strconv.Atoi(c.Query("page", "1")); err == nil && page > 0 {
		filter.Page = page
	}

	if pageSize, err := strconv.Atoi(c.Query("page_size", "20")); err == nil && pageSize > 0 {
		filter.PageSize = pageSize
	}

	response, err := h.service.GetAllAttendance(c.Context(), schoolID, filter)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GetAttendanceByID handles getting a single attendance record
// @Summary Get attendance by ID
// @Description Get detailed information about a specific attendance record
// @Tags Attendance
// @Produce json
// @Param id path int true "Attendance ID"
// @Success 200 {object} AttendanceResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/attendance/{id} [get]
func (h *Handler) GetAttendanceByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Invalid attendance ID",
			},
		})
	}

	response, err := h.service.GetAttendanceByID(c.Context(), uint(id))
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GetStudentAttendance handles getting attendance for a specific student
// @Summary Get student attendance
// @Description Get attendance records for a specific student
// @Tags Attendance
// @Produce json
// @Param studentId path int true "Student ID"
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} []AttendanceResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/attendance/student/{studentId} [get]
func (h *Handler) GetStudentAttendance(c *fiber.Ctx) error {
	studentID, err := strconv.ParseUint(c.Params("studentId"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Invalid student ID",
			},
		})
	}

	var startDate, endDate time.Time

	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = parsed
		}
	}

	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endDate = parsed
		}
	}

	responses, err := h.service.GetStudentAttendance(c.Context(), uint(studentID), startDate, endDate)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    responses,
	})
}

// GetClassAttendance handles getting attendance for a specific class
// @Summary Get class attendance
// @Description Get attendance records for a specific class on a given date
// @Tags Attendance
// @Produce json
// @Param classId path int true "Class ID"
// @Param date query string false "Date (YYYY-MM-DD), defaults to today"
// @Success 200 {object} ClassAttendanceSummaryResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/attendance/class/{classId} [get]
func (h *Handler) GetClassAttendance(c *fiber.Ctx) error {
	classID, err := strconv.ParseUint(c.Params("classId"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Invalid class ID",
			},
		})
	}

	date := time.Now()
	if dateStr := c.Query("date"); dateStr != "" {
		if parsed, err := time.Parse("2006-01-02", dateStr); err == nil {
			date = parsed
		}
	}

	response, err := h.service.GetClassAttendance(c.Context(), uint(classID), date)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GetSchoolSummary handles getting school-wide attendance summary
// @Summary Get school attendance summary
// @Description Get attendance summary for the entire school on a given date
// @Tags Attendance
// @Produce json
// @Param date query string false "Date (YYYY-MM-DD), defaults to today"
// @Success 200 {object} SchoolAttendanceSummaryResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/attendance/summary [get]
func (h *Handler) GetSchoolSummary(c *fiber.Ctx) error {
	schoolID, ok := c.Locals("school_id").(uint)
	if !ok {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTHZ_TENANT_REQUIRED",
				"message": "School context is required",
			},
		})
	}

	date := time.Now()
	if dateStr := c.Query("date"); dateStr != "" {
		if parsed, err := time.Parse("2006-01-02", dateStr); err == nil {
			date = parsed
		}
	}

	response, err := h.service.GetSchoolAttendanceSummary(c.Context(), schoolID, date)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// DeleteAttendance handles deleting an attendance record
// @Summary Delete attendance record
// @Description Delete a specific attendance record
// @Tags Attendance
// @Produce json
// @Param id path int true "Attendance ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/attendance/{id} [delete]
func (h *Handler) DeleteAttendance(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Invalid attendance ID",
			},
		})
	}

	if err := h.service.DeleteAttendance(c.Context(), uint(id)); err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Attendance record deleted successfully",
	})
}

// handleError handles service errors and returns appropriate HTTP responses
func (h *Handler) handleError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, ErrAttendanceNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "NOT_FOUND_ATTENDANCE",
				"message": "Attendance record not found",
			},
		})
	case errors.Is(err, ErrStudentNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "NOT_FOUND_STUDENT",
				"message": "Student not found",
			},
		})
	case errors.Is(err, ErrInvalidRFIDCode):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_RFID",
				"message": "Invalid RFID code or student not found",
			},
		})
	case errors.Is(err, device.ErrInvalidAPIKey):
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTH_INVALID_API_KEY",
				"message": "Invalid API key",
			},
		})
	case errors.Is(err, ErrAPIKeyRequired):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "API key is required",
			},
		})
	case errors.Is(err, ErrRFIDCodeRequired):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "RFID code is required",
			},
		})
	case errors.Is(err, ErrStudentIDRequired):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "Student ID is required",
			},
		})
	case errors.Is(err, ErrDateRequired):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "Date is required",
			},
		})
	case errors.Is(err, ErrInvalidDate):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Invalid date format. Use YYYY-MM-DD",
			},
		})
	case errors.Is(err, ErrInvalidTime):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Invalid time format. Use HH:MM",
			},
		})
	case errors.Is(err, ErrCheckOutBeforeIn):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_TIME",
				"message": "Check-out time cannot be before check-in time",
			},
		})
	case errors.Is(err, ErrAlreadyCheckedOut):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_ALREADY_CHECKED_OUT",
				"message": "Student has already checked out today",
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
