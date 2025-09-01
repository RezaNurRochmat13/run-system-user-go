package routes

import (
	userRoutes "runs-system-user-go/module/user/routes"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v1")

	userRoutes.SetupUserRoutes(api)
}
