package todoController

import (
	todoModel "todolist-ilcs-api/module/todo/model"
	todoService "todolist-ilcs-api/module/todo/service"

	"github.com/gofiber/fiber/v2"
)

func GetAllTodos(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	status := c.Query("status")
	search := c.Query("search", "")

	todos, _, err := todoService.GetPaginatedTodos(page, limit, status, search)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": err.Error(), "data": nil})
	}

	if len(todos) == 0 {
		return c.Status(200).JSON(fiber.Map{"status": "error", "message": "No todos found", "data": nil})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Todos fetched successfully", "data": todos, "page": page, "limit": limit})
}

func GetSingleTodo(c *fiber.Ctx) error {
	id := c.Params("id")

	todo, err := todoService.GetTodoByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": err.Error(), "data": nil})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Todo found", "data": todo})
}

func CreateNewTodo(c *fiber.Ctx) error {
	var todo todoModel.Todo
	if err := c.BodyParser(&todo); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Invalid input", "data": nil})
	}

	if err := todoService.CreateTodo(&todo); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": err.Error(), "data": nil})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Todo created successfully", "data": todo})
}

func UpdateTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	var data map[string]interface{}

	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Invalid input", "data": nil})
	}

	todo, err := todoService.UpdateTodo(id, data)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": err.Error(), "data": nil})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Todo updated successfully", "data": todo})
}

func DeleteTodo(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := todoService.DeleteTodo(id); err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": err.Error(), "data": nil})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Todo deleted successfully"})
}
