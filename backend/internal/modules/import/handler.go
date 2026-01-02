package importmodule

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"github.com/school-management/backend/internal/middleware"
)

// Handler handles HTTP requests for import operations
type Handler struct {
	service Service
}

// NewHandler creates a new import handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers import routes
// Requirements: 1.1, 1.2, 2.1
func (h *Handler) RegisterRoutes(router fiber.Router) {
	importGroup := router.Group("/import")

	// Template download routes
	// Requirements: 1.1, 1.2
	importGroup.Get("/template/students", h.DownloadStudentTemplate)
	importGroup.Get("/template/parents", h.DownloadParentTemplate)

	// Import routes
	// Requirements: 2.1
	importGroup.Post("/students", h.ImportStudents)
	importGroup.Post("/parents", h.ImportParents)
}

// DownloadStudentTemplate handles downloading the student import template
// @Summary Download student import template
// @Description Download an Excel template for student bulk import
// @Tags Import
// @Produce application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Success 200 {file} file "Excel template file"
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/school/import/template/students [get]
func (h *Handler) DownloadStudentTemplate(c *fiber.Ctx) error {
	_, ok := middleware.GetTenantID(c)
	if !ok {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTHZ_TENANT_REQUIRED",
				"message": "Konteks sekolah diperlukan",
			},
		})
	}

	data, err := h.service.GenerateStudentTemplate()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "INTERNAL_ERROR",
				"message": "Gagal membuat template",
			},
		})
	}

	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", "attachment; filename=template_import_siswa.xlsx")

	return c.Send(data)
}

// DownloadParentTemplate handles downloading the parent import template
// @Summary Download parent import template
// @Description Download an Excel template for parent bulk import
// @Tags Import
// @Produce application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Success 200 {file} file "Excel template file"
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/school/import/template/parents [get]
func (h *Handler) DownloadParentTemplate(c *fiber.Ctx) error {
	_, ok := middleware.GetTenantID(c)
	if !ok {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTHZ_TENANT_REQUIRED",
				"message": "Konteks sekolah diperlukan",
			},
		})
	}

	data, err := h.service.GenerateParentTemplate()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "INTERNAL_ERROR",
				"message": "Gagal membuat template",
			},
		})
	}

	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", "attachment; filename=template_import_orangtua.xlsx")

	return c.Send(data)
}

// ImportStudents handles student bulk import from Excel file
// @Summary Import students from Excel
// @Description Import multiple students from an Excel file
// @Tags Import
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Excel file (.xlsx)"
// @Success 200 {object} ImportResult
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/school/import/students [post]
func (h *Handler) ImportStudents(c *fiber.Ctx) error {
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

	// Get uploaded file
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_FILE_REQUIRED",
				"message": "File wajib diunggah",
			},
		})
	}

	// Open file
	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_FILE_INVALID",
				"message": "Gagal membuka file",
			},
		})
	}
	defer file.Close()

	// Import students
	result, err := h.service.ImportStudents(c.Context(), schoolID, file, fileHeader.Size)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    result,
		"message": "Import selesai",
	})
}

// ImportParents handles parent bulk import from Excel file
// @Summary Import parents from Excel
// @Description Import multiple parents from an Excel file
// @Tags Import
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Excel file (.xlsx)"
// @Success 200 {object} ImportResult
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/school/import/parents [post]
func (h *Handler) ImportParents(c *fiber.Ctx) error {
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

	// Get uploaded file
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_FILE_REQUIRED",
				"message": "File wajib diunggah",
			},
		})
	}

	// Open file
	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VAL_FILE_INVALID",
				"message": "Gagal membuka file",
			},
		})
	}
	defer file.Close()

	// Import parents
	result, err := h.service.ImportParents(c.Context(), schoolID, file, fileHeader.Size)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    result,
		"message": "Import selesai",
	})
}

// handleError handles errors and returns appropriate HTTP responses
func (h *Handler) handleError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, ErrInvalidFileFormat):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "INVALID_FILE_FORMAT",
				"message": err.Error(),
			},
		})
	case errors.Is(err, ErrFileTooLarge):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "FILE_TOO_LARGE",
				"message": err.Error(),
			},
		})
	case errors.Is(err, ErrEmptyFile):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "EMPTY_FILE",
				"message": err.Error(),
			},
		})
	case errors.Is(err, ErrInvalidHeader):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "INVALID_HEADER",
				"message": err.Error(),
			},
		})
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "INTERNAL_ERROR",
				"message": "Terjadi kesalahan internal",
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
