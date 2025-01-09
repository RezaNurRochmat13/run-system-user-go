package todoController

import (
	"todolist-ilcs-api/database"
	todoModel "todolist-ilcs-api/module/todo/model"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetAllTodos(c *fiber.Ctx) error {
	db := database.DB
	var todos []todoModel.Todo

	db.Find(&todos)

	if len(todos) == 0 {
		return c.Status(200).JSON(fiber.Map{"status": "error", "message": "No todos present", "data": nil})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Todos fetched successfully", "data": todos})
}

func GetSingleTodo(c *fiber.Ctx) error {
    db := database.DB
    var todo todoModel.Todo

    // Read the param noteId
    id := c.Params("id")

    // Find the note with the given Id
    db.Find(&todo, "id = ?", id)

    // If no such note present return an error
    if todo.ID == uuid.Nil {
        return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No note present", "data": nil})
    }

    // Return the note with the Id
    return c.JSON(fiber.Map{"status": "success", "message": "Notes Found", "data": todo})
}


func CreateNewTodo(c *fiber.Ctx) error {
    db := database.DB
    todo := new(todoModel.Todo)

    // Store the body in the note and return error if encountered
    err := c.BodyParser(todo)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
    }

    // Add a uuid to the note
    todo.ID = uuid.New()
    // Create the Note and return error if encountered
    err = db.Create(&todo).Error
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not create note", "data": err})
    }

    // Return the created note
    return c.JSON(fiber.Map{"status": "success", "message": "Created Note", "data": todo})
}

func UpdateTodo(c *fiber.Ctx) error {
    type updateTodoPayload struct {
        Title    string `json:"title"`
        Description string `json:"description"`
        Status     string `json:"status"`
        DueDate    string `json:"due_date"`
    }

    db := database.DB
    var todo todoModel.Todo

    // Read the param noteId
    id := c.Params("id")

    // Find the note with the given Id
    db.Find(&todo, "id = ?", id)

    // If no such note present return an error
    if todo.ID == uuid.Nil {
        return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No note present", "data": nil})
    }

    // Store the body containing the updated data and return error if encountered
    var updateNoteData updateTodoPayload
    err := c.BodyParser(&updateNoteData)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
    }

    // Edit the note
    todo.Title = updateNoteData.Title
    todo.Description = updateNoteData.Description
    todo.Status = updateNoteData.Status
    todo.DueDate = updateNoteData.DueDate

    // Save the Changes
    db.Save(&todo)

    // Return the updated note
    return c.JSON(fiber.Map{"status": "success", "message": "Notes Found", "data": todo})
}

func DeleteTodo(c *fiber.Ctx) error {
    db := database.DB
    var todo todoModel.Todo

    // Read the param noteId
    id := c.Params("noteId")

    // Find the note with the given Id
    db.Find(&todo, "id = ?", id)

    // If no such note present return an error
    if todo.ID == uuid.Nil {
        return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No note present", "data": nil})
    }

    // Delete the note and return error if encountered
    err := db.Delete(&todo, "id = ?", id).Error

    if err != nil {
        return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Failed to delete note", "data": nil})
    }

    // Return success message
    return c.JSON(fiber.Map{"status": "success", "message": "Deleted Note"})
}


