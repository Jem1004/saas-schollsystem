package middleware

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/school-management/backend/internal/domain/models"
)

// TenantMiddleware creates a middleware that enforces tenant isolation
// Requirements: 1.4 - WHEN data is queried, THE System SHALL filter results by school_id to ensure tenant isolation
// Requirements: 1.5 - IF a request attempts to access data from another tenant, THEN THE System SHALL reject the request
func TenantMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, _ := c.Locals("role").(string)
		username, _ := c.Locals("username").(string)
		userID, _ := c.Locals("userID").(uint)

		// Super admin can access all tenants - but they need to specify which school
		// For settings, super admin should not access without school context
		if role == string(models.RoleSuperAdmin) {
			// Check if there's a school_id in query or header for super admin
			// For now, just pass through - handler will check for tenantID
			return c.Next()
		}

		// Get school ID from context (set by auth middleware)
		schoolIDVal := c.Locals("schoolID")
		
		// Check if schoolID exists and is not nil
		var schoolID *uint
		if schoolIDVal != nil {
			switch v := schoolIDVal.(type) {
			case *uint:
				if v != nil {
					schoolID = v
				}
			case uint:
				schoolID = &v
			}
		}

		// For non-super_admin users, school_id must be present
		if schoolID == nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    "AUTHZ_TENANT_REQUIRED",
					"message": "Konteks sekolah diperlukan. Pastikan akun Anda terhubung dengan sekolah.",
					"debug": fiber.Map{
						"user_id":  userID,
						"username": username,
						"role":     role,
						"hint":     "User tidak memiliki school_id di token JWT. Silakan login ulang atau hubungi administrator.",
					},
				},
			})
		}

		// Store tenant ID for use in queries
		// Using both tenantID and school_id for backward compatibility
		c.Locals("tenantID", *schoolID)
		c.Locals("school_id", *schoolID)

		return c.Next()
	}
}

// TenantScope returns a GORM scope that filters by tenant
// Requirements: 1.4 - Data queries SHALL filter by school_id
func TenantScope(schoolID uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("school_id = ?", schoolID)
	}
}

// TenantScopeFromContext returns a GORM scope using the tenant ID from context
func TenantScopeFromContext(c *fiber.Ctx) func(db *gorm.DB) *gorm.DB {
	tenantID, ok := c.Locals("tenantID").(uint)
	if !ok {
		// Return a scope that matches nothing if no tenant ID
		return func(db *gorm.DB) *gorm.DB {
			return db.Where("1 = 0")
		}
	}
	return TenantScope(tenantID)
}

// GetTenantID extracts tenant ID from context
func GetTenantID(c *fiber.Ctx) (uint, bool) {
	tenantID, ok := c.Locals("tenantID").(uint)
	return tenantID, ok
}

// IsSuperAdmin checks if the current user is a super admin
func IsSuperAdmin(c *fiber.Ctx) bool {
	role, ok := c.Locals("role").(string)
	return ok && role == string(models.RoleSuperAdmin)
}

// ValidateTenantAccess validates that the user can access the specified tenant
// Requirements: 1.5 - Cross-tenant access SHALL be rejected
func ValidateTenantAccess(c *fiber.Ctx, targetSchoolID uint) error {
	// Super admin can access any tenant
	if IsSuperAdmin(c) {
		return nil
	}

	tenantID, ok := GetTenantID(c)
	if !ok {
		return fiber.NewError(fiber.StatusForbidden, "Konteks sekolah diperlukan")
	}

	if tenantID != targetSchoolID {
		return fiber.NewError(fiber.StatusForbidden, "Akses ke sekolah ini tidak diizinkan")
	}

	return nil
}

// TenantContext holds tenant information for use in services
type TenantContext struct {
	SchoolID    uint
	IsSuperAdmin bool
}

// GetTenantContext extracts tenant context from fiber context
func GetTenantContext(c *fiber.Ctx) *TenantContext {
	role, _ := c.Locals("role").(string)
	isSuperAdmin := role == string(models.RoleSuperAdmin)

	var schoolID uint
	if tenantID, ok := c.Locals("tenantID").(uint); ok {
		schoolID = tenantID
	} else if sid, ok := c.Locals("schoolID").(*uint); ok && sid != nil {
		schoolID = *sid
	}

	return &TenantContext{
		SchoolID:    schoolID,
		IsSuperAdmin: isSuperAdmin,
	}
}

// WithTenantContext adds tenant context to a context.Context
func WithTenantContext(ctx context.Context, tc *TenantContext) context.Context {
	return context.WithValue(ctx, tenantContextKey, tc)
}

// TenantContextFromContext extracts tenant context from context.Context
func TenantContextFromContext(ctx context.Context) *TenantContext {
	tc, ok := ctx.Value(tenantContextKey).(*TenantContext)
	if !ok {
		return nil
	}
	return tc
}

type contextKey string

const tenantContextKey contextKey = "tenantContext"
