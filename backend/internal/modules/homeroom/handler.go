package homeroom

import (
	"errors"
	"fmt"
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

// RegisterRoutesWithoutGroup registers Homeroom routes without creating a sub-group
func (h *Handler) RegisterRoutesWithoutGroup(router fiber.Router) {
	// Dashboard routes for Wali Kelas
	router.Get("/stats", h.GetHomeroomStats)
	router.Get("/my-class", h.GetMyClass)
	router.Get("/students", h.GetClassStudents)
	router.Get("/attendance", h.GetClassAttendance)
	router.Get("/schedules", h.GetActiveSchedules)

	// Manual attendance for Wali Kelas
	router.Post("/attendance/manual", h.RecordManualAttendance)
	router.Put("/attendance/:id", h.UpdateAttendance)

	// Grade CRUD for Wali Kelas
	router.Get("/grades", h.GetGrades)
	router.Post("/grades", h.CreateGrade)
	router.Post("/grades/batch", h.CreateBatchGrades)
	router.Get("/grades/:id", h.GetGradeByID)
	router.Put("/grades/:id", h.UpdateGrade)
	router.Delete("/grades/:id", h.DeleteGrade)

	// Student grades
	router.Get("/students/:studentId/grades", h.GetStudentGrades)

	// Note CRUD
	router.Get("/notes", h.GetNotes)
	router.Post("/notes", h.CreateNote)
	router.Get("/notes/:id", h.GetNoteByID)
	router.Put("/notes/:id", h.UpdateNote)
	router.Delete("/notes/:id", h.DeleteNote)

	// Student notes
	router.Get("/students/:studentId/notes", h.GetStudentNotes)
	router.Get("/students/:studentId/notes/summary", h.GetStudentNoteSummary)

	// Class notes
	router.Get("/class/:classId/notes", h.GetClassNotes)
}

// ==================== Dashboard Handlers ====================

// GetHomeroomStats handles getting dashboard statistics for wali kelas
// @Summary Get homeroom dashboard stats
// @Description Get dashboard statistics for wali kelas including attendance, grades, and notes
// @Tags Homeroom
// @Produce json
// @Success 200 {object} HomeroomStatsResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/homeroom/stats [get]
func (h *Handler) GetHomeroomStats(c *fiber.Ctx) error {
	schoolID, ok := c.Locals("school_id").(uint)
	if !ok {
		return h.tenantRequiredError(c)
	}

	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return h.authRequiredError(c)
	}

	response, err := h.service.GetHomeroomStats(c.Context(), schoolID, userID)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GetMyClass handles getting class information for wali kelas
// @Summary Get my class info
// @Description Get information about the class assigned to the wali kelas
// @Tags Homeroom
// @Produce json
// @Success 200 {object} ClassInfoResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/homeroom/my-class [get]
func (h *Handler) GetMyClass(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return h.authRequiredError(c)
	}

	response, err := h.service.GetMyClass(c.Context(), userID)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GetClassStudents handles getting students in wali kelas's class
// @Summary Get class students
// @Description Get list of students in the wali kelas's assigned class
// @Tags Homeroom
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(20)
// @Success 200 {object} ClassStudentListResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/homeroom/students [get]
func (h *Handler) GetClassStudents(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return h.authRequiredError(c)
	}

	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("pageSize", 20)

	response, err := h.service.GetClassStudents(c.Context(), userID, page, pageSize)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GetClassAttendance handles getting attendance for wali kelas's class
// @Summary Get class attendance
// @Description Get attendance records for the wali kelas's class on a specific date
// @Tags Homeroom
// @Produce json
// @Param date query string true "Date in YYYY-MM-DD format"
// @Success 200 {object} ClassAttendanceListResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/homeroom/attendance [get]
func (h *Handler) GetClassAttendance(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return h.authRequiredError(c)
	}

	date := c.Query("date")
	if date == "" {
		// Default to today
		date = c.Context().Time().Format("2006-01-02")
	}

	response, err := h.service.GetClassAttendance(c.Context(), userID, date)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GetActiveSchedules handles getting active schedules for the school
// @Summary Get active schedules
// @Description Get active attendance schedules for the school on a specific date
// @Tags Homeroom
// @Produce json
// @Param date query string false "Date in YYYY-MM-DD format (default: today)"
// @Success 200 {object} []ScheduleResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/homeroom/schedules [get]
func (h *Handler) GetActiveSchedules(c *fiber.Ctx) error {
	schoolID, ok := c.Locals("school_id").(uint)
	if !ok {
		return h.tenantRequiredError(c)
	}

	date := c.Query("date")
	if date == "" {
		date = c.Context().Time().Format("2006-01-02")
	}

	response, err := h.service.GetActiveSchedules(c.Context(), schoolID, date)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// RecordManualAttendance handles recording manual attendance
// @Summary Record manual attendance
// @Description Record manual attendance for a student in wali kelas's class
// @Tags Homeroom
// @Accept json
// @Produce json
// @Param request body ManualAttendanceRequest true "Attendance data"
// @Success 201 {object} StudentAttendanceResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/homeroom/attendance/manual [post]
func (h *Handler) RecordManualAttendance(c *fiber.Ctx) error {
	schoolID, ok := c.Locals("school_id").(uint)
	if !ok {
		return h.tenantRequiredError(c)
	}

	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return h.authRequiredError(c)
	}

	var req ManualAttendanceRequest
	if err := c.BodyParser(&req); err != nil {
		return h.invalidBodyError(c)
	}

	response, err := h.service.RecordManualAttendance(c.Context(), schoolID, userID, req)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    response,
		"message": "Absensi manual berhasil dicatat",
	})
}

// UpdateAttendance handles updating attendance record
// @Summary Update attendance
// @Description Update an existing attendance record
// @Tags Homeroom
// @Accept json
// @Produce json
// @Param id path int true "Attendance ID"
// @Param request body UpdateAttendanceRequest true "Updated attendance data"
// @Success 200 {object} StudentAttendanceResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/homeroom/attendance/{id} [put]
func (h *Handler) UpdateAttendance(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return h.authRequiredError(c)
	}

	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return h.invalidIDError(c, "attendance")
	}

	var req UpdateAttendanceRequest
	if err := c.BodyParser(&req); err != nil {
		return h.invalidBodyError(c)
	}

	response, err := h.service.UpdateAttendance(c.Context(), userID, uint(id), req)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
		"message": "Absensi berhasil diperbarui",
	})
}

// ==================== Grade Handlers ====================

// GetGrades handles listing grades for wali kelas's class
// @Summary List grades
// @Description Get a paginated list of grades for the wali kelas's class
// @Tags Homeroom
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(20)
// @Success 200 {object} GradeListResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/homeroom/grades [get]
func (h *Handler) GetGrades(c *fiber.Ctx) error {
	schoolID, ok := c.Locals("school_id").(uint)
	if !ok {
		return h.tenantRequiredError(c)
	}

	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return h.authRequiredError(c)
	}

	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("pageSize", 20)

	response, err := h.service.GetClassGrades(c.Context(), schoolID, userID, page, pageSize)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// CreateGrade handles creating a new grade
// @Summary Create grade
// @Description Create a new grade for a student in wali kelas's class
// @Tags Homeroom
// @Accept json
// @Produce json
// @Param request body CreateGradeRequest true "Grade data"
// @Success 201 {object} GradeResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/homeroom/grades [post]
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
		"message": "Nilai berhasil ditambahkan",
	})
}

// CreateBatchGrades handles creating multiple grades at once
// @Summary Create batch grades
// @Description Create multiple grades for students in wali kelas's class
// @Tags Homeroom
// @Accept json
// @Produce json
// @Param request body BatchGradeRequest true "Batch grade data"
// @Success 201 {object} []GradeResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/homeroom/grades/batch [post]
func (h *Handler) CreateBatchGrades(c *fiber.Ctx) error {
	schoolID, ok := c.Locals("school_id").(uint)
	if !ok {
		return h.tenantRequiredError(c)
	}

	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return h.authRequiredError(c)
	}

	var req BatchGradeRequest
	if err := c.BodyParser(&req); err != nil {
		return h.invalidBodyError(c)
	}

	responses, err := h.service.CreateBatchGrades(c.Context(), schoolID, userID, req)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    responses,
		"message": "Nilai berhasil ditambahkan",
	})
}

// GetGradeByID handles getting a single grade
// @Summary Get grade by ID
// @Description Get detailed information about a specific grade
// @Tags Homeroom
// @Produce json
// @Param id path int true "Grade ID"
// @Success 200 {object} GradeResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/homeroom/grades/{id} [get]
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

// UpdateGrade handles updating a grade
// @Summary Update grade
// @Description Update an existing grade
// @Tags Homeroom
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
// @Router /api/v1/homeroom/grades/{id} [put]
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
// @Tags Homeroom
// @Produce json
// @Param id path int true "Grade ID"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/homeroom/grades/{id} [delete]
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

// GetStudentGrades handles getting grades for a specific student
// @Summary Get student grades
// @Description Get all grades for a specific student in wali kelas's class
// @Tags Homeroom
// @Produce json
// @Param studentId path int true "Student ID"
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(50)
// @Success 200 {object} GradeListResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/homeroom/students/{studentId}/grades [get]
func (h *Handler) GetStudentGrades(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return h.authRequiredError(c)
	}

	studentID, err := strconv.ParseUint(c.Params("studentId"), 10, 32)
	if err != nil {
		return h.invalidIDError(c, "student")
	}

	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("pageSize", 50)

	response, err := h.service.GetStudentGrades(c.Context(), userID, uint(studentID), page, pageSize)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
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
	// Log raw body first
	fmt.Printf("[DEBUG] CreateNote - Raw body: %s\n", string(c.Body()))
	fmt.Printf("[DEBUG] CreateNote - Content-Type: %s\n", c.Get("Content-Type"))

	schoolID, ok := c.Locals("school_id").(uint)
	if !ok {
		fmt.Printf("[DEBUG] CreateNote - school_id not found in locals\n")
		return h.tenantRequiredError(c)
	}
	fmt.Printf("[DEBUG] CreateNote - schoolID: %d\n", schoolID)

	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		fmt.Printf("[DEBUG] CreateNote - user_id not found in locals\n")
		return h.authRequiredError(c)
	}
	fmt.Printf("[DEBUG] CreateNote - userID: %d\n", userID)

	var req CreateNoteRequest
	if err := c.BodyParser(&req); err != nil {
		// Log the raw body for debugging
		fmt.Printf("[DEBUG] CreateNote - Body parse error: %v\n", err)
		return h.invalidBodyError(c)
	}

	// Log the parsed request
	fmt.Printf("[DEBUG] CreateNote - Parsed request: StudentID=%d, Content=%s\n", req.StudentID, req.Content)

	response, err := h.service.CreateNote(c.Context(), schoolID, userID, req)
	if err != nil {
		fmt.Printf("[DEBUG] CreateNote - Service error: %v\n", err)
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
// @Router /api/v1/homeroom/students/{studentId}/notes [get]
func (h *Handler) GetStudentNotes(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return h.authRequiredError(c)
	}

	studentID, err := strconv.ParseUint(c.Params("studentId"), 10, 32)
	if err != nil {
		return h.invalidIDError(c, "student")
	}

	// Validate teacher has access to this student
	if err := h.service.ValidateTeacherAccess(c.Context(), userID, uint(studentID)); err != nil {
		return h.handleError(c, err)
	}

	responses, err := h.service.GetStudentNotes(c.Context(), uint(studentID))
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"notes": responses,
		},
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
	case errors.Is(err, ErrAttendanceNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "NOT_FOUND_ATTENDANCE",
				"message": "Data absensi tidak ditemukan",
			},
		})
	case errors.Is(err, ErrAttendanceAlreadyExists):
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "CONFLICT_ATTENDANCE",
				"message": "Data absensi untuk siswa ini pada tanggal tersebut sudah ada",
			},
		})
	case errors.Is(err, ErrInvalidStatus):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_STATUS",
				"message": "Status absensi tidak valid",
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
