package main

import (
	"runs-system-user-go/database"
	"runs-system-user-go/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Start a new Fiber App
	app := fiber.New()

	app.Use(logger.New())

	// Connect to the database
	database.ConnectDatabase()

	// Connect to Redis
	database.SetupConnectRedis()

	// Setup Routes
	routes.SetupRoutes(app)


	// Send string back for GET calls to the endpoint '/'
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("API is up and running")
	})

	// Listen on port 3000
	app.Listen(":3000")
}
