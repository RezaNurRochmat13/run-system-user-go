package userRoutes

import (
	"runs-system-user-go/middleware"
	userController "runs-system-user-go/module/user/controller"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupUserRoutes(db *gorm.DB, router fiber.Router) {
	userRoutes := router.Group("/users")

	// Read all Users
	userRoutes.Get("/", middleware.AuthMiddleware(db), userController.GetAllUsers)

	// Create new User
	userRoutes.Post("/", userController.CreateNewUser)

	// Get User by ID
	userRoutes.Get("/:id", userController.GetSingleUser)

	// Update User by ID
	userRoutes.Put("/:id", userController.UpdateUser)

	// Delete User by ID
	userRoutes.Delete("/:id", userController.DeleteUser)
}
