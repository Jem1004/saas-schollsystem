package bk

import (
	"errors"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Handler handles HTTP requests for BK (Bimbingan Konseling) management
type Handler struct {
	service Service
}

// NewHandler creates a new BK handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers BK routes for Guru BK (full access)
func (h *Handler) RegisterRoutes(router fiber.Router) {
	bk := router.Group("/bk")

	// Dashboard
	bk.Get("/dashboard", h.GetDashboard)

	// Violations
	violations := bk.Group("/violations")
	violations.Get("", h.GetViolations)
	violations.Post("", h.CreateViolation)
	violations.Get("/:id", h.GetViolationByID)
	violations.Delete("/:id", h.DeleteViolation)

	// Achievements
	achievements := bk.Group("/achievements")
	achievements.Get("", h.GetAchievements)
	achievements.Post("", h.CreateAchievement)
	achievements.Get("/:id", h.GetAchievementByID)
	achievements.Delete("/:id", h.DeleteAchievement)
	achievements.Get("/student/:studentId/points", h.GetStudentAchievementPoints)

	// Permits
	permits := bk.Group("/permits")
	permits.Get("", h.GetPermits)
	permits.Post("", h.CreatePermit)
	permits.Get("/:id", h.GetPermitByID)
	permits.Post("/:id/return", h.RecordReturn)
	permits.Get("/:id/document", h.GetPermitDocument)
	permits.Delete("/:id", h.DeletePermit)

	// Counseling Notes
	counseling := bk.Group("/counseling")
	counseling.Get("", h.GetCounselingNotes)
	counseling.Post("", h.CreateCounselingNote)
	counseling.Get("/:id", h.GetCounselingNoteByID)
	counseling.Put("/:id", h.UpdateCounselingNote)
	counseling.Delete("/:id", h.DeleteCounselingNote)

	// Student BK Profile
	bk.Get("/students/:studentId/profile", h.GetStudentBKProfile)
	bk.Get("/students/:studentId/violations", h.GetStudentViolations)
	bk.Get("/students/:studentId/achievements", h.GetStudentAchievements)
	bk.Get("/students/:studentId/permits", h.GetStudentPermits)
	bk.Get("/students/:studentId/counseling", h.GetStudentCounselingNotes)
}

// RegisterReadOnlyRoutes registers read-only BK routes for Wali Kelas
// Requirements: 6.5 - WHEN a Wali_Kelas views BK data, THE System SHALL provide read-only access
func (h *Handler) RegisterReadOnlyRoutes(router fiber.Router) {
	bk := router.Group("/bk")

	// Read-only access to violations, achievements, permits
	bk.Get("/violations", h.GetViolations)
	bk.Get("/violations/:id", h.GetViolationByID)
	bk.Get("/achievements", h.GetAchievements)
	bk.Get("/achievements/:id", h.GetAchievementByID)
	bk.Get("/achievements/student/:studentId/points", h.GetStudentAchievementPoints)
	bk.Get("/permits", h.GetPermits)
	bk.Get("/permits/:id", h.GetPermitByID)

	// Counseling notes - only parent_summary visible
	// Requirements: 9.4 - WHEN a Wali_Kelas views counseling data, THE System SHALL show only parent_summary
	bk.Get("/counseling", h.GetCounselingNotesReadOnly)
	bk.Get("/counseling/:id", h.GetCounselingNoteByIDReadOnly)

	// Student BK Profile (read-only, no internal notes)
	bk.Get("/students/:studentId/profile", h.GetStudentBKProfileReadOnly)
	bk.Get("/students/:studentId/violations", h.GetStudentViolations)
	bk.Get("/students/:studentId/achievements", h.GetStudentAchievements)
	bk.Get("/students/:studentId/permits", h.GetStudentPermits)
	bk.Get("/students/:studentId/counseling", h.GetStudentCounselingNotesReadOnly)
}

// ==================== Dashboard ====================

// GetDashboard handles getting BK dashboard data
// @Summary Get BK dashboard
// @Description Get BK dashboard with overview statistics
// @Tags BK
// @Produce json
// @Success 200 {object} BKDashboardResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/bk/dashboard [get]
func (h *Handler) GetDashboard(c *fiber.Ctx) error {
	schoolID, ok := c.Locals("school_id").(uint)
	if !ok {
		return h.tenantRequiredError(c)
	}

	response, err := h.service.GetBKDashboard(c.Context(), schoolID)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// ==================== Violations ====================

// CreateViolation handles creating a new violation
// @Summary Create violation
// @Description Record a new student violation
// @Tags BK - Violations
// @Accept json
// @Produce json
// @Param request body CreateViolationRequest true "Violation data"
// @Success 201 {object} ViolationResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/bk/violations [post]
func (h *Handler) CreateViolation(c *fiber.Ctx) error {
	schoolID, ok := c.Locals("school_id").(uint)
	if !ok {
		return h.tenantRequiredError(c)
	}

	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return h.authRequiredError(c)
	}

	var req CreateViolationRequest
	if err := c.BodyParser(&req); err != nil {
		return h.invalidBodyError(c)
	}

	response, err := h.service.CreateViolation(c.Context(), schoolID, userID, req)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    response,
		"message": "Violation recorded successfully",
	})
}

// GetViolations handles listing violations
// @Summary List violations
// @Description Get a paginated list of violations
// @Tags BK - Violations
// @Produce json
// @Param student_id query int false "Filter by student ID"
// @Param class_id query int false "Filter by class ID"
// @Param level query string false "Filter by level (ringan, sedang, berat)"
// @Param category query string false "Filter by category"
// @Param start_date query string false "Filter by start date (YYYY-MM-DD)"
// @Param end_date query string false "Filter by end date (YYYY-MM-DD)"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} ViolationListResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/bk/violations [get]
func (h *Handler) GetViolations(c *fiber.Ctx) error {
	schoolID, ok := c.Locals("school_id").(uint)
	if !ok {
		return h.tenantRequiredError(c)
	}

	filter := h.parseViolationFilter(c)

	response, err := h.service.GetViolations(c.Context(), schoolID, filter)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GetViolationByID handles getting a single violation
// @Summary Get violation by ID
// @Description Get detailed information about a specific violation
// @Tags BK - Violations
// @Produce json
// @Param id path int true "Violation ID"
// @Success 200 {object} ViolationResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/bk/violations/{id} [get]
func (h *Handler) GetViolationByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return h.invalidIDError(c, "violation")
	}

	response, err := h.service.GetViolationByID(c.Context(), uint(id))
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GetStudentViolations handles getting violations for a specific student
// @Summary Get student violations
// @Description Get all violations for a specific student
// @Tags BK - Violations
// @Produce json
// @Param studentId path int true "Student ID"
// @Success 200 {object} []ViolationResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/bk/students/{studentId}/violations [get]
func (h *Handler) GetStudentViolations(c *fiber.Ctx) error {
	studentID, err := strconv.ParseUint(c.Params("studentId"), 10, 32)
	if err != nil {
		return h.invalidIDError(c, "student")
	}

	responses, err := h.service.GetStudentViolations(c.Context(), uint(studentID))
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    responses,
	})
}

// DeleteViolation handles deleting a violation
// @Summary Delete violation
// @Description Delete a specific violation record
// @Tags BK - Violations
// @Produce json
// @Param id path int true "Violation ID"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/bk/violations/{id} [delete]
func (h *Handler) DeleteViolation(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return h.invalidIDError(c, "violation")
	}

	if err := h.service.DeleteViolation(c.Context(), uint(id)); err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Violation deleted successfully",
	})
}


// ==================== Achievements ====================

// CreateAchievement handles creating a new achievement
// @Summary Create achievement
// @Description Record a new student achievement
// @Tags BK - Achievements
// @Accept json
// @Produce json
// @Param request body CreateAchievementRequest true "Achievement data"
// @Success 201 {object} AchievementResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/bk/achievements [post]
func (h *Handler) CreateAchievement(c *fiber.Ctx) error {
	schoolID, ok := c.Locals("school_id").(uint)
	if !ok {
		return h.tenantRequiredError(c)
	}

	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return h.authRequiredError(c)
	}

	var req CreateAchievementRequest
	if err := c.BodyParser(&req); err != nil {
		return h.invalidBodyError(c)
	}

	response, err := h.service.CreateAchievement(c.Context(), schoolID, userID, req)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    response,
		"message": "Achievement recorded successfully",
	})
}

// GetAchievements handles listing achievements
// @Summary List achievements
// @Description Get a paginated list of achievements
// @Tags BK - Achievements
// @Produce json
// @Param student_id query int false "Filter by student ID"
// @Param class_id query int false "Filter by class ID"
// @Param start_date query string false "Filter by start date (YYYY-MM-DD)"
// @Param end_date query string false "Filter by end date (YYYY-MM-DD)"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} AchievementListResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/bk/achievements [get]
func (h *Handler) GetAchievements(c *fiber.Ctx) error {
	schoolID, ok := c.Locals("school_id").(uint)
	if !ok {
		return h.tenantRequiredError(c)
	}

	filter := h.parseAchievementFilter(c)

	response, err := h.service.GetAchievements(c.Context(), schoolID, filter)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GetAchievementByID handles getting a single achievement
// @Summary Get achievement by ID
// @Description Get detailed information about a specific achievement
// @Tags BK - Achievements
// @Produce json
// @Param id path int true "Achievement ID"
// @Success 200 {object} AchievementResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/bk/achievements/{id} [get]
func (h *Handler) GetAchievementByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return h.invalidIDError(c, "achievement")
	}

	response, err := h.service.GetAchievementByID(c.Context(), uint(id))
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GetStudentAchievements handles getting achievements for a specific student
// @Summary Get student achievements
// @Description Get all achievements for a specific student
// @Tags BK - Achievements
// @Produce json
// @Param studentId path int true "Student ID"
// @Success 200 {object} []AchievementResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/bk/students/{studentId}/achievements [get]
func (h *Handler) GetStudentAchievements(c *fiber.Ctx) error {
	studentID, err := strconv.ParseUint(c.Params("studentId"), 10, 32)
	if err != nil {
		return h.invalidIDError(c, "student")
	}

	responses, err := h.service.GetStudentAchievements(c.Context(), uint(studentID))
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    responses,
	})
}

// GetStudentAchievementPoints handles getting total achievement points for a student
// @Summary Get student achievement points
// @Description Get total achievement points for a specific student
// @Tags BK - Achievements
// @Produce json
// @Param studentId path int true "Student ID"
// @Success 200 {object} AchievementPointsResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/bk/achievements/student/{studentId}/points [get]
func (h *Handler) GetStudentAchievementPoints(c *fiber.Ctx) error {
	studentID, err := strconv.ParseUint(c.Params("studentId"), 10, 32)
	if err != nil {
		return h.invalidIDError(c, "student")
	}

	response, err := h.service.GetStudentAchievementPoints(c.Context(), uint(studentID))
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// DeleteAchievement handles deleting an achievement
// @Summary Delete achievement
// @Description Delete a specific achievement record
// @Tags BK - Achievements
// @Produce json
// @Param id path int true "Achievement ID"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/bk/achievements/{id} [delete]
func (h *Handler) DeleteAchievement(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return h.invalidIDError(c, "achievement")
	}

	if err := h.service.DeleteAchievement(c.Context(), uint(id)); err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Achievement deleted successfully",
	})
}

// ==================== Permits ====================

// CreatePermit handles creating a new exit permit
// @Summary Create exit permit
// @Description Create a new school exit permit for a student
// @Tags BK - Permits
// @Accept json
// @Produce json
// @Param request body CreatePermitRequest true "Permit data"
// @Success 201 {object} PermitResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/bk/permits [post]
func (h *Handler) CreatePermit(c *fiber.Ctx) error {
	schoolID, ok := c.Locals("school_id").(uint)
	if !ok {
		return h.tenantRequiredError(c)
	}

	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return h.authRequiredError(c)
	}

	var req CreatePermitRequest
	if err := c.BodyParser(&req); err != nil {
		return h.invalidBodyError(c)
	}

	response, err := h.service.CreatePermit(c.Context(), schoolID, userID, req)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    response,
		"message": "Exit permit created successfully",
	})
}

// GetPermits handles listing permits
// @Summary List permits
// @Description Get a paginated list of exit permits
// @Tags BK - Permits
// @Produce json
// @Param student_id query int false "Filter by student ID"
// @Param class_id query int false "Filter by class ID"
// @Param teacher_id query int false "Filter by responsible teacher ID"
// @Param has_returned query bool false "Filter by return status"
// @Param start_date query string false "Filter by start date (YYYY-MM-DD)"
// @Param end_date query string false "Filter by end date (YYYY-MM-DD)"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} PermitListResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/bk/permits [get]
func (h *Handler) GetPermits(c *fiber.Ctx) error {
	schoolID, ok := c.Locals("school_id").(uint)
	if !ok {
		return h.tenantRequiredError(c)
	}

	filter := h.parsePermitFilter(c)

	response, err := h.service.GetPermits(c.Context(), schoolID, filter)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GetPermitByID handles getting a single permit
// @Summary Get permit by ID
// @Description Get detailed information about a specific permit
// @Tags BK - Permits
// @Produce json
// @Param id path int true "Permit ID"
// @Success 200 {object} PermitResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/bk/permits/{id} [get]
func (h *Handler) GetPermitByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return h.invalidIDError(c, "permit")
	}

	response, err := h.service.GetPermitByID(c.Context(), uint(id))
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GetStudentPermits handles getting permits for a specific student
// @Summary Get student permits
// @Description Get all permits for a specific student
// @Tags BK - Permits
// @Produce json
// @Param studentId path int true "Student ID"
// @Success 200 {object} []PermitResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/bk/students/{studentId}/permits [get]
func (h *Handler) GetStudentPermits(c *fiber.Ctx) error {
	studentID, err := strconv.ParseUint(c.Params("studentId"), 10, 32)
	if err != nil {
		return h.invalidIDError(c, "student")
	}

	responses, err := h.service.GetStudentPermits(c.Context(), uint(studentID))
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    responses,
	})
}

// RecordReturn handles recording student return time
// @Summary Record student return
// @Description Record the return time for a student with an exit permit
// @Tags BK - Permits
// @Accept json
// @Produce json
// @Param id path int true "Permit ID"
// @Param request body RecordReturnRequest true "Return data"
// @Success 200 {object} PermitResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/bk/permits/{id}/return [post]
func (h *Handler) RecordReturn(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return h.invalidIDError(c, "permit")
	}

	var req RecordReturnRequest
	if err := c.BodyParser(&req); err != nil {
		return h.invalidBodyError(c)
	}

	// Default to current time if not provided
	if req.ReturnTime.IsZero() {
		req.ReturnTime = time.Now()
	}

	response, err := h.service.RecordReturn(c.Context(), uint(id), req)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
		"message": "Return time recorded successfully",
	})
}

// GetPermitDocument handles getting permit document data
// @Summary Get permit document
// @Description Get permit document data for printing/PDF generation
// @Tags BK - Permits
// @Produce json
// @Param id path int true "Permit ID"
// @Success 200 {object} PermitDocumentData
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/bk/permits/{id}/document [get]
func (h *Handler) GetPermitDocument(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return h.invalidIDError(c, "permit")
	}

	response, err := h.service.GetPermitDocument(c.Context(), uint(id))
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// DeletePermit handles deleting a permit
// @Summary Delete permit
// @Description Delete a specific permit record
// @Tags BK - Permits
// @Produce json
// @Param id path int true "Permit ID"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/bk/permits/{id} [delete]
func (h *Handler) DeletePermit(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return h.invalidIDError(c, "permit")
	}

	if err := h.service.DeletePermit(c.Context(), uint(id)); err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Permit deleted successfully",
	})
}


// ==================== Counseling Notes ====================

// CreateCounselingNote handles creating a new counseling note
// @Summary Create counseling note
// @Description Create a new counseling note for a student
// @Tags BK - Counseling
// @Accept json
// @Produce json
// @Param request body CreateCounselingNoteRequest true "Counseling note data"
// @Success 201 {object} CounselingNoteFullResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/bk/counseling [post]
func (h *Handler) CreateCounselingNote(c *fiber.Ctx) error {
	schoolID, ok := c.Locals("school_id").(uint)
	if !ok {
		return h.tenantRequiredError(c)
	}

	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return h.authRequiredError(c)
	}

	var req CreateCounselingNoteRequest
	if err := c.BodyParser(&req); err != nil {
		return h.invalidBodyError(c)
	}

	response, err := h.service.CreateCounselingNote(c.Context(), schoolID, userID, req)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    response,
		"message": "Counseling note created successfully",
	})
}

// GetCounselingNotes handles listing counseling notes (with internal notes for Guru BK)
// @Summary List counseling notes
// @Description Get a paginated list of counseling notes with internal notes
// @Tags BK - Counseling
// @Produce json
// @Param student_id query int false "Filter by student ID"
// @Param class_id query int false "Filter by class ID"
// @Param start_date query string false "Filter by start date (YYYY-MM-DD)"
// @Param end_date query string false "Filter by end date (YYYY-MM-DD)"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} CounselingNoteFullListResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/bk/counseling [get]
func (h *Handler) GetCounselingNotes(c *fiber.Ctx) error {
	schoolID, ok := c.Locals("school_id").(uint)
	if !ok {
		return h.tenantRequiredError(c)
	}

	filter := h.parseCounselingNoteFilter(c)

	response, err := h.service.GetCounselingNotes(c.Context(), schoolID, filter, true)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GetCounselingNotesReadOnly handles listing counseling notes (without internal notes)
// Requirements: 9.4 - WHEN a Wali_Kelas views counseling data, THE System SHALL show only parent_summary
func (h *Handler) GetCounselingNotesReadOnly(c *fiber.Ctx) error {
	schoolID, ok := c.Locals("school_id").(uint)
	if !ok {
		return h.tenantRequiredError(c)
	}

	filter := h.parseCounselingNoteFilter(c)

	response, err := h.service.GetCounselingNotes(c.Context(), schoolID, filter, false)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GetCounselingNoteByID handles getting a single counseling note (with internal note)
// @Summary Get counseling note by ID
// @Description Get detailed information about a specific counseling note
// @Tags BK - Counseling
// @Produce json
// @Param id path int true "Counseling Note ID"
// @Success 200 {object} CounselingNoteFullResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/bk/counseling/{id} [get]
func (h *Handler) GetCounselingNoteByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return h.invalidIDError(c, "counseling note")
	}

	response, err := h.service.GetCounselingNoteByID(c.Context(), uint(id), true)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GetCounselingNoteByIDReadOnly handles getting a single counseling note (without internal note)
// Requirements: 9.4 - WHEN a Wali_Kelas views counseling data, THE System SHALL show only parent_summary
func (h *Handler) GetCounselingNoteByIDReadOnly(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return h.invalidIDError(c, "counseling note")
	}

	response, err := h.service.GetCounselingNoteByID(c.Context(), uint(id), false)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GetStudentCounselingNotes handles getting counseling notes for a specific student (with internal notes)
// @Summary Get student counseling notes
// @Description Get all counseling notes for a specific student
// @Tags BK - Counseling
// @Produce json
// @Param studentId path int true "Student ID"
// @Success 200 {object} []CounselingNoteFullResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/bk/students/{studentId}/counseling [get]
func (h *Handler) GetStudentCounselingNotes(c *fiber.Ctx) error {
	studentID, err := strconv.ParseUint(c.Params("studentId"), 10, 32)
	if err != nil {
		return h.invalidIDError(c, "student")
	}

	responses, err := h.service.GetStudentCounselingNotes(c.Context(), uint(studentID), true)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    responses,
	})
}

// GetStudentCounselingNotesReadOnly handles getting counseling notes for a specific student (without internal notes)
// Requirements: 9.4 - WHEN a Wali_Kelas views counseling data, THE System SHALL show only parent_summary
func (h *Handler) GetStudentCounselingNotesReadOnly(c *fiber.Ctx) error {
	studentID, err := strconv.ParseUint(c.Params("studentId"), 10, 32)
	if err != nil {
		return h.invalidIDError(c, "student")
	}

	responses, err := h.service.GetStudentCounselingNotes(c.Context(), uint(studentID), false)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    responses,
	})
}

// UpdateCounselingNote handles updating a counseling note
// @Summary Update counseling note
// @Description Update an existing counseling note
// @Tags BK - Counseling
// @Accept json
// @Produce json
// @Param id path int true "Counseling Note ID"
// @Param request body UpdateCounselingNoteRequest true "Updated counseling note data"
// @Success 200 {object} CounselingNoteFullResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/bk/counseling/{id} [put]
func (h *Handler) UpdateCounselingNote(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return h.invalidIDError(c, "counseling note")
	}

	var req UpdateCounselingNoteRequest
	if err := c.BodyParser(&req); err != nil {
		return h.invalidBodyError(c)
	}

	response, err := h.service.UpdateCounselingNote(c.Context(), uint(id), req)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
		"message": "Counseling note updated successfully",
	})
}

// DeleteCounselingNote handles deleting a counseling note
// @Summary Delete counseling note
// @Description Delete a specific counseling note
// @Tags BK - Counseling
// @Produce json
// @Param id path int true "Counseling Note ID"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/bk/counseling/{id} [delete]
func (h *Handler) DeleteCounselingNote(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return h.invalidIDError(c, "counseling note")
	}

	if err := h.service.DeleteCounselingNote(c.Context(), uint(id)); err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Counseling note deleted successfully",
	})
}

// ==================== Student BK Profile ====================

// GetStudentBKProfile handles getting a student's complete BK profile (with internal notes)
// @Summary Get student BK profile
// @Description Get complete BK profile for a student including violations, achievements, permits, and counseling notes
// @Tags BK - Profile
// @Produce json
// @Param studentId path int true "Student ID"
// @Success 200 {object} StudentBKProfileFullResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/bk/students/{studentId}/profile [get]
func (h *Handler) GetStudentBKProfile(c *fiber.Ctx) error {
	studentID, err := strconv.ParseUint(c.Params("studentId"), 10, 32)
	if err != nil {
		return h.invalidIDError(c, "student")
	}

	response, err := h.service.GetStudentBKProfile(c.Context(), uint(studentID), true)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GetStudentBKProfileReadOnly handles getting a student's BK profile (without internal notes)
// Requirements: 6.5 - WHEN a Wali_Kelas views BK data, THE System SHALL provide read-only access
func (h *Handler) GetStudentBKProfileReadOnly(c *fiber.Ctx) error {
	studentID, err := strconv.ParseUint(c.Params("studentId"), 10, 32)
	if err != nil {
		return h.invalidIDError(c, "student")
	}

	response, err := h.service.GetStudentBKProfile(c.Context(), uint(studentID), false)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}


// ==================== Helper Functions ====================

func (h *Handler) parseViolationFilter(c *fiber.Ctx) ViolationFilter {
	filter := ViolationFilter{
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

	if level := c.Query("level"); level != "" {
		filter.Level = &level
	}

	if category := c.Query("category"); category != "" {
		filter.Category = &category
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

func (h *Handler) parseAchievementFilter(c *fiber.Ctx) AchievementFilter {
	filter := AchievementFilter{
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

func (h *Handler) parsePermitFilter(c *fiber.Ctx) PermitFilter {
	filter := PermitFilter{
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

	if hasReturnedStr := c.Query("has_returned"); hasReturnedStr != "" {
		hasReturned := hasReturnedStr == "true"
		filter.HasReturned = &hasReturned
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

func (h *Handler) parseCounselingNoteFilter(c *fiber.Ctx) CounselingNoteFilter {
	filter := CounselingNoteFilter{
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
			"message": "School context is required",
		},
	})
}

func (h *Handler) authRequiredError(c *fiber.Ctx) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"success": false,
		"error": fiber.Map{
			"code":    "AUTH_REQUIRED",
			"message": "Authentication is required",
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
	case errors.Is(err, ErrViolationNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "NOT_FOUND_VIOLATION",
				"message": "Violation not found",
			},
		})
	case errors.Is(err, ErrAchievementNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "NOT_FOUND_ACHIEVEMENT",
				"message": "Achievement not found",
			},
		})
	case errors.Is(err, ErrPermitNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "NOT_FOUND_PERMIT",
				"message": "Permit not found",
			},
		})
	case errors.Is(err, ErrCounselingNoteNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "NOT_FOUND_COUNSELING_NOTE",
				"message": "Counseling note not found",
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
	case errors.Is(err, ErrUserNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "NOT_FOUND_USER",
				"message": "User not found",
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
	case errors.Is(err, ErrCategoryRequired):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "Category is required",
			},
		})
	case errors.Is(err, ErrLevelRequired):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "Level is required",
			},
		})
	case errors.Is(err, ErrDescriptionRequired):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "Description is required",
			},
		})
	case errors.Is(err, ErrTitleRequired):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "Title is required",
			},
		})
	case errors.Is(err, ErrPointRequired):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_VALUE",
				"message": "Point must be greater than 0",
			},
		})
	case errors.Is(err, ErrReasonRequired):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "Reason is required",
			},
		})
	case errors.Is(err, ErrExitTimeRequired):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "Exit time is required",
			},
		})
	case errors.Is(err, ErrResponsibleTeacherRequired):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "Responsible teacher is required",
			},
		})
	case errors.Is(err, ErrInternalNoteRequired):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "Internal note is required",
			},
		})
	case errors.Is(err, ErrReturnTimeRequired):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "Return time is required",
			},
		})
	case errors.Is(err, ErrAlreadyReturned):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_ALREADY_RETURNED",
				"message": "Student has already returned",
			},
		})
	case errors.Is(err, ErrStudentNotInSchool):
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTHZ_TENANT_MISMATCH",
				"message": "Student does not belong to this school",
			},
		})
	case errors.Is(err, ErrTeacherNotInSchool):
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTHZ_TENANT_MISMATCH",
				"message": "Teacher does not belong to this school",
			},
		})
	case errors.Is(err, ErrInvalidViolationLevel):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_VALUE",
				"message": "Invalid violation level. Must be one of: ringan, sedang, berat",
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
