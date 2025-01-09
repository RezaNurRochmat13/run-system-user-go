package todoRoutes

import (
	todoController "todolist-ilcs-api/module/todo/controller"

	"github.com/gofiber/fiber/v2"
)

func SetupTodoRoutes(router fiber.Router) {
	todo := router.Group("/todos")
	
	// Read all Todos
	todo.Get("/", todoController.GetAllTodos)
	// Read one Todo
	todo.Get("/:id", todoController.GetSingleTodo)
	// Create a Todo
	todo.Post("/", todoController.CreateNewTodo)
	// Update one Todo
	todo.Put("/:id", todoController.UpdateTodo)
	// Delete one Todo
	todo.Delete("/:id", todoController.DeleteTodo)
}
