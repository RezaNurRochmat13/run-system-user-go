package authController

import (
	authModel "runs-system-user-go/module/auth/model"
	authService "runs-system-user-go/module/auth/service"

	"github.com/gofiber/fiber/v2"
)

func Register(c *fiber.Ctx) error {
	var userRegister authModel.UserRegister
	if err := c.BodyParser(&userRegister); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Invalid input", "data": nil})
	}

	if err := authService.RegisterUser(&userRegister); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to register user", "data": nil})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "User registered successfully"})
}

func Login(c *fiber.Ctx) error {
	var userLogin authModel.UserLogin
	if err := c.BodyParser(&userLogin); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Invalid input", "data": nil})
	}

	authResponse, err := authService.LoginUser(&userLogin)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to login user", "data": nil})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "User logged in successfully", "data": authResponse})
}
