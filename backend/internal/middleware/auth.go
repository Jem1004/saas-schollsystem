package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/school-management/backend/internal/modules/auth"
)

// AuthMiddleware creates a middleware that validates JWT tokens
// Requirements: 4.5 - THE System SHALL enforce role-based access control for all protected resources
func AuthMiddleware(jwtManager *auth.JWTManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    "AUTH_TOKEN_MISSING",
					"message": "Authorization header is required",
				},
			})
		}

		// Check Bearer prefix
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    "AUTH_TOKEN_INVALID",
					"message": "Invalid authorization header format",
				},
			})
		}

		tokenString := parts[1]

		// Validate token
		claims, err := jwtManager.ValidateAccessToken(tokenString)
		if err != nil {
			return handleTokenError(c, err)
		}

		// Store claims in context for use by handlers
		// Using both camelCase and snake_case for backward compatibility
		c.Locals("userID", claims.UserID)
		c.Locals("user_id", claims.UserID)
		c.Locals("schoolID", claims.SchoolID)
		c.Locals("role", claims.Role)
		c.Locals("username", claims.Username)
		c.Locals("claims", claims)

		return c.Next()
	}
}

// OptionalAuthMiddleware creates a middleware that validates JWT tokens if present
// but allows requests without tokens to proceed
func OptionalAuthMiddleware(jwtManager *auth.JWTManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Next()
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return c.Next()
		}

		tokenString := parts[1]
		claims, err := jwtManager.ValidateAccessToken(tokenString)
		if err != nil {
			return c.Next()
		}

		c.Locals("userID", claims.UserID)
		c.Locals("user_id", claims.UserID)
		c.Locals("schoolID", claims.SchoolID)
		c.Locals("role", claims.Role)
		c.Locals("username", claims.Username)
		c.Locals("claims", claims)

		return c.Next()
	}
}

// handleTokenError handles JWT token validation errors
func handleTokenError(c *fiber.Ctx, err error) error {
	switch err {
	case auth.ErrTokenExpired:
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTH_TOKEN_EXPIRED",
				"message": "Token has expired",
			},
		})
	case auth.ErrTokenMalformed:
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTH_TOKEN_MALFORMED",
				"message": "Token is malformed",
			},
		})
	default:
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "AUTH_TOKEN_INVALID",
				"message": "Invalid token",
			},
		})
	}
}

// GetUserID extracts user ID from context
func GetUserID(c *fiber.Ctx) (uint, bool) {
	userID, ok := c.Locals("userID").(uint)
	return userID, ok
}

// GetSchoolID extracts school ID from context
func GetSchoolID(c *fiber.Ctx) (*uint, bool) {
	schoolID, ok := c.Locals("schoolID").(*uint)
	return schoolID, ok
}

// GetRole extracts role from context
func GetRole(c *fiber.Ctx) (string, bool) {
	role, ok := c.Locals("role").(string)
	return role, ok
}

// GetClaims extracts all claims from context
func GetClaims(c *fiber.Ctx) (*auth.TokenClaims, bool) {
	claims, ok := c.Locals("claims").(*auth.TokenClaims)
	return claims, ok
}
