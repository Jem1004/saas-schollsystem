package middleware

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/school-management/backend/internal/domain/models"
	"github.com/school-management/backend/internal/policy"
)

// RBACConfig holds configuration for RBAC middleware
type RBACConfig struct {
	AccessPolicy policy.AccessPolicy
}

// RoleMiddleware creates a middleware that restricts access to specific roles
// Requirements: 4.5 - THE System SHALL enforce role-based access control for all protected resources
func RoleMiddleware(allowedRoles ...models.UserRole) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, ok := c.Locals("role").(string)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    "AUTHZ_ROLE_DENIED",
					"message": "Access denied",
				},
			})
		}

		userRole := models.UserRole(role)
		for _, allowedRole := range allowedRoles {
			if userRole == allowedRole {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTHZ_ROLE_DENIED",
				"message": "You do not have permission to access this resource",
			},
		})
	}
}

// GetUserContext extracts user context from fiber context for access policy checks
func GetUserContext(c *fiber.Ctx) *policy.UserContext {
	userID, _ := c.Locals("userID").(uint)
	schoolID, _ := c.Locals("schoolID").(*uint)
	role, _ := c.Locals("role").(string)
	classID, _ := c.Locals("assignedClassID").(*uint)

	return &policy.UserContext{
		UserID:   userID,
		SchoolID: schoolID,
		Role:     models.UserRole(role),
		ClassID:  classID,
	}
}

// SetUserClassContext middleware that loads the assigned class for wali kelas users
func SetUserClassContext(accessPolicy policy.AccessPolicy) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, ok := c.Locals("role").(string)
		if !ok {
			return c.Next()
		}

		// Only load class context for wali kelas
		if models.UserRole(role) == models.RoleWaliKelas {
			userID, ok := c.Locals("userID").(uint)
			if ok {
				classID, err := accessPolicy.GetUserAssignedClassID(c.Context(), userID)
				if err == nil && classID != nil {
					c.Locals("assignedClassID", classID)
				}
			}
		}

		return c.Next()
	}
}

// StudentAccessMiddleware validates access to student resources
// Requirements: 4.5 - Role-based access control for student data
func StudentAccessMiddleware(accessPolicy policy.AccessPolicy) fiber.Handler {
	return func(c *fiber.Ctx) error {
		studentIDStr := c.Params("studentId")
		if studentIDStr == "" {
			return c.Next()
		}

		studentID, err := strconv.ParseUint(studentIDStr, 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    "VAL_INVALID_FORMAT",
					"message": "Invalid student ID format",
				},
			})
		}

		userCtx := GetUserContext(c)
		canAccess, err := accessPolicy.CanAccessStudent(c.Context(), userCtx, uint(studentID))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    "INTERNAL_ERROR",
					"message": "Failed to validate access",
				},
			})
		}

		if !canAccess {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    "AUTHZ_STUDENT_ACCESS_DENIED",
					"message": "You do not have permission to access this student's data",
				},
			})
		}

		return c.Next()
	}
}

// SuperAdminOnly restricts access to super admin only
func SuperAdminOnly() fiber.Handler {
	return RoleMiddleware(models.RoleSuperAdmin)
}

// AdminSekolahOnly restricts access to school admin only
func AdminSekolahOnly() fiber.Handler {
	return RoleMiddleware(models.RoleAdminSekolah)
}

// GuruBKOnly restricts access to counseling teacher only
func GuruBKOnly() fiber.Handler {
	return RoleMiddleware(models.RoleGuruBK)
}

// WaliKelasOnly restricts access to homeroom teacher only
func WaliKelasOnly() fiber.Handler {
	return RoleMiddleware(models.RoleWaliKelas)
}

// SchoolStaffOnly restricts access to school staff (admin, teachers)
func SchoolStaffOnly() fiber.Handler {
	return RoleMiddleware(
		models.RoleAdminSekolah,
		models.RoleGuruBK,
		models.RoleWaliKelas,
		models.RoleGuru,
	)
}

// TeachersOnly restricts access to teachers (BK, homeroom, regular)
func TeachersOnly() fiber.Handler {
	return RoleMiddleware(
		models.RoleGuruBK,
		models.RoleWaliKelas,
		models.RoleGuru,
	)
}

// AdminOrSuperAdmin restricts access to admin roles only
func AdminOrSuperAdmin() fiber.Handler {
	return RoleMiddleware(
		models.RoleSuperAdmin,
		models.RoleAdminSekolah,
	)
}

// BKWriteAccess restricts write access to BK data to authorized roles
// Requirements: 6.5 - Wali_Kelas has read-only access to BK data
func BKWriteAccess() fiber.Handler {
	return RoleMiddleware(
		models.RoleSuperAdmin,
		models.RoleAdminSekolah,
		models.RoleGuruBK,
	)
}

// GradeWriteAccess restricts write access to grades
// Requirements: 10.5 - Wali_Kelas can only input grades for students in their assigned class
func GradeWriteAccess() fiber.Handler {
	return RoleMiddleware(
		models.RoleSuperAdmin,
		models.RoleAdminSekolah,
		models.RoleWaliKelas,
	)
}

// HomeroomNoteWriteAccess restricts write access to homeroom notes
// Requirements: 11.4 - Wali_Kelas can only create notes for students in their assigned class
func HomeroomNoteWriteAccess() fiber.Handler {
	return RoleMiddleware(
		models.RoleSuperAdmin,
		models.RoleAdminSekolah,
		models.RoleWaliKelas,
	)
}

// BKAccessMiddleware provides access control for BK (counseling) data
// Requirements: 6.5 - WHEN a Wali_Kelas views BK data, THE System SHALL provide read-only access
// Requirements: 9.3 - THE System SHALL keep internal_note private and accessible only to Guru_BK
// Requirements: 9.4 - WHEN a Wali_Kelas views counseling data, THE System SHALL show only parent_summary
func BKAccessMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, ok := c.Locals("role").(string)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    "AUTHZ_ROLE_DENIED",
					"message": "Access denied",
				},
			})
		}

		userRole := models.UserRole(role)

		// Determine access level based on role
		switch userRole {
		case models.RoleSuperAdmin, models.RoleAdminSekolah:
			// Full access for admin roles
			c.Locals("bkAccessLevel", "full")
		case models.RoleGuruBK:
			// Full access including internal notes
			c.Locals("bkAccessLevel", "full")
			c.Locals("canViewInternalNotes", true)
		case models.RoleWaliKelas:
			// Read-only access, no internal notes
			c.Locals("bkAccessLevel", "readonly")
			c.Locals("canViewInternalNotes", false)
		case models.RoleParent:
			// Limited access - only parent summary
			c.Locals("bkAccessLevel", "limited")
			c.Locals("canViewInternalNotes", false)
		case models.RoleStudent:
			// Limited access - summary only
			c.Locals("bkAccessLevel", "limited")
			c.Locals("canViewInternalNotes", false)
		default:
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    "AUTHZ_ROLE_DENIED",
					"message": "You do not have permission to access BK data",
				},
			})
		}

		return c.Next()
	}
}

// BKAccessMiddlewareWithPolicy provides access control for BK data using access policy
func BKAccessMiddlewareWithPolicy(accessPolicy policy.AccessPolicy) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userCtx := GetUserContext(c)
		accessLevel := accessPolicy.CanAccessBKData(userCtx)

		if accessLevel == policy.AccessLevelNone {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    "AUTHZ_ROLE_DENIED",
					"message": "You do not have permission to access BK data",
				},
			})
		}

		c.Locals("bkAccessLevel", string(accessLevel))
		c.Locals("canViewInternalNotes", accessPolicy.CanViewInternalNotes(userCtx))
		c.Locals("canModifyBKData", accessPolicy.CanModifyBKData(userCtx))

		return c.Next()
	}
}

// GradeAccessMiddleware provides access control for grade data
// Requirements: 10.5 - THE System SHALL validate that Wali_Kelas can only input grades for students in their assigned class
func GradeAccessMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, ok := c.Locals("role").(string)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    "AUTHZ_ROLE_DENIED",
					"message": "Access denied",
				},
			})
		}

		userRole := models.UserRole(role)

		switch userRole {
		case models.RoleSuperAdmin, models.RoleAdminSekolah:
			c.Locals("gradeAccessLevel", "full")
		case models.RoleWaliKelas:
			// Can only modify grades for their class
			c.Locals("gradeAccessLevel", "class_only")
		case models.RoleParent, models.RoleStudent:
			c.Locals("gradeAccessLevel", "readonly")
		default:
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    "AUTHZ_ROLE_DENIED",
					"message": "You do not have permission to access grade data",
				},
			})
		}

		return c.Next()
	}
}

// GradeAccessMiddlewareWithPolicy provides access control for grade data using access policy
func GradeAccessMiddlewareWithPolicy(accessPolicy policy.AccessPolicy) fiber.Handler {
	return func(c *fiber.Ctx) error {
		studentIDStr := c.Params("studentId")
		if studentIDStr == "" {
			// No student context, just set basic access level
			userCtx := GetUserContext(c)
			switch userCtx.Role {
			case models.RoleSuperAdmin, models.RoleAdminSekolah:
				c.Locals("gradeAccessLevel", string(policy.AccessLevelFull))
			case models.RoleWaliKelas:
				c.Locals("gradeAccessLevel", "class_only")
			case models.RoleParent, models.RoleStudent:
				c.Locals("gradeAccessLevel", string(policy.AccessLevelReadOnly))
			default:
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"success": false,
					"error": fiber.Map{
						"code":    "AUTHZ_ROLE_DENIED",
						"message": "You do not have permission to access grade data",
					},
				})
			}
			return c.Next()
		}

		studentID, err := strconv.ParseUint(studentIDStr, 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    "VAL_INVALID_FORMAT",
					"message": "Invalid student ID format",
				},
			})
		}

		userCtx := GetUserContext(c)
		accessLevel, err := accessPolicy.CanAccessGrade(c.Context(), userCtx, uint(studentID))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    "INTERNAL_ERROR",
					"message": "Failed to validate access",
				},
			})
		}

		if accessLevel == policy.AccessLevelNone {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    "AUTHZ_GRADE_ACCESS_DENIED",
					"message": "You do not have permission to access this student's grades",
				},
			})
		}

		c.Locals("gradeAccessLevel", string(accessLevel))
		return c.Next()
	}
}

// HomeroomNoteAccessMiddleware provides access control for homeroom notes
// Requirements: 11.4 - THE System SHALL validate that Wali_Kelas can only create notes for students in their assigned class
func HomeroomNoteAccessMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, ok := c.Locals("role").(string)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    "AUTHZ_ROLE_DENIED",
					"message": "Access denied",
				},
			})
		}

		userRole := models.UserRole(role)

		switch userRole {
		case models.RoleSuperAdmin, models.RoleAdminSekolah:
			c.Locals("homeroomNoteAccessLevel", "full")
		case models.RoleWaliKelas:
			c.Locals("homeroomNoteAccessLevel", "class_only")
		case models.RoleGuruBK:
			c.Locals("homeroomNoteAccessLevel", "readonly")
		case models.RoleParent, models.RoleStudent:
			c.Locals("homeroomNoteAccessLevel", "readonly")
		default:
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    "AUTHZ_ROLE_DENIED",
					"message": "You do not have permission to access homeroom notes",
				},
			})
		}

		return c.Next()
	}
}

// HomeroomNoteAccessMiddlewareWithPolicy provides access control for homeroom notes using access policy
func HomeroomNoteAccessMiddlewareWithPolicy(accessPolicy policy.AccessPolicy) fiber.Handler {
	return func(c *fiber.Ctx) error {
		studentIDStr := c.Params("studentId")
		if studentIDStr == "" {
			// No student context, just set basic access level
			userCtx := GetUserContext(c)
			switch userCtx.Role {
			case models.RoleSuperAdmin, models.RoleAdminSekolah:
				c.Locals("homeroomNoteAccessLevel", string(policy.AccessLevelFull))
			case models.RoleWaliKelas:
				c.Locals("homeroomNoteAccessLevel", "class_only")
			case models.RoleGuruBK, models.RoleParent, models.RoleStudent:
				c.Locals("homeroomNoteAccessLevel", string(policy.AccessLevelReadOnly))
			default:
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"success": false,
					"error": fiber.Map{
						"code":    "AUTHZ_ROLE_DENIED",
						"message": "You do not have permission to access homeroom notes",
					},
				})
			}
			return c.Next()
		}

		studentID, err := strconv.ParseUint(studentIDStr, 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    "VAL_INVALID_FORMAT",
					"message": "Invalid student ID format",
				},
			})
		}

		userCtx := GetUserContext(c)
		accessLevel, err := accessPolicy.CanAccessHomeroomNote(c.Context(), userCtx, uint(studentID))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    "INTERNAL_ERROR",
					"message": "Failed to validate access",
				},
			})
		}

		if accessLevel == policy.AccessLevelNone {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    "AUTHZ_HOMEROOM_ACCESS_DENIED",
					"message": "You do not have permission to access this student's homeroom notes",
				},
			})
		}

		c.Locals("homeroomNoteAccessLevel", string(accessLevel))
		return c.Next()
	}
}

// ValidateClassAccess middleware validates that wali kelas can only access their assigned class
// Requirements: 10.5, 11.4 - Wali_Kelas can only access students in their assigned class
func ValidateClassAccess(accessPolicy policy.AccessPolicy) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userCtx := GetUserContext(c)

		// Only apply to wali kelas
		if userCtx.Role != models.RoleWaliKelas {
			return c.Next()
		}

		studentIDStr := c.Params("studentId")
		if studentIDStr == "" {
			return c.Next()
		}

		studentID, err := strconv.ParseUint(studentIDStr, 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    "VAL_INVALID_FORMAT",
					"message": "Invalid student ID format",
				},
			})
		}

		// Check if student is in wali kelas's class
		inClass, err := accessPolicy.IsStudentInUserClass(c.Context(), userCtx, uint(studentID))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    "INTERNAL_ERROR",
					"message": "Failed to validate class access",
				},
			})
		}

		if !inClass {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    "AUTHZ_CLASS_ACCESS_DENIED",
					"message": "You can only access students in your assigned class",
				},
			})
		}

		return c.Next()
	}
}

// CanViewInternalNotes checks if the current user can view internal counseling notes
func CanViewInternalNotes(c *fiber.Ctx) bool {
	canView, ok := c.Locals("canViewInternalNotes").(bool)
	return ok && canView
}

// GetBKAccessLevel returns the BK access level for the current user
func GetBKAccessLevel(c *fiber.Ctx) string {
	level, ok := c.Locals("bkAccessLevel").(string)
	if !ok {
		return ""
	}
	return level
}

// IsBKReadOnly checks if the current user has read-only access to BK data
func IsBKReadOnly(c *fiber.Ctx) bool {
	return GetBKAccessLevel(c) == "readonly"
}

// GetGradeAccessLevel returns the grade access level for the current user
func GetGradeAccessLevel(c *fiber.Ctx) string {
	level, ok := c.Locals("gradeAccessLevel").(string)
	if !ok {
		return ""
	}
	return level
}
