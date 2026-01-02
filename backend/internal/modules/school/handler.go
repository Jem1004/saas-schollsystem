package school

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/school-management/backend/internal/middleware"
)

// Handler handles HTTP requests for school management (classes, students, parents)
type Handler struct {
	service Service
}

// NewHandler creates a new school handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers school routes (Admin Sekolah)
// Requirements: 3.1, 3.2, 3.3 - Admin Sekolah manages school data
func (h *Handler) RegisterRoutes(router fiber.Router) {
	// Stats route
	router.Get("/stats", h.GetStats)

	// Devices route (for admin sekolah to get their school's devices)
	router.Get("/devices", h.GetSchoolDevices)

	// Class routes
	classes := router.Group("/classes")
	classes.Post("", h.CreateClass)
	classes.Get("", h.GetClasses)
	classes.Get("/:id", h.GetClass)
	classes.Put("/:id", h.UpdateClass)
	classes.Delete("/:id", h.DeleteClass)
	classes.Post("/:id/homeroom-teacher", h.AssignHomeroomTeacher)
	classes.Get("/:id/students", h.GetClassStudents)

	// Student routes
	students := router.Group("/students")
	students.Post("", h.CreateStudent)
	students.Get("", h.GetStudents)
	students.Get("/without-class", h.GetStudentsWithoutClass)
	students.Get("/search", h.SearchStudents)
	students.Post("/bulk-assign-class", h.BulkAssignClass)
	students.Get("/:id", h.GetStudent)
	students.Put("/:id", h.UpdateStudent)
	students.Delete("/:id", h.DeleteStudent)
	students.Post("/:id/account", h.CreateStudentAccount)
	students.Post("/:id/reset-password", h.ResetStudentPassword)
	students.Post("/:id/clear-rfid", h.ClearStudentRFID)

	// Parent routes
	parents := router.Group("/parents")
	parents.Post("", h.CreateParent)
	parents.Get("", h.GetParents)
	parents.Get("/:id", h.GetParent)
	parents.Put("/:id", h.UpdateParent)
	parents.Delete("/:id", h.DeleteParent)
	parents.Post("/:id/students", h.LinkParentToStudents)
	parents.Post("/:id/reset-password", h.ResetParentPassword)

	// User routes (school staff)
	users := router.Group("/users")
	users.Post("", h.CreateUser)
	users.Get("", h.GetUsers)
	users.Get("/:id", h.GetUser)
	users.Put("/:id", h.UpdateUser)
	users.Delete("/:id", h.DeleteUser)
	users.Post("/:id/reset-password", h.ResetUserPassword)
}

// GetStats handles getting school statistics for dashboard
// @Summary Get school statistics
// @Description Get statistics for the school dashboard (students, classes, teachers, parents, attendance)
// @Tags School
// @Produce json
// @Success 200 {object} SchoolStatsResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/stats [get]
func (h *Handler) GetStats(c *fiber.Ctx) error {
	schoolID, ok := middleware.GetTenantID(c)
	if !ok {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTHZ_TENANT_REQUIRED",
				"message": "Konteks sekolah diperlukan",
			},
		})
	}

	stats, err := h.service.GetStats(c.Context(), schoolID)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    stats,
	})
}

// ==================== Class Handlers ====================

// CreateClass handles class creation
// @Summary Create a new class
// @Description Create a new class in the school
// @Tags Classes
// @Accept json
// @Produce json
// @Param request body CreateClassRequest true "Class data"
// @Success 201 {object} ClassResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/classes [post]
func (h *Handler) CreateClass(c *fiber.Ctx) error {
	schoolID, ok := middleware.GetTenantID(c)
	if !ok {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTHZ_TENANT_REQUIRED",
				"message": "Konteks sekolah diperlukan",
			},
		})
	}

	var req CreateClassRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Format data tidak valid",
			},
		})
	}

	response, err := h.service.CreateClass(c.Context(), schoolID, req)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    response,
		"message": "Kelas berhasil dibuat",
	})
}

// GetClasses handles listing all classes
// @Summary List all classes
// @Description Get a paginated list of all classes in the school
// @Tags Classes
// @Produce json
// @Param name query string false "Filter by class name"
// @Param grade query int false "Filter by grade"
// @Param year query string false "Filter by academic year"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} ClassListResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/classes [get]
func (h *Handler) GetClasses(c *fiber.Ctx) error {
	schoolID, ok := middleware.GetTenantID(c)
	if !ok {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTHZ_TENANT_REQUIRED",
				"message": "Konteks sekolah diperlukan",
			},
		})
	}

	filter := DefaultClassFilter()
	filter.Name = c.Query("name")
	filter.Year = c.Query("year")

	if gradeStr := c.Query("grade"); gradeStr != "" {
		if grade, err := strconv.Atoi(gradeStr); err == nil {
			filter.Grade = &grade
		}
	}
	if page, err := strconv.Atoi(c.Query("page", "1")); err == nil && page > 0 {
		filter.Page = page
	}
	if pageSize, err := strconv.Atoi(c.Query("page_size", "20")); err == nil && pageSize > 0 {
		filter.PageSize = pageSize
	}

	response, err := h.service.GetAllClasses(c.Context(), schoolID, filter)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GetClass handles getting a single class
// @Summary Get a class by ID
// @Description Get detailed information about a specific class
// @Tags Classes
// @Produce json
// @Param id path int true "Class ID"
// @Success 200 {object} ClassResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/classes/{id} [get]
func (h *Handler) GetClass(c *fiber.Ctx) error {
	schoolID, ok := middleware.GetTenantID(c)
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
				"message": "ID kelas tidak valid",
			},
		})
	}

	response, err := h.service.GetClassByID(c.Context(), schoolID, uint(id))
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}


// UpdateClass handles updating a class
// @Summary Update a class
// @Description Update class information
// @Tags Classes
// @Accept json
// @Produce json
// @Param id path int true "Class ID"
// @Param request body UpdateClassRequest true "Class data to update"
// @Success 200 {object} ClassResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/classes/{id} [put]
func (h *Handler) UpdateClass(c *fiber.Ctx) error {
	schoolID, ok := middleware.GetTenantID(c)
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
				"message": "ID kelas tidak valid",
			},
		})
	}

	var req UpdateClassRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Format data tidak valid",
			},
		})
	}

	response, err := h.service.UpdateClass(c.Context(), schoolID, uint(id), req)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
		"message": "Kelas berhasil diperbarui",
	})
}

// DeleteClass handles deleting a class
// @Summary Delete a class
// @Description Delete a class from the school
// @Tags Classes
// @Produce json
// @Param id path int true "Class ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/classes/{id} [delete]
func (h *Handler) DeleteClass(c *fiber.Ctx) error {
	schoolID, ok := middleware.GetTenantID(c)
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
				"message": "ID kelas tidak valid",
			},
		})
	}

	if err := h.service.DeleteClass(c.Context(), schoolID, uint(id)); err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Kelas berhasil dihapus",
	})
}

// AssignHomeroomTeacher handles assigning a homeroom teacher to a class
// @Summary Assign homeroom teacher
// @Description Assign a homeroom teacher to a class
// @Tags Classes
// @Accept json
// @Produce json
// @Param id path int true "Class ID"
// @Param request body AssignHomeroomTeacherRequest true "Teacher assignment data"
// @Success 200 {object} ClassResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/classes/{id}/homeroom-teacher [post]
func (h *Handler) AssignHomeroomTeacher(c *fiber.Ctx) error {
	schoolID, ok := middleware.GetTenantID(c)
	if !ok {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTHZ_TENANT_REQUIRED",
				"message": "Konteks sekolah diperlukan",
			},
		})
	}

	classID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "ID kelas tidak valid",
			},
		})
	}

	var req AssignHomeroomTeacherRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Format data tidak valid",
			},
		})
	}

	if req.TeacherID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "ID guru wajib diisi",
			},
		})
	}

	response, err := h.service.AssignHomeroomTeacher(c.Context(), schoolID, uint(classID), req.TeacherID)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
		"message": "Wali kelas berhasil ditugaskan",
	})
}

// GetClassStudents handles getting all students in a class
// @Summary Get class students
// @Description Get all students in a specific class
// @Tags Classes
// @Produce json
// @Param id path int true "Class ID"
// @Success 200 {object} []StudentResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/classes/{id}/students [get]
func (h *Handler) GetClassStudents(c *fiber.Ctx) error {
	schoolID, ok := middleware.GetTenantID(c)
	if !ok {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTHZ_TENANT_REQUIRED",
				"message": "Konteks sekolah diperlukan",
			},
		})
	}

	classID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "ID kelas tidak valid",
			},
		})
	}

	students, err := h.service.GetStudentsByClass(c.Context(), schoolID, uint(classID))
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    students,
	})
}


// ==================== Student Handlers ====================

// CreateStudent handles student creation
// @Summary Create a new student
// @Description Create a new student in the school
// @Tags Students
// @Accept json
// @Produce json
// @Param request body CreateStudentRequest true "Student data"
// @Success 201 {object} StudentResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/students [post]
func (h *Handler) CreateStudent(c *fiber.Ctx) error {
	schoolID, ok := middleware.GetTenantID(c)
	if !ok {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTHZ_TENANT_REQUIRED",
				"message": "Konteks sekolah diperlukan",
			},
		})
	}

	var req CreateStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Format data tidak valid",
			},
		})
	}

	response, err := h.service.CreateStudent(c.Context(), schoolID, req)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    response,
		"message": "Siswa berhasil ditambahkan",
	})
}

// GetStudents handles listing all students
// @Summary List all students
// @Description Get a paginated list of all students in the school
// @Tags Students
// @Produce json
// @Param name query string false "Filter by student name"
// @Param nis query string false "Filter by NIS"
// @Param nisn query string false "Filter by NISN"
// @Param class_id query int false "Filter by class ID"
// @Param is_active query bool false "Filter by active status"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} StudentListResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/students [get]
func (h *Handler) GetStudents(c *fiber.Ctx) error {
	schoolID, ok := middleware.GetTenantID(c)
	if !ok {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTHZ_TENANT_REQUIRED",
				"message": "Konteks sekolah diperlukan",
			},
		})
	}

	filter := DefaultStudentFilter()
	filter.Name = c.Query("name")
	filter.NIS = c.Query("nis")
	filter.NISN = c.Query("nisn")

	if classIDStr := c.Query("class_id"); classIDStr != "" {
		if classID, err := strconv.ParseUint(classIDStr, 10, 32); err == nil {
			cid := uint(classID)
			filter.ClassID = &cid
		}
	}
	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		isActive := isActiveStr == "true"
		filter.IsActive = &isActive
	}
	if page, err := strconv.Atoi(c.Query("page", "1")); err == nil && page > 0 {
		filter.Page = page
	}
	if pageSize, err := strconv.Atoi(c.Query("page_size", "20")); err == nil && pageSize > 0 {
		filter.PageSize = pageSize
	}

	response, err := h.service.GetAllStudents(c.Context(), schoolID, filter)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GetStudent handles getting a single student
// @Summary Get a student by ID
// @Description Get detailed information about a specific student
// @Tags Students
// @Produce json
// @Param id path int true "Student ID"
// @Success 200 {object} StudentResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/students/{id} [get]
func (h *Handler) GetStudent(c *fiber.Ctx) error {
	schoolID, ok := middleware.GetTenantID(c)
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
				"message": "ID siswa tidak valid",
			},
		})
	}

	response, err := h.service.GetStudentByID(c.Context(), schoolID, uint(id))
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// UpdateStudent handles updating a student
// @Summary Update a student
// @Description Update student information
// @Tags Students
// @Accept json
// @Produce json
// @Param id path int true "Student ID"
// @Param request body UpdateStudentRequest true "Student data to update"
// @Success 200 {object} StudentResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/students/{id} [put]
func (h *Handler) UpdateStudent(c *fiber.Ctx) error {
	schoolID, ok := middleware.GetTenantID(c)
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
				"message": "ID siswa tidak valid",
			},
		})
	}

	var req UpdateStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Format data tidak valid",
			},
		})
	}

	response, err := h.service.UpdateStudent(c.Context(), schoolID, uint(id), req)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
		"message": "Siswa berhasil diperbarui",
	})
}

// DeleteStudent handles deleting a student
// @Summary Delete a student
// @Description Delete a student from the school
// @Tags Students
// @Produce json
// @Param id path int true "Student ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/students/{id} [delete]
func (h *Handler) DeleteStudent(c *fiber.Ctx) error {
	schoolID, ok := middleware.GetTenantID(c)
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
				"message": "ID siswa tidak valid",
			},
		})
	}

	if err := h.service.DeleteStudent(c.Context(), schoolID, uint(id)); err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Siswa berhasil dihapus",
	})
}


// ==================== Parent Handlers ====================

// CreateParent handles parent creation
// @Summary Create a new parent
// @Description Create a new parent with user account in the school
// @Tags Parents
// @Accept json
// @Produce json
// @Param request body CreateParentRequest true "Parent data"
// @Success 201 {object} ParentResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/parents [post]
func (h *Handler) CreateParent(c *fiber.Ctx) error {
	schoolID, ok := middleware.GetTenantID(c)
	if !ok {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTHZ_TENANT_REQUIRED",
				"message": "Konteks sekolah diperlukan",
			},
		})
	}

	var req CreateParentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Format data tidak valid",
			},
		})
	}

	response, err := h.service.CreateParent(c.Context(), schoolID, req)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    response,
		"message": "Orang tua berhasil ditambahkan",
	})
}

// GetParents handles listing all parents
// @Summary List all parents
// @Description Get a paginated list of all parents in the school
// @Tags Parents
// @Produce json
// @Param name query string false "Filter by parent name"
// @Param phone query string false "Filter by phone number"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} ParentListResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/parents [get]
func (h *Handler) GetParents(c *fiber.Ctx) error {
	schoolID, ok := middleware.GetTenantID(c)
	if !ok {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTHZ_TENANT_REQUIRED",
				"message": "Konteks sekolah diperlukan",
			},
		})
	}

	filter := DefaultParentFilter()
	filter.Name = c.Query("name")
	filter.Phone = c.Query("phone")

	if page, err := strconv.Atoi(c.Query("page", "1")); err == nil && page > 0 {
		filter.Page = page
	}
	if pageSize, err := strconv.Atoi(c.Query("page_size", "20")); err == nil && pageSize > 0 {
		filter.PageSize = pageSize
	}

	response, err := h.service.GetAllParents(c.Context(), schoolID, filter)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GetParent handles getting a single parent
// @Summary Get a parent by ID
// @Description Get detailed information about a specific parent
// @Tags Parents
// @Produce json
// @Param id path int true "Parent ID"
// @Success 200 {object} ParentResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/parents/{id} [get]
func (h *Handler) GetParent(c *fiber.Ctx) error {
	schoolID, ok := middleware.GetTenantID(c)
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
				"message": "ID orang tua tidak valid",
			},
		})
	}

	response, err := h.service.GetParentByID(c.Context(), schoolID, uint(id))
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// UpdateParent handles updating a parent
// @Summary Update a parent
// @Description Update parent information
// @Tags Parents
// @Accept json
// @Produce json
// @Param id path int true "Parent ID"
// @Param request body UpdateParentRequest true "Parent data to update"
// @Success 200 {object} ParentResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/parents/{id} [put]
func (h *Handler) UpdateParent(c *fiber.Ctx) error {
	schoolID, ok := middleware.GetTenantID(c)
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
				"message": "ID orang tua tidak valid",
			},
		})
	}

	var req UpdateParentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Format data tidak valid",
			},
		})
	}

	response, err := h.service.UpdateParent(c.Context(), schoolID, uint(id), req)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
		"message": "Orang tua berhasil diperbarui",
	})
}

// DeleteParent handles deleting a parent
// @Summary Delete a parent
// @Description Delete a parent from the school
// @Tags Parents
// @Produce json
// @Param id path int true "Parent ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/parents/{id} [delete]
func (h *Handler) DeleteParent(c *fiber.Ctx) error {
	schoolID, ok := middleware.GetTenantID(c)
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
				"message": "ID orang tua tidak valid",
			},
		})
	}

	if err := h.service.DeleteParent(c.Context(), schoolID, uint(id)); err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Orang tua berhasil dihapus",
	})
}

// LinkParentToStudents handles linking a parent to students
// @Summary Link parent to students
// @Description Link a parent to one or more students
// @Tags Parents
// @Accept json
// @Produce json
// @Param id path int true "Parent ID"
// @Param request body LinkParentStudentRequest true "Student IDs to link"
// @Success 200 {object} ParentResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/parents/{id}/students [post]
func (h *Handler) LinkParentToStudents(c *fiber.Ctx) error {
	schoolID, ok := middleware.GetTenantID(c)
	if !ok {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTHZ_TENANT_REQUIRED",
				"message": "Konteks sekolah diperlukan",
			},
		})
	}

	parentID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "ID orang tua tidak valid",
			},
		})
	}

	var req LinkParentStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Format data tidak valid",
			},
		})
	}

	if len(req.StudentIDs) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "Minimal satu ID siswa wajib diisi",
			},
		})
	}

	response, err := h.service.LinkParentToStudents(c.Context(), schoolID, uint(parentID), req.StudentIDs)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
		"message": "Orang tua berhasil dihubungkan dengan siswa",
	})
}

// ResetParentPassword handles resetting a parent's password
// @Summary Reset parent password
// @Description Reset a parent's password to default
// @Tags Parents
// @Produce json
// @Param id path int true "Parent ID"
// @Success 200 {object} ResetPasswordResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/parents/{id}/reset-password [post]
func (h *Handler) ResetParentPassword(c *fiber.Ctx) error {
	schoolID, ok := middleware.GetTenantID(c)
	if !ok {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTHZ_TENANT_REQUIRED",
				"message": "Konteks sekolah diperlukan",
			},
		})
	}

	parentID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "ID orang tua tidak valid",
			},
		})
	}

	response, err := h.service.ResetParentPassword(c.Context(), schoolID, uint(parentID))
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
		"message": "Password berhasil direset",
	})
}

// CreateStudentAccount handles creating a user account for a student
// @Summary Create student account
// @Description Create a user account for a student to login to mobile app
// @Tags Students
// @Produce json
// @Param id path int true "Student ID"
// @Success 201 {object} StudentResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/students/{id}/account [post]
func (h *Handler) CreateStudentAccount(c *fiber.Ctx) error {
	schoolID, ok := middleware.GetTenantID(c)
	if !ok {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTHZ_TENANT_REQUIRED",
				"message": "Konteks sekolah diperlukan",
			},
		})
	}

	studentID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "ID siswa tidak valid",
			},
		})
	}

	response, err := h.service.CreateStudentAccount(c.Context(), schoolID, uint(studentID))
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    response,
		"message": "Akun siswa berhasil dibuat",
	})
}

// ResetStudentPassword handles resetting a student's password
// @Summary Reset student password
// @Description Reset a student's password to default
// @Tags Students
// @Produce json
// @Param id path int true "Student ID"
// @Success 200 {object} ResetPasswordResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/students/{id}/reset-password [post]
func (h *Handler) ResetStudentPassword(c *fiber.Ctx) error {
	schoolID, ok := middleware.GetTenantID(c)
	if !ok {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTHZ_TENANT_REQUIRED",
				"message": "Konteks sekolah diperlukan",
			},
		})
	}

	studentID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "ID siswa tidak valid",
			},
		})
	}

	response, err := h.service.ResetStudentPassword(c.Context(), schoolID, uint(studentID))
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
		"message": "Password berhasil direset",
	})
}


// ==================== Error Handling ====================

// handleError handles service errors and returns appropriate HTTP responses
func (h *Handler) handleError(c *fiber.Ctx, err error) error {
	switch {
	// Class errors
	case errors.Is(err, ErrClassNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "NOT_FOUND_CLASS",
				"message": "Kelas tidak ditemukan",
			},
		})
	case errors.Is(err, ErrDuplicateClass):
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_DUPLICATE_ENTRY",
				"message": "Kelas dengan nama ini sudah ada untuk tingkat dan tahun yang sama",
			},
		})
	case errors.Is(err, ErrClassHasStudents):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_CONSTRAINT_VIOLATION",
				"message": "Tidak dapat menghapus kelas yang masih memiliki siswa. Pindahkan atau hapus siswa terlebih dahulu.",
			},
		})

	// Student errors
	case errors.Is(err, ErrStudentNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "NOT_FOUND_STUDENT",
				"message": "Siswa tidak ditemukan",
			},
		})
	case errors.Is(err, ErrDuplicateNISN):
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_DUPLICATE_NISN",
				"message": "Siswa dengan NISN ini sudah terdaftar di sistem",
			},
		})
	case errors.Is(err, ErrDuplicateNIS):
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_DUPLICATE_NIS",
				"message": "Siswa dengan NIS ini sudah terdaftar di sekolah ini",
			},
		})

	// Parent errors
	case errors.Is(err, ErrParentNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "NOT_FOUND_PARENT",
				"message": "Orang tua tidak ditemukan",
			},
		})

	// Teacher errors
	case errors.Is(err, ErrTeacherNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "NOT_FOUND_TEACHER",
				"message": "Guru tidak ditemukan",
			},
		})
	case errors.Is(err, ErrInvalidTeacher):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_TEACHER",
				"message": "User bukan guru yang valid untuk ditugaskan sebagai wali kelas",
			},
		})

	// User errors
	case errors.Is(err, ErrUserNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "NOT_FOUND_USER",
				"message": "User tidak ditemukan",
			},
		})

	// Validation errors
	case errors.Is(err, ErrNameRequired):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "Nama wajib diisi",
			},
		})
	case errors.Is(err, ErrNISRequired):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "NIS wajib diisi",
			},
		})
	case errors.Is(err, ErrNISNRequired):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "NISN wajib diisi",
			},
		})
	case errors.Is(err, ErrClassIDRequired):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "ID kelas wajib diisi",
			},
		})
	case errors.Is(err, ErrGradeRequired):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "Tingkat wajib diisi dan harus lebih dari 0",
			},
		})
	case errors.Is(err, ErrYearRequired):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "Tahun ajaran wajib diisi",
			},
		})
	case errors.Is(err, ErrPasswordRequired):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "Password wajib diisi",
			},
		})
	case errors.Is(err, ErrPasswordTooShort):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Password minimal 8 karakter",
			},
		})
	case errors.Is(err, ErrStudentIDsRequired):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "Minimal satu ID siswa wajib diisi",
			},
		})
	case errors.Is(err, ErrPhoneRequired):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "Nomor HP wajib diisi",
			},
		})
	case errors.Is(err, ErrStudentHasAccount):
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_DUPLICATE_ACCOUNT",
				"message": "Siswa sudah memiliki akun",
			},
		})
	case errors.Is(err, ErrStudentNoAccount):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_NO_ACCOUNT",
				"message": "Siswa belum memiliki akun",
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

// ==================== User Handlers ====================

// GetUsers handles listing all users (school staff)
// @Summary List all users
// @Description Get a paginated list of all users (school staff) in the school
// @Tags Users
// @Produce json
// @Param name query string false "Filter by name"
// @Param role query string false "Filter by role"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} UserListResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/school/users [get]
func (h *Handler) GetUsers(c *fiber.Ctx) error {
	schoolID, ok := middleware.GetTenantID(c)
	if !ok {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTHZ_TENANT_REQUIRED",
				"message": "Konteks sekolah diperlukan",
			},
		})
	}

	filter := DefaultUserFilter()
	filter.Name = c.Query("name")
	filter.Role = c.Query("role")

	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		isActive := isActiveStr == "true"
		filter.IsActive = &isActive
	}
	if page, err := strconv.Atoi(c.Query("page", "1")); err == nil && page > 0 {
		filter.Page = page
	}
	if pageSize, err := strconv.Atoi(c.Query("page_size", "20")); err == nil && pageSize > 0 {
		filter.PageSize = pageSize
	}

	response, err := h.service.GetAllUsers(c.Context(), schoolID, filter)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GetUser handles getting a single user
// @Summary Get a user by ID
// @Description Get detailed information about a specific user
// @Tags Users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} UserResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/school/users/{id} [get]
func (h *Handler) GetUser(c *fiber.Ctx) error {
	schoolID, ok := middleware.GetTenantID(c)
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
				"message": "ID user tidak valid",
			},
		})
	}

	response, err := h.service.GetUserByID(c.Context(), schoolID, uint(id))
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// CreateUser handles user creation
// @Summary Create a new user
// @Description Create a new user (school staff) in the school
// @Tags Users
// @Accept json
// @Produce json
// @Param request body CreateUserRequest true "User data"
// @Success 201 {object} UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/school/users [post]
func (h *Handler) CreateUser(c *fiber.Ctx) error {
	schoolID, ok := middleware.GetTenantID(c)
	if !ok {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTHZ_TENANT_REQUIRED",
				"message": "Konteks sekolah diperlukan",
			},
		})
	}

	var req CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Format data tidak valid",
			},
		})
	}

	response, err := h.service.CreateUser(c.Context(), schoolID, req)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    response,
		"message": "User berhasil dibuat",
	})
}

// UpdateUser handles updating a user
// @Summary Update a user
// @Description Update user information
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body UpdateUserRequest true "User data to update"
// @Success 200 {object} UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/school/users/{id} [put]
func (h *Handler) UpdateUser(c *fiber.Ctx) error {
	schoolID, ok := middleware.GetTenantID(c)
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
				"message": "ID user tidak valid",
			},
		})
	}

	var req UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Format data tidak valid",
			},
		})
	}

	response, err := h.service.UpdateUser(c.Context(), schoolID, uint(id), req)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
		"message": "User berhasil diperbarui",
	})
}

// DeleteUser handles deleting a user
// @Summary Delete a user
// @Description Delete a user from the school
// @Tags Users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/school/users/{id} [delete]
func (h *Handler) DeleteUser(c *fiber.Ctx) error {
	schoolID, ok := middleware.GetTenantID(c)
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
				"message": "ID user tidak valid",
			},
		})
	}

	if err := h.service.DeleteUser(c.Context(), schoolID, uint(id)); err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "User berhasil dihapus",
	})
}

// ResetUserPassword handles resetting a user's password
// @Summary Reset user password
// @Description Reset a user's password and generate a new temporary password
// @Tags Users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} ResetPasswordResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/school/users/{id}/reset-password [post]
func (h *Handler) ResetUserPassword(c *fiber.Ctx) error {
	schoolID, ok := middleware.GetTenantID(c)
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
				"message": "ID user tidak valid",
			},
		})
	}

	response, err := h.service.ResetUserPassword(c.Context(), schoolID, uint(id))
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
		"message": "Password berhasil direset",
	})
}

// GetSchoolDevices handles getting all devices for the school
// @Summary Get school devices
// @Description Get all RFID devices registered for this school
// @Tags School
// @Produce json
// @Success 200 {object} []DeviceResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/school/devices [get]
func (h *Handler) GetSchoolDevices(c *fiber.Ctx) error {
	schoolID, ok := middleware.GetTenantID(c)
	if !ok {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTHZ_TENANT_REQUIRED",
				"message": "Konteks sekolah diperlukan",
			},
		})
	}

	devices, err := h.service.GetSchoolDevices(c.Context(), schoolID)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    devices,
	})
}

// ClearStudentRFID handles clearing a student's RFID code
// @Summary Clear student RFID
// @Description Clear the RFID code from a student (unpair the card)
// @Tags Students
// @Produce json
// @Param id path int true "Student ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/school/students/{id}/clear-rfid [post]
func (h *Handler) ClearStudentRFID(c *fiber.Ctx) error {
	schoolID, ok := middleware.GetTenantID(c)
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
				"message": "ID siswa tidak valid",
			},
		})
	}

	if err := h.service.ClearStudentRFID(c.Context(), schoolID, uint(id)); err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Kartu RFID berhasil dihapus dari siswa",
	})
}

// ==================== Bulk Operations Handlers ====================

// GetStudentsWithoutClass handles getting students without class assignment
// @Summary Get students without class
// @Description Get all students that don't have a class assigned (for bulk assignment after import)
// @Tags Students
// @Produce json
// @Success 200 {object} []StudentResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/school/students/without-class [get]
func (h *Handler) GetStudentsWithoutClass(c *fiber.Ctx) error {
	schoolID, ok := middleware.GetTenantID(c)
	if !ok {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTHZ_TENANT_REQUIRED",
				"message": "Konteks sekolah diperlukan",
			},
		})
	}

	students, err := h.service.GetStudentsWithoutClass(c.Context(), schoolID)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    students,
	})
}

// BulkAssignClass handles bulk class assignment for multiple students
// @Summary Bulk assign class to students
// @Description Assign a class to multiple students at once (for students imported without class)
// @Tags Students
// @Accept json
// @Produce json
// @Param request body BulkAssignClassRequest true "Bulk assignment data"
// @Success 200 {object} BulkAssignClassResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/school/students/bulk-assign-class [post]
func (h *Handler) BulkAssignClass(c *fiber.Ctx) error {
	schoolID, ok := middleware.GetTenantID(c)
	if !ok {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTHZ_TENANT_REQUIRED",
				"message": "Konteks sekolah diperlukan",
			},
		})
	}

	var req BulkAssignClassRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Format data tidak valid",
			},
		})
	}

	if len(req.StudentIDs) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "Minimal satu ID siswa wajib diisi",
			},
		})
	}

	if req.ClassID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "ID kelas wajib diisi",
			},
		})
	}

	response, err := h.service.BulkAssignClass(c.Context(), schoolID, req)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
		"message": response.Message,
	})
}

// SearchStudents handles searching students by NISN or name for parent linking
// @Summary Search students
// @Description Search students by NISN or name for linking to parents
// @Tags Students
// @Produce json
// @Param q query string true "Search query (NISN or name)"
// @Success 200 {object} []StudentSearchResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/school/students/search [get]
func (h *Handler) SearchStudents(c *fiber.Ctx) error {
	schoolID, ok := middleware.GetTenantID(c)
	if !ok {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTHZ_TENANT_REQUIRED",
				"message": "Konteks sekolah diperlukan",
			},
		})
	}

	query := c.Query("q")
	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "Parameter pencarian (q) wajib diisi",
			},
		})
	}

	students, err := h.service.SearchStudents(c.Context(), schoolID, query)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    students,
	})
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Success bool `json:"success"`
	Error   struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}
