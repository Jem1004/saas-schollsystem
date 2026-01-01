package homeroom

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Handler handles HTTP requests for Homeroom Note management
type Handler struct {
	service Service
}

// NewHandler creates a new Homeroom handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers Homeroom routes for Wali Kelas
// Requirements: 11.1, 11.3, 11.4, 11.5
func (h *Handler) RegisterRoutes(router fiber.Router) {
	notes := router.Group("/homeroom")

	// Note CRUD
	notes.Get("", h.GetNotes)
	notes.Post("", h.CreateNote)
	notes.Get("/:id", h.GetNoteByID)
	notes.Put("/:id", h.UpdateNote)
	notes.Delete("/:id", h.DeleteNote)

	// Student notes
	notes.Get("/student/:studentId", h.GetStudentNotes)
	notes.Get("/student/:studentId/summary", h.GetStudentNoteSummary)

	// Class notes
	notes.Get("/class/:classId", h.GetClassNotes)
}

// ==================== Note Handlers ====================

// CreateNote handles creating a new homeroom note
// @Summary Create homeroom note
// @Description Create a new homeroom note for a student
// @Tags Homeroom
// @Accept json
// @Produce json
// @Param request body CreateNoteRequest true "Note data"
// @Success 201 {object} NoteResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/homeroom [post]
func (h *Handler) CreateNote(c *fiber.Ctx) error {
	schoolID, ok := c.Locals("school_id").(uint)
	if !ok {
		return h.tenantRequiredError(c)
	}

	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return h.authRequiredError(c)
	}

	var req CreateNoteRequest
	if err := c.BodyParser(&req); err != nil {
		return h.invalidBodyError(c)
	}

	response, err := h.service.CreateNote(c.Context(), schoolID, userID, req)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    response,
		"message": "Catatan wali kelas berhasil dibuat",
	})
}

// GetNotes handles listing homeroom notes
// @Summary List homeroom notes
// @Description Get a paginated list of homeroom notes
// @Tags Homeroom
// @Produce json
// @Param student_id query int false "Filter by student ID"
// @Param class_id query int false "Filter by class ID"
// @Param teacher_id query int false "Filter by teacher ID"
// @Param start_date query string false "Filter by start date (YYYY-MM-DD)"
// @Param end_date query string false "Filter by end date (YYYY-MM-DD)"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} NoteListResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/homeroom [get]
func (h *Handler) GetNotes(c *fiber.Ctx) error {
	schoolID, ok := c.Locals("school_id").(uint)
	if !ok {
		return h.tenantRequiredError(c)
	}

	filter := h.parseNoteFilter(c)

	response, err := h.service.GetNotes(c.Context(), schoolID, filter)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}


// GetNoteByID handles getting a single homeroom note
// @Summary Get homeroom note by ID
// @Description Get detailed information about a specific homeroom note
// @Tags Homeroom
// @Produce json
// @Param id path int true "Note ID"
// @Success 200 {object} NoteResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/homeroom/{id} [get]
func (h *Handler) GetNoteByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return h.invalidIDError(c, "note")
	}

	response, err := h.service.GetNoteByID(c.Context(), uint(id))
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GetStudentNotes handles getting homeroom notes for a specific student
// @Summary Get student homeroom notes
// @Description Get all homeroom notes for a specific student
// @Tags Homeroom
// @Produce json
// @Param studentId path int true "Student ID"
// @Success 200 {object} []NoteResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/homeroom/student/{studentId} [get]
func (h *Handler) GetStudentNotes(c *fiber.Ctx) error {
	studentID, err := strconv.ParseUint(c.Params("studentId"), 10, 32)
	if err != nil {
		return h.invalidIDError(c, "student")
	}

	responses, err := h.service.GetStudentNotes(c.Context(), uint(studentID))
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    responses,
	})
}

// GetStudentNoteSummary handles getting note summary for a student
// @Summary Get student note summary
// @Description Get note summary (total notes, last note date) for a student
// @Tags Homeroom
// @Produce json
// @Param studentId path int true "Student ID"
// @Success 200 {object} StudentNoteSummary
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/homeroom/student/{studentId}/summary [get]
func (h *Handler) GetStudentNoteSummary(c *fiber.Ctx) error {
	studentID, err := strconv.ParseUint(c.Params("studentId"), 10, 32)
	if err != nil {
		return h.invalidIDError(c, "student")
	}

	response, err := h.service.GetStudentNoteSummary(c.Context(), uint(studentID))
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GetClassNotes handles getting homeroom notes for all students in a class
// @Summary Get class homeroom notes
// @Description Get all homeroom notes for students in a specific class
// @Tags Homeroom
// @Produce json
// @Param classId path int true "Class ID"
// @Param start_date query string false "Filter by start date (YYYY-MM-DD)"
// @Param end_date query string false "Filter by end date (YYYY-MM-DD)"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} NoteListResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/homeroom/class/{classId} [get]
func (h *Handler) GetClassNotes(c *fiber.Ctx) error {
	classID, err := strconv.ParseUint(c.Params("classId"), 10, 32)
	if err != nil {
		return h.invalidIDError(c, "class")
	}

	filter := h.parseNoteFilter(c)

	response, err := h.service.GetClassNotes(c.Context(), uint(classID), filter)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// UpdateNote handles updating a homeroom note
// @Summary Update homeroom note
// @Description Update an existing homeroom note
// @Tags Homeroom
// @Accept json
// @Produce json
// @Param id path int true "Note ID"
// @Param request body UpdateNoteRequest true "Updated note data"
// @Success 200 {object} NoteResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/homeroom/{id} [put]
func (h *Handler) UpdateNote(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return h.invalidIDError(c, "note")
	}

	var req UpdateNoteRequest
	if err := c.BodyParser(&req); err != nil {
		return h.invalidBodyError(c)
	}

	response, err := h.service.UpdateNote(c.Context(), uint(id), req)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
		"message": "Catatan wali kelas berhasil diperbarui",
	})
}

// DeleteNote handles deleting a homeroom note
// @Summary Delete homeroom note
// @Description Delete a specific homeroom note
// @Tags Homeroom
// @Produce json
// @Param id path int true "Note ID"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/homeroom/{id} [delete]
func (h *Handler) DeleteNote(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return h.invalidIDError(c, "note")
	}

	if err := h.service.DeleteNote(c.Context(), uint(id)); err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Catatan wali kelas berhasil dihapus",
	})
}

// ==================== Helper Functions ====================

func (h *Handler) parseNoteFilter(c *fiber.Ctx) NoteFilter {
	filter := NoteFilter{
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

	if teacherIDStr := c.Query("teacher_id"); teacherIDStr != "" {
		if teacherID, err := strconv.ParseUint(teacherIDStr, 10, 32); err == nil {
			id := uint(teacherID)
			filter.TeacherID = &id
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
	case errors.Is(err, ErrNoteNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "NOT_FOUND_NOTE",
				"message": "Catatan wali kelas tidak ditemukan",
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
	case errors.Is(err, ErrContentRequired):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "Konten wajib diisi",
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
