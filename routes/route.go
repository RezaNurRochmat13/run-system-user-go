package routes

import (
	todoRoutes "todolist-ilcs-api/module/todo/routes"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v1")

	todoRoutes.SetupTodoRoutes(api)
}
