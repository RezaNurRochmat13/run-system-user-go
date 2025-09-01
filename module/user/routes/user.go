package userRoutes

import (
	userController "runs-system-user-go/module/user/controller"

	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(router fiber.Router) {
	userRoutes := router.Group("/users")
	
	// Read all Users
	userRoutes.Get("/", userController.GetAllUsers)
	
	// Create new User
	userRoutes.Post("/", userController.CreateNewUser)
	
	// Get User by ID
	userRoutes.Get("/:id", userController.GetSingleUser)
	
	// Update User by ID
	userRoutes.Put("/:id", userController.UpdateUser)
	
	// Delete User by ID
	userRoutes.Delete("/:id", userController.DeleteUser)
}
