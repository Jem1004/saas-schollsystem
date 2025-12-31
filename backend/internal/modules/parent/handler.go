package parent

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/school-management/backend/internal/middleware"
)

// Handler handles HTTP requests for parent API
type Handler struct {
	service Service
}

// NewHandler creates a new parent handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers parent routes (Parent Mobile App)
// Requirements: 12.2, 14.4, 15.1, 15.2 - Parent data access
func (h *Handler) RegisterRoutes(router fiber.Router) {
	parent := router.Group("/parent")

	// Children management
	parent.Get("/children", h.GetChildren)
	parent.Get("/children/:id/dashboard", h.GetChildDashboard)

	// Attendance
	parent.Get("/children/:id/attendance", h.GetChildAttendance)
	parent.Get("/children/:id/attendance/summary", h.GetChildAttendanceSummary)

	// Grades
	parent.Get("/children/:id/grades", h.GetChildGrades)
	parent.Get("/children/:id/grades/summary", h.GetChildGradeSummary)

	// Homeroom notes
	parent.Get("/children/:id/notes", h.GetChildNotes)

	// BK information
	parent.Get("/children/:id/bk", h.GetChildBKInfo)
}

// GetChildren handles getting all linked children
// @Summary Get linked children
// @Description Get all children linked to the authenticated parent
// @Tags Parent
// @Produce json
// @Success 200 {object} ChildListResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/parent/children [get]
func (h *Handler) GetChildren(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTH_REQUIRED",
				"message": "Authentication required",
			},
		})
	}

	response, err := h.service.GetLinkedChildren(c.Context(), userID)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GetChildDashboard handles getting dashboard data for a child
// @Summary Get child dashboard
// @Description Get dashboard data for a specific child
// @Tags Parent
// @Produce json
// @Param id path int true "Student ID"
// @Success 200 {object} ChildDashboardResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/parent/children/{id}/dashboard [get]
func (h *Handler) GetChildDashboard(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTH_REQUIRED",
				"message": "Authentication required",
			},
		})
	}

	studentID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Invalid student ID",
			},
		})
	}

	response, err := h.service.GetChildDashboard(c.Context(), userID, uint(studentID))
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GetChildAttendance handles getting attendance records for a child
// @Summary Get child attendance
// @Description Get attendance records for a specific child
// @Tags Parent
// @Produce json
// @Param id path int true "Student ID"
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} ChildAttendanceListResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/parent/children/{id}/attendance [get]
func (h *Handler) GetChildAttendance(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTH_REQUIRED",
				"message": "Authentication required",
			},
		})
	}

	studentID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Invalid student ID",
			},
		})
	}

	filter := DefaultAttendanceFilter()
	filter.StartDate = c.Query("start_date")
	filter.EndDate = c.Query("end_date")
	if page, err := strconv.Atoi(c.Query("page", "1")); err == nil && page > 0 {
		filter.Page = page
	}
	if pageSize, err := strconv.Atoi(c.Query("page_size", "20")); err == nil && pageSize > 0 {
		filter.PageSize = pageSize
	}

	response, err := h.service.GetChildAttendance(c.Context(), userID, uint(studentID), filter)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GetChildAttendanceSummary handles getting attendance summary for a child
// @Summary Get child attendance summary
// @Description Get attendance summary for a specific child
// @Tags Parent
// @Produce json
// @Param id path int true "Student ID"
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} AttendanceSummaryResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/parent/children/{id}/attendance/summary [get]
func (h *Handler) GetChildAttendanceSummary(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTH_REQUIRED",
				"message": "Authentication required",
			},
		})
	}

	studentID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Invalid student ID",
			},
		})
	}

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	response, err := h.service.GetChildAttendanceSummary(c.Context(), userID, uint(studentID), startDate, endDate)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GetChildGrades handles getting grades for a child
// @Summary Get child grades
// @Description Get grades for a specific child
// @Tags Parent
// @Produce json
// @Param id path int true "Student ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} ChildGradeListResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/parent/children/{id}/grades [get]
func (h *Handler) GetChildGrades(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTH_REQUIRED",
				"message": "Authentication required",
			},
		})
	}

	studentID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Invalid student ID",
			},
		})
	}

	filter := DefaultGradeFilter()
	if page, err := strconv.Atoi(c.Query("page", "1")); err == nil && page > 0 {
		filter.Page = page
	}
	if pageSize, err := strconv.Atoi(c.Query("page_size", "20")); err == nil && pageSize > 0 {
		filter.PageSize = pageSize
	}

	response, err := h.service.GetChildGrades(c.Context(), userID, uint(studentID), filter)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GetChildGradeSummary handles getting grade summary for a child
// @Summary Get child grade summary
// @Description Get grade summary for a specific child
// @Tags Parent
// @Produce json
// @Param id path int true "Student ID"
// @Success 200 {object} GradeSummaryResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/parent/children/{id}/grades/summary [get]
func (h *Handler) GetChildGradeSummary(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTH_REQUIRED",
				"message": "Authentication required",
			},
		})
	}

	studentID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Invalid student ID",
			},
		})
	}

	response, err := h.service.GetChildGradeSummary(c.Context(), userID, uint(studentID))
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GetChildNotes handles getting homeroom notes for a child
// @Summary Get child homeroom notes
// @Description Get homeroom notes for a specific child
// @Tags Parent
// @Produce json
// @Param id path int true "Student ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} ChildNoteListResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/parent/children/{id}/notes [get]
func (h *Handler) GetChildNotes(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTH_REQUIRED",
				"message": "Authentication required",
			},
		})
	}

	studentID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Invalid student ID",
			},
		})
	}

	filter := DefaultNoteFilter()
	if page, err := strconv.Atoi(c.Query("page", "1")); err == nil && page > 0 {
		filter.Page = page
	}
	if pageSize, err := strconv.Atoi(c.Query("page_size", "20")); err == nil && pageSize > 0 {
		filter.PageSize = pageSize
	}

	response, err := h.service.GetChildNotes(c.Context(), userID, uint(studentID), filter)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GetChildBKInfo handles getting BK information for a child
// @Summary Get child BK information
// @Description Get BK information (violations, achievements, permits, counseling) for a specific child
// @Tags Parent
// @Produce json
// @Param id path int true "Student ID"
// @Success 200 {object} ChildBKInfoResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/parent/children/{id}/bk [get]
func (h *Handler) GetChildBKInfo(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTH_REQUIRED",
				"message": "Authentication required",
			},
		})
	}

	studentID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Invalid student ID",
			},
		})
	}

	response, err := h.service.GetChildBKInfo(c.Context(), userID, uint(studentID))
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
	case errors.Is(err, ErrParentNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "NOT_FOUND_PARENT",
				"message": "Parent not found",
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
	case errors.Is(err, ErrNotLinked):
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTHZ_NOT_LINKED",
				"message": "You do not have access to this student's data",
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
