package student

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/school-management/backend/internal/middleware"
)

// Handler handles HTTP requests for student API
type Handler struct {
	service Service
}

// NewHandler creates a new student handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers student routes (Student Mobile App)
// Requirements: 16.1-16.5 - Student self-monitoring
func (h *Handler) RegisterRoutes(router fiber.Router) {
	student := router.Group("/student")

	// Profile
	student.Get("/profile", h.GetProfile)
	student.Get("/dashboard", h.GetDashboard)
	student.Get("/summary", h.GetSummary)

	// Attendance
	student.Get("/attendance", h.GetAttendance)
	student.Get("/attendance/summary", h.GetAttendanceSummary)

	// Grades
	student.Get("/grades", h.GetGrades)
	student.Get("/grades/summary", h.GetGradeSummary)

	// BK information
	student.Get("/bk", h.GetBKInfo)
}

// GetProfile handles getting the student's profile
// @Summary Get student profile
// @Description Get the authenticated student's profile
// @Tags Student
// @Produce json
// @Success 200 {object} StudentProfileResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/student/profile [get]
func (h *Handler) GetProfile(c *fiber.Ctx) error {
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

	response, err := h.service.GetProfile(c.Context(), userID)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GetDashboard handles getting the student's dashboard
// @Summary Get student dashboard
// @Description Get the authenticated student's dashboard with summary and recent data
// @Tags Student
// @Produce json
// @Success 200 {object} DashboardResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/student/dashboard [get]
func (h *Handler) GetDashboard(c *fiber.Ctx) error {
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

	response, err := h.service.GetDashboard(c.Context(), userID)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GetSummary handles getting the student's summary statistics
// @Summary Get student summary
// @Description Get the authenticated student's summary statistics
// @Tags Student
// @Produce json
// @Success 200 {object} StudentSummaryResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/student/summary [get]
func (h *Handler) GetSummary(c *fiber.Ctx) error {
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

	response, err := h.service.GetSummary(c.Context(), userID)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GetAttendance handles getting the student's attendance records
// @Summary Get student attendance
// @Description Get the authenticated student's attendance records
// @Tags Student
// @Produce json
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} AttendanceListResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/student/attendance [get]
func (h *Handler) GetAttendance(c *fiber.Ctx) error {
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

	filter := DefaultAttendanceFilter()
	filter.StartDate = c.Query("start_date")
	filter.EndDate = c.Query("end_date")
	if page, err := strconv.Atoi(c.Query("page", "1")); err == nil && page > 0 {
		filter.Page = page
	}
	if pageSize, err := strconv.Atoi(c.Query("page_size", "20")); err == nil && pageSize > 0 {
		filter.PageSize = pageSize
	}

	response, err := h.service.GetAttendance(c.Context(), userID, filter)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GetAttendanceSummary handles getting the student's attendance summary
// @Summary Get student attendance summary
// @Description Get the authenticated student's attendance summary
// @Tags Student
// @Produce json
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} AttendanceSummaryResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/student/attendance/summary [get]
func (h *Handler) GetAttendanceSummary(c *fiber.Ctx) error {
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

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	response, err := h.service.GetAttendanceSummary(c.Context(), userID, startDate, endDate)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GetGrades handles getting the student's grades
// @Summary Get student grades
// @Description Get the authenticated student's grades
// @Tags Student
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} GradeListResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/student/grades [get]
func (h *Handler) GetGrades(c *fiber.Ctx) error {
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

	filter := DefaultGradeFilter()
	if page, err := strconv.Atoi(c.Query("page", "1")); err == nil && page > 0 {
		filter.Page = page
	}
	if pageSize, err := strconv.Atoi(c.Query("page_size", "20")); err == nil && pageSize > 0 {
		filter.PageSize = pageSize
	}

	response, err := h.service.GetGrades(c.Context(), userID, filter)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GetGradeSummary handles getting the student's grade summary
// @Summary Get student grade summary
// @Description Get the authenticated student's grade summary
// @Tags Student
// @Produce json
// @Success 200 {object} GradeSummaryResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/student/grades/summary [get]
func (h *Handler) GetGradeSummary(c *fiber.Ctx) error {
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

	response, err := h.service.GetGradeSummary(c.Context(), userID)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GetBKInfo handles getting the student's BK information
// @Summary Get student BK information
// @Description Get the authenticated student's BK information (achievements and violations summary)
// @Tags Student
// @Produce json
// @Success 200 {object} BKInfoResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/student/bk [get]
func (h *Handler) GetBKInfo(c *fiber.Ctx) error {
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

	response, err := h.service.GetBKInfo(c.Context(), userID)
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
	case errors.Is(err, ErrStudentNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "NOT_FOUND_STUDENT",
				"message": "Student not found",
			},
		})
	case errors.Is(err, ErrUserNotStudent):
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTHZ_NOT_STUDENT",
				"message": "User is not a student",
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
