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
	students.Get("/:id", h.GetStudent)
	students.Put("/:id", h.UpdateStudent)
	students.Delete("/:id", h.DeleteStudent)

	// Parent routes
	parents := router.Group("/parents")
	parents.Post("", h.CreateParent)
	parents.Get("", h.GetParents)
	parents.Get("/:id", h.GetParent)
	parents.Put("/:id", h.UpdateParent)
	parents.Delete("/:id", h.DeleteParent)
	parents.Post("/:id/students", h.LinkParentToStudents)
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
				"message": "Tenant context is required",
			},
		})
	}

	var req CreateClassRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Invalid request body",
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
		"message": "Class created successfully",
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
				"message": "Tenant context is required",
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
				"message": "Tenant context is required",
			},
		})
	}

	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Invalid class ID",
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
				"message": "Tenant context is required",
			},
		})
	}

	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Invalid class ID",
			},
		})
	}

	var req UpdateClassRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Invalid request body",
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
		"message": "Class updated successfully",
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
				"message": "Tenant context is required",
			},
		})
	}

	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Invalid class ID",
			},
		})
	}

	if err := h.service.DeleteClass(c.Context(), schoolID, uint(id)); err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Class deleted successfully",
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
				"message": "Tenant context is required",
			},
		})
	}

	classID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Invalid class ID",
			},
		})
	}

	var req AssignHomeroomTeacherRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Invalid request body",
			},
		})
	}

	if req.TeacherID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "Teacher ID is required",
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
		"message": "Homeroom teacher assigned successfully",
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
				"message": "Tenant context is required",
			},
		})
	}

	classID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Invalid class ID",
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
				"message": "Tenant context is required",
			},
		})
	}

	var req CreateStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Invalid request body",
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
		"message": "Student created successfully",
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
				"message": "Tenant context is required",
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
				"message": "Tenant context is required",
			},
		})
	}

	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Invalid student ID",
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
				"message": "Tenant context is required",
			},
		})
	}

	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Invalid student ID",
			},
		})
	}

	var req UpdateStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Invalid request body",
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
		"message": "Student updated successfully",
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
				"message": "Tenant context is required",
			},
		})
	}

	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Invalid student ID",
			},
		})
	}

	if err := h.service.DeleteStudent(c.Context(), schoolID, uint(id)); err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Student deleted successfully",
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
				"message": "Tenant context is required",
			},
		})
	}

	var req CreateParentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Invalid request body",
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
		"message": "Parent created successfully",
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
				"message": "Tenant context is required",
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
				"message": "Tenant context is required",
			},
		})
	}

	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Invalid parent ID",
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
				"message": "Tenant context is required",
			},
		})
	}

	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Invalid parent ID",
			},
		})
	}

	var req UpdateParentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Invalid request body",
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
		"message": "Parent updated successfully",
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
				"message": "Tenant context is required",
			},
		})
	}

	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Invalid parent ID",
			},
		})
	}

	if err := h.service.DeleteParent(c.Context(), schoolID, uint(id)); err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Parent deleted successfully",
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
				"message": "Tenant context is required",
			},
		})
	}

	parentID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Invalid parent ID",
			},
		})
	}

	var req LinkParentStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Invalid request body",
			},
		})
	}

	if len(req.StudentIDs) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "At least one student ID is required",
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
		"message": "Parent linked to students successfully",
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
				"message": "Class not found",
			},
		})
	case errors.Is(err, ErrDuplicateClass):
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_DUPLICATE_ENTRY",
				"message": "A class with this name already exists for this grade and year",
			},
		})
	case errors.Is(err, ErrClassHasStudents):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_CONSTRAINT_VIOLATION",
				"message": "Cannot delete class with students. Please reassign or remove students first.",
			},
		})

	// Student errors
	case errors.Is(err, ErrStudentNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "NOT_FOUND_STUDENT",
				"message": "Student not found",
			},
		})
	case errors.Is(err, ErrDuplicateNISN):
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_DUPLICATE_NISN",
				"message": "A student with this NISN already exists in the system",
			},
		})
	case errors.Is(err, ErrDuplicateNIS):
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_DUPLICATE_NIS",
				"message": "A student with this NIS already exists in this school",
			},
		})

	// Parent errors
	case errors.Is(err, ErrParentNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "NOT_FOUND_PARENT",
				"message": "Parent not found",
			},
		})

	// Teacher errors
	case errors.Is(err, ErrTeacherNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "NOT_FOUND_TEACHER",
				"message": "Teacher not found",
			},
		})
	case errors.Is(err, ErrInvalidTeacher):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_TEACHER",
				"message": "User is not a valid teacher for homeroom assignment",
			},
		})

	// Validation errors
	case errors.Is(err, ErrNameRequired):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "Name is required",
			},
		})
	case errors.Is(err, ErrNISRequired):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "NIS is required",
			},
		})
	case errors.Is(err, ErrNISNRequired):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "NISN is required",
			},
		})
	case errors.Is(err, ErrClassIDRequired):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "Class ID is required",
			},
		})
	case errors.Is(err, ErrGradeRequired):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "Grade is required and must be greater than 0",
			},
		})
	case errors.Is(err, ErrYearRequired):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "Year is required",
			},
		})
	case errors.Is(err, ErrPasswordRequired):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "Password is required",
			},
		})
	case errors.Is(err, ErrPasswordTooShort):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_INVALID_FORMAT",
				"message": "Password must be at least 8 characters",
			},
		})
	case errors.Is(err, ErrStudentIDsRequired):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_REQUIRED_FIELD",
				"message": "At least one student ID is required",
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
