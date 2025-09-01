package userController

import (
	userModel "runs-system-user-go/module/user/model"
	userService "runs-system-user-go/module/user/service"

	"github.com/gofiber/fiber/v2"
)

func GetAllUsers(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	users, _, err := userService.GetPaginatedUsers(page, limit, "", "")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": err.Error(), "data": nil})
	}

	if len(users) == 0 {
		return c.Status(200).JSON(fiber.Map{"status": "error", "message": "No users found", "data": nil})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Users fetched successfully", "data": users, "page": page, "limit": limit})
}

func GetSingleUser(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := userService.GetUserByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": err.Error(), "data": nil})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "User found", "data": user})
}

func CreateNewUser(c *fiber.Ctx) error {
	var user userModel.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Invalid input", "data": nil})
	}

	if err := userService.CreateUser(&user); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": err.Error(), "data": nil})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "User created successfully", "data": user})
}

func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var data map[string]interface{}

	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Invalid input", "data": nil})
	}

	user, err := userService.UpdateUser(id, data)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": err.Error(), "data": nil})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "User updated successfully", "data": user})
}

func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := userService.DeleteUser(id); err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": err.Error(), "data": nil})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "User deleted successfully"})
}
