package grade

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Handler handles HTTP requests for Grade management
type Handler struct {
	service Service
}

// NewHandler creates a new Grade handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers Grade routes for Wali Kelas
// Requirements: 10.1, 10.2, 10.4, 10.5
func (h *Handler) RegisterRoutes(router fiber.Router) {
	grades := router.Group("/grades")

	// Grade CRUD
	grades.Get("", h.GetGrades)
	grades.Post("", h.CreateGrade)
	grades.Get("/:id", h.GetGradeByID)
	grades.Put("/:id", h.UpdateGrade)
	grades.Delete("/:id", h.DeleteGrade)

	// Student grades
	grades.Get("/student/:studentId", h.GetStudentGrades)
	grades.Get("/student/:studentId/summary", h.GetStudentGradeSummary)

	// Class grades
	grades.Get("/class/:classId", h.GetClassGrades)
}

// ==================== Grade Handlers ====================

// CreateGrade handles creating a new grade
// @Summary Create grade
// @Description Input a new grade for a student
// @Tags Grades
// @Accept json
// @Produce json
// @Param request body CreateGradeRequest true "Grade data"
// @Success 201 {object} GradeResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/grades [post]
func (h *Handler) CreateGrade(c *fiber.Ctx) error {
	schoolID, ok := c.Locals("school_id").(uint)
	if !ok {
		return h.tenantRequiredError(c)
	}

	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return h.authRequiredError(c)
	}

	var req CreateGradeRequest
	if err := c.BodyParser(&req); err != nil {
		return h.invalidBodyError(c)
	}

	response, err := h.service.CreateGrade(c.Context(), schoolID, userID, req)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    response,
		"message": "Nilai berhasil dicatat",
	})
}

// GetGrades handles listing grades
// @Summary List grades
// @Description Get a paginated list of grades
// @Tags Grades
// @Produce json
// @Param student_id query int false "Filter by student ID"
// @Param class_id query int false "Filter by class ID"
// @Param start_date query string false "Filter by start date (YYYY-MM-DD)"
// @Param end_date query string false "Filter by end date (YYYY-MM-DD)"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} GradeListResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/grades [get]
func (h *Handler) GetGrades(c *fiber.Ctx) error {
	schoolID, ok := c.Locals("school_id").(uint)
	if !ok {
		return h.tenantRequiredError(c)
	}

	filter := h.parseGradeFilter(c)

	response, err := h.service.GetGrades(c.Context(), schoolID, filter)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}


// GetGradeByID handles getting a single grade
// @Summary Get grade by ID
// @Description Get detailed information about a specific grade
// @Tags Grades
// @Produce json
// @Param id path int true "Grade ID"
// @Success 200 {object} GradeResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/grades/{id} [get]
func (h *Handler) GetGradeByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return h.invalidIDError(c, "grade")
	}

	response, err := h.service.GetGradeByID(c.Context(), uint(id))
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GetStudentGrades handles getting grades for a specific student
// @Summary Get student grades
// @Description Get all grades for a specific student sorted by date
// @Tags Grades
// @Produce json
// @Param studentId path int true "Student ID"
// @Success 200 {object} []GradeResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/grades/student/{studentId} [get]
func (h *Handler) GetStudentGrades(c *fiber.Ctx) error {
	studentID, err := strconv.ParseUint(c.Params("studentId"), 10, 32)
	if err != nil {
		return h.invalidIDError(c, "student")
	}

	responses, err := h.service.GetStudentGrades(c.Context(), uint(studentID))
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    responses,
	})
}

// GetStudentGradeSummary handles getting grade summary for a student
// @Summary Get student grade summary
// @Description Get grade summary (total grades, average score) for a student
// @Tags Grades
// @Produce json
// @Param studentId path int true "Student ID"
// @Success 200 {object} StudentGradeSummary
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/grades/student/{studentId}/summary [get]
func (h *Handler) GetStudentGradeSummary(c *fiber.Ctx) error {
	studentID, err := strconv.ParseUint(c.Params("studentId"), 10, 32)
	if err != nil {
		return h.invalidIDError(c, "student")
	}

	response, err := h.service.GetStudentGradeSummary(c.Context(), uint(studentID))
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GetClassGrades handles getting grades for all students in a class
// @Summary Get class grades
// @Description Get all grades for students in a specific class
// @Tags Grades
// @Produce json
// @Param classId path int true "Class ID"
// @Param start_date query string false "Filter by start date (YYYY-MM-DD)"
// @Param end_date query string false "Filter by end date (YYYY-MM-DD)"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} GradeListResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/grades/class/{classId} [get]
func (h *Handler) GetClassGrades(c *fiber.Ctx) error {
	classID, err := strconv.ParseUint(c.Params("classId"), 10, 32)
	if err != nil {
		return h.invalidIDError(c, "class")
	}

	filter := h.parseGradeFilter(c)

	response, err := h.service.GetClassGrades(c.Context(), uint(classID), filter)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// UpdateGrade handles updating a grade
// @Summary Update grade
// @Description Update an existing grade
// @Tags Grades
// @Accept json
// @Produce json
// @Param id path int true "Grade ID"
// @Param request body UpdateGradeRequest true "Updated grade data"
// @Success 200 {object} GradeResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/grades/{id} [put]
func (h *Handler) UpdateGrade(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return h.invalidIDError(c, "grade")
	}

	var req UpdateGradeRequest
	if err := c.BodyParser(&req); err != nil {
		return h.invalidBodyError(c)
	}

	response, err := h.service.UpdateGrade(c.Context(), uint(id), req)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
		"message": "Nilai berhasil diperbarui",
	})
}

// DeleteGrade handles deleting a grade
// @Summary Delete grade
// @Description Delete a specific grade record
// @Tags Grades
// @Produce json
// @Param id path int true "Grade ID"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/grades/{id} [delete]
func (h *Handler) DeleteGrade(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return h.invalidIDError(c, "grade")
	}

	if err := h.service.DeleteGrade(c.Context(), uint(id)); err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Nilai berhasil dihapus",
	})
}

// ==================== Helper Functions ====================

func (h *Handler) parseGradeFilter(c *fiber.Ctx) GradeFilter {
	filter := GradeFilter{
		Page:     1,
		PageSize: 20,
	}

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

	if page, err := strconv.Atoi(c.Query("page", "1")); err == nil && page > 0 {
		filter.Page = page
	}

	if pageSize, err := strconv.Atoi(c.Query("page_size", "20")); err == nil && pageSize > 0 {
		filter.PageSize = pageSize
	}

	return filter
}

// ==================== Error Handlers ====================

func (h *Handler) tenantRequiredError(c *fiber.Ctx) error {
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
		"success": false,
		"error": fiber.Map{
			"code":    "AUTHZ_TENANT_REQUIRED",
			"message": "Konteks sekolah diperlukan",
		},
	})
}

func (h *Handler) authRequiredError(c *fiber.Ctx) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"success": false,
		"error": fiber.Map{
			"code":    "AUTH_REQUIRED",
			"message": "Autentikasi diperlukan",
		},
	})
}

func (h *Handler) invalidBodyError(c *fiber.Ctx) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"success": false,
		"error": fiber.Map{
			"code":    "VAL_INVALID_FORMAT",
			"message": "Format data tidak valid",
		},
	})
}

func (h *Handler) invalidIDError(c *fiber.Ctx, resource string) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"success": false,
		"error": fiber.Map{
			"code":    "VAL_INVALID_FORMAT",
			"message": "Invalid " + resource + " ID",
		},
	})
}

func (h *Handler) handleError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, ErrGradeNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "NOT_FOUND_GRADE",
				"message": "Nilai tidak ditemukan",
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
	case errors.Is(err, ErrStudentIDRequired):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "ID siswa wajib diisi",
			},
		})
	case errors.Is(err, ErrTitleRequired):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "Judul wajib diisi",
			},
		})
	case errors.Is(err, ErrScoreInvalid):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_VALUE",
				"message": "Nilai harus antara 0 dan 100",
			},
		})
	case errors.Is(err, ErrStudentNotInSchool):
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTHZ_TENANT_MISMATCH",
				"message": "Siswa bukan dari sekolah ini",
			},
		})
	case errors.Is(err, ErrStudentNotInClass):
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTHZ_CLASS_MISMATCH",
				"message": "Siswa bukan dari kelas yang Anda ampu",
			},
		})
	case errors.Is(err, ErrNoClassAssigned):
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTHZ_NO_CLASS",
				"message": "Tidak ada kelas yang ditugaskan untuk guru ini",
			},
		})
	case errors.Is(err, ErrNotAuthorized):
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTHZ_NOT_AUTHORIZED",
				"message": "Tidak memiliki izin untuk melakukan aksi ini",
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
