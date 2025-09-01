package middleware

import (
	userModel "runs-system-user-go/module/user/model"
	"runs-system-user-go/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func AuthMiddleware(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get header
		authHeader := c.Get("Authorization")
		if !strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
			return c.Status(fiber.StatusUnauthorized).
				JSON(fiber.Map{"error": "missing or invalid Authorization header"})
		}

		// Extract token
		tokenString := strings.TrimSpace(authHeader[len("Bearer "):])

		// Parse JWT
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).
				JSON(fiber.Map{"error": "invalid or expired token"})
		}

		// Load user from DB
		var user userModel.User
		if err := db.First(&user, claims.UserID).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).
				JSON(fiber.Map{"error": "user not found"})
		}

		// Attach user to context
		c.Locals("user", &user)

		// Continue
		return c.Next()
	}
}
