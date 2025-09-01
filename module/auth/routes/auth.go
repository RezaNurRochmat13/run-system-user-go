package authRoutes

import (
	authController "runs-system-user-go/module/auth/controller"

	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(router fiber.Router) {
	authRoutes := router.Group("/auth")

	authRoutes.Post("/register", authController.Register)
	authRoutes.Post("/login", authController.Login)
}
