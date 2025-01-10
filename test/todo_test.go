package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"todolist-ilcs-api/database"
	todoModel "todolist-ilcs-api/module/todo/model"
	todoRoutes "todolist-ilcs-api/module/todo/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

var db = database.DB

func setupTestApp() *fiber.App {
	os.Setenv("DB_NAME", "todolist_test") // Use test database
	database.ConnectDatabase()
	db.Exec("TRUNCATE todos RESTART IDENTITY") // Clear test database
	app := fiber.New()
	todoRoutes.SetupTodoRoutes(app)
	return app
}

func TestTodoCRUD(t *testing.T) {
	app := setupTestApp()

	t.Run("Create Todo", func(t *testing.T) {
		todo := map[string]interface{}{"title": "Test Todo","description": "Test Description", "status": "todo", "due_date": "2023-01-01"}
		body, _ := json.Marshal(todo)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/todos", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var result todoModel.Todo
		json.NewDecoder(resp.Body).Decode(&result)
		assert.Equal(t, "Test Todo", result.Title)
	})

	t.Run("Get Todos with Pagination, Search, and Filter", func(t *testing.T) {
		// Insert mock data into the database
		mockTodos := []todoModel.Todo{
			{Title: "Test Todo 1", Description: "Important task", Status: "pending"},
			{Title: "Test Todo 2", Description: "Another important task", Status: "completed"},
			{Title: "Urgent Task", Description: "Needs attention", Status: "pending"},
		}
		for _, todo := range mockTodos {
			db.Create(&todo) // Ensure the database is set up and accessible
		}
	
		// Test fetching with pagination
		req := httptest.NewRequest(http.MethodGet, "/api/v1/todos?page=1&limit=2", nil)
		resp, _ := app.Test(req)
	
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	
		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		data, _ := json.Marshal(response["data"]) // Extract "data" field
		var todos []todoModel.Todo
		json.Unmarshal(data, &todos)
	
		assert.Len(t, todos, 2) // Ensure pagination works (2 items per page)
		assert.Equal(t, "Test Todo 1", todos[0].Title)
		assert.Equal(t, "Test Todo 2", todos[1].Title)
	
		// Test fetching with filter by status (completed)
		req = httptest.NewRequest(http.MethodGet, "/api/v1/todos?status=completed", nil)
		resp, _ = app.Test(req)
	
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	
		json.NewDecoder(resp.Body).Decode(&response)
		data, _ = json.Marshal(response["data"])
		todos = nil
		json.Unmarshal(data, &todos)
	
		assert.Len(t, todos, 1) // Only one todo with "done = true"
		assert.Equal(t, "Test Todo 2", todos[0].Title)
	
		// Test searching by title and description
		req = httptest.NewRequest(http.MethodGet, "/api/v1/todos?search=urgent", nil)
		resp, _ = app.Test(req)
	
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	
		json.NewDecoder(resp.Body).Decode(&response)
		data, _ = json.Marshal(response["data"])
		todos = nil
		json.Unmarshal(data, &todos)
	
		assert.Len(t, todos, 1) // Only one todo matches "urgent"
		assert.Equal(t, "Urgent Task", todos[0].Title)
	})

	t.Run("Get Single Todo", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/todos/1", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var todo todoModel.Todo
		json.NewDecoder(resp.Body).Decode(&todo)
		assert.Equal(t, "Test Todo", todo.Title)
	})

	t.Run("Update Todo", func(t *testing.T) {
		todo := map[string]interface{}{"title": "Updated Todo", "description": "Updated Description", "status": "todo", "due_date": "2023-01-01"}
		body, _ := json.Marshal(todo)

		req := httptest.NewRequest(http.MethodPatch, "/api/v1/todos/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result todoModel.Todo
		json.NewDecoder(resp.Body).Decode(&result)
		assert.Equal(t, "Updated Todo", result.Title)
		assert.Equal(t, "Updated Description", result.Description)
	})

	t.Run("Delete Todo", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/todos/1", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		req = httptest.NewRequest(http.MethodGet, "/api/v1/todos/1", nil)
		resp, _ = app.Test(req)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}