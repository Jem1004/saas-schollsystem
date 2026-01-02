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
	repo    Repository
}

// NewHandler creates a new attendance handler
func NewHandler(service Service, repo Repository) *Handler {
	return &Handler{
		service: service,
		repo:    repo,
	}
}

// RegisterRoutes registers attendance routes for admin/teachers
func (h *Handler) RegisterRoutes(router fiber.Router) {
	attendance := router.Group("/attendance")
	attendance.Get("", h.GetAttendance)
	attendance.Get("/summary", h.GetSchoolSummary)
	attendance.Get("/export", h.ExportAttendance)
	attendance.Get("/monthly-recap", h.GetMonthlyRecap)
	attendance.Get("/monthly-recap/export", h.ExportMonthlyRecap)
	attendance.Get("/class/:classId", h.GetClassAttendance)
	attendance.Get("/student/:studentId", h.GetStudentAttendance)
	attendance.Get("/:id", h.GetAttendanceByID)
	attendance.Post("/manual", h.RecordManualAttendance)
	attendance.Post("/manual/bulk", h.RecordBulkManualAttendance)
	attendance.Delete("/:id", h.DeleteAttendance)
}

// RegisterRoutesWithoutGroup registers attendance routes without creating a sub-group
func (h *Handler) RegisterRoutesWithoutGroup(router fiber.Router) {
	router.Get("", h.GetAttendance)
	router.Get("/summary", h.GetSchoolSummary)
	router.Get("/export", h.ExportAttendance)
	router.Get("/monthly-recap", h.GetMonthlyRecap)
	router.Get("/monthly-recap/export", h.ExportMonthlyRecap)
	router.Get("/class/:classId", h.GetClassAttendance)
	router.Get("/student/:studentId", h.GetStudentAttendance)
	router.Get("/:id", h.GetAttendanceByID)
	router.Post("/manual", h.RecordManualAttendance)
	router.Post("/manual/bulk", h.RecordBulkManualAttendance)
	router.Delete("/:id", h.DeleteAttendance)
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
				"message": "Format data tidak valid",
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
				"message": "Konteks sekolah diperlukan",
			},
		})
	}

	var req ManualAttendanceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Format data tidak valid",
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
		"message": "Kehadiran berhasil dicatat",
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
				"message": "Konteks sekolah diperlukan",
			},
		})
	}

	var req BulkManualAttendanceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Format data tidak valid",
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
		"message": "Kehadiran massal berhasil dicatat",
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
				"message": "Konteks sekolah diperlukan",
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
				"message": "ID kehadiran tidak valid",
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
				"message": "ID siswa tidak valid",
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
				"message": "ID kelas tidak valid",
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
				"message": "Konteks sekolah diperlukan",
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
				"message": "ID kehadiran tidak valid",
			},
		})
	}

	if err := h.service.DeleteAttendance(c.Context(), uint(id)); err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Data kehadiran berhasil dihapus",
	})
}

// ExportAttendance handles exporting attendance data to Excel
// @Summary Export attendance to Excel
// @Description Export attendance records to Excel file with optional filters
// @Tags Attendance
// @Produce application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Param start_date query string true "Start date (YYYY-MM-DD)"
// @Param end_date query string true "End date (YYYY-MM-DD)"
// @Param class_id query int false "Filter by class ID"
// @Success 200 {file} file "Excel file"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/attendance/export [get]
func (h *Handler) ExportAttendance(c *fiber.Ctx) error {
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

	// Get school name from context (set by auth middleware)
	schoolName, _ := c.Locals("school_name").(string)
	if schoolName == "" {
		schoolName = "school"
	}

	// Parse filter parameters
	filter := ExportFilter{
		StartDate: c.Query("start_date"),
		EndDate:   c.Query("end_date"),
	}

	// Parse class_id if provided
	if classIDStr := c.Query("class_id"); classIDStr != "" {
		if classID, err := strconv.ParseUint(classIDStr, 10, 32); err == nil {
			id := uint(classID)
			filter.ClassID = &id
		}
	}

	// Requirements: 2.7 - Wali_Kelas only sees their assigned class
	role, _ := c.Locals("role").(string)
	if role == "wali_kelas" {
		userID, _ := c.Locals("userID").(uint)
		if userID > 0 {
			class, err := h.repo.FindClassByHomeroomTeacher(c.Context(), schoolID, userID)
			if err == nil && class != nil {
				filter.ClassID = &class.ID
			}
		}
	}

	// Generate Excel file
	excelData, filename, err := h.service.ExportAttendanceToExcel(c.Context(), schoolID, schoolName, filter)
	if err != nil {
		return h.handleError(c, err)
	}

	// Set response headers for file download
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", "attachment; filename="+filename)
	c.Set("Content-Length", strconv.Itoa(len(excelData)))

	return c.Send(excelData)
}

// GetMonthlyRecap handles getting monthly attendance recap
// @Summary Get monthly attendance recap
// @Description Get monthly attendance summary per student
// @Tags Attendance
// @Produce json
// @Param month query int true "Month (1-12)"
// @Param year query int true "Year (e.g., 2024)"
// @Param class_id query int false "Filter by class ID"
// @Success 200 {object} MonthlyRecapResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/attendance/monthly-recap [get]
func (h *Handler) GetMonthlyRecap(c *fiber.Ctx) error {
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

	// Parse filter parameters
	month, err := strconv.Atoi(c.Query("month"))
	if err != nil || month < 1 || month > 12 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Bulan harus antara 1-12",
			},
		})
	}

	year, err := strconv.Atoi(c.Query("year"))
	if err != nil || year < 2000 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Tahun tidak valid",
			},
		})
	}

	filter := MonthlyRecapFilter{
		Month: month,
		Year:  year,
	}

	// Parse class_id if provided
	if classIDStr := c.Query("class_id"); classIDStr != "" {
		if classID, err := strconv.ParseUint(classIDStr, 10, 32); err == nil {
			id := uint(classID)
			filter.ClassID = &id
		}
	}

	// Requirements: 2.7 - Wali_Kelas only sees their assigned class
	role, _ := c.Locals("role").(string)
	if role == "wali_kelas" {
		userID, _ := c.Locals("userID").(uint)
		if userID > 0 {
			class, err := h.repo.FindClassByHomeroomTeacher(c.Context(), schoolID, userID)
			if err == nil && class != nil {
				filter.ClassID = &class.ID
			}
		}
	}

	// Get monthly recap
	response, err := h.service.GetMonthlyRecap(c.Context(), schoolID, filter)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// ExportMonthlyRecap handles exporting monthly recap to Excel
// @Summary Export monthly recap to Excel
// @Description Export monthly attendance recap to Excel file
// @Tags Attendance
// @Produce application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Param month query int true "Month (1-12)"
// @Param year query int true "Year (e.g., 2024)"
// @Param class_id query int false "Filter by class ID"
// @Success 200 {file} file "Excel file"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/attendance/monthly-recap/export [get]
func (h *Handler) ExportMonthlyRecap(c *fiber.Ctx) error {
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

	// Get school name from context
	schoolName, _ := c.Locals("school_name").(string)
	if schoolName == "" {
		schoolName = "school"
	}

	// Parse filter parameters
	month, err := strconv.Atoi(c.Query("month"))
	if err != nil || month < 1 || month > 12 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Bulan harus antara 1-12",
			},
		})
	}

	year, err := strconv.Atoi(c.Query("year"))
	if err != nil || year < 2000 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Tahun tidak valid",
			},
		})
	}

	filter := MonthlyRecapFilter{
		Month: month,
		Year:  year,
	}

	// Parse class_id if provided
	if classIDStr := c.Query("class_id"); classIDStr != "" {
		if classID, err := strconv.ParseUint(classIDStr, 10, 32); err == nil {
			id := uint(classID)
			filter.ClassID = &id
		}
	}

	// Requirements: 2.7 - Wali_Kelas only sees their assigned class
	role, _ := c.Locals("role").(string)
	if role == "wali_kelas" {
		userID, _ := c.Locals("userID").(uint)
		if userID > 0 {
			class, err := h.repo.FindClassByHomeroomTeacher(c.Context(), schoolID, userID)
			if err == nil && class != nil {
				filter.ClassID = &class.ID
			}
		}
	}

	// Generate Excel file
	excelData, filename, err := h.service.ExportMonthlyRecapToExcel(c.Context(), schoolID, schoolName, filter)
	if err != nil {
		return h.handleError(c, err)
	}

	// Set response headers for file download
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", "attachment; filename="+filename)
	c.Set("Content-Length", strconv.Itoa(len(excelData)))

	return c.Send(excelData)
}

// handleError handles service errors and returns appropriate HTTP responses
func (h *Handler) handleError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, ErrAttendanceNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "NOT_FOUND_ATTENDANCE",
				"message": "Data kehadiran tidak ditemukan",
			},
		})
	case errors.Is(err, ErrStudentNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "NOT_FOUND_STUDENT",
				"message": "Siswa tidak ditemukan",
			},
		})
	case errors.Is(err, ErrInvalidRFIDCode):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_RFID",
				"message": "Kode RFID tidak valid atau siswa tidak ditemukan",
			},
		})
	case errors.Is(err, device.ErrInvalidAPIKey):
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTH_INVALID_API_KEY",
				"message": "API key tidak valid",
			},
		})
	case errors.Is(err, ErrAPIKeyRequired):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "API key wajib diisi",
			},
		})
	case errors.Is(err, ErrRFIDCodeRequired):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "Kode RFID wajib diisi",
			},
		})
	case errors.Is(err, ErrStudentIDRequired):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "ID siswa wajib diisi",
			},
		})
	case errors.Is(err, ErrDateRequired):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "Tanggal wajib diisi",
			},
		})
	case errors.Is(err, ErrInvalidDate):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Format tanggal tidak valid. Gunakan YYYY-MM-DD",
			},
		})
	case errors.Is(err, ErrInvalidTime):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Format waktu tidak valid. Gunakan HH:MM",
			},
		})
	case errors.Is(err, ErrCheckOutBeforeIn):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_TIME",
				"message": "Waktu check-out tidak boleh sebelum waktu check-in",
			},
		})
	case errors.Is(err, ErrAlreadyCheckedOut):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_ALREADY_CHECKED_OUT",
				"message": "Siswa sudah melakukan check-out hari ini",
			},
		})
	case errors.Is(err, ErrOutsideAttendanceWindow):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_NO_SCHEDULE",
				"message": "Tidak ada jadwal absensi untuk waktu ini",
			},
		})
	case errors.Is(err, ErrAlreadyCheckedIn):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_ALREADY_CHECKED_IN",
				"message": "Anda sudah absen untuk jadwal ini",
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
