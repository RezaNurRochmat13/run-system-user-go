package test

import (
	"os"
	"runs-system-user-go/database"
	userRoutes "runs-system-user-go/module/user/routes"

	"github.com/gofiber/fiber/v2"
)

var db = database.DB

func SetupTestApp() *fiber.App {
	os.Setenv("DB_NAME", "userlist_test") // Use test database
	database.ConnectDatabase()
	db.Exec("TRUNCATE users RESTART IDENTITY") // Clear test database
	app := fiber.New()
	userRoutes.SetupUserRoutes(app)
	return app
}
