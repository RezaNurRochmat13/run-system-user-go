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

		req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var result todoModel.Todo
		json.NewDecoder(resp.Body).Decode(&result)
		assert.Equal(t, "Test Todo", result.Title)
	})

	t.Run("Get Todos", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/todos", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var todos []todoModel.Todo
		json.NewDecoder(resp.Body).Decode(&todos)
		assert.Len(t, todos, 1)
		assert.Equal(t, "Test Todo", todos[0].Title)
	})

	t.Run("Get Single Todo", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/todos/1", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var todo todoModel.Todo
		json.NewDecoder(resp.Body).Decode(&todo)
		assert.Equal(t, "Test Todo", todo.Title)
	})

	t.Run("Update Todo", func(t *testing.T) {
		todo := map[string]interface{}{"title": "Updated Todo", "description": "Updated Description", "status": "todo", "due_date": "2023-01-01"}
		body, _ := json.Marshal(todo)

		req := httptest.NewRequest(http.MethodPatch, "/todos/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result todoModel.Todo
		json.NewDecoder(resp.Body).Decode(&result)
		assert.Equal(t, "Updated Todo", result.Title)
		assert.Equal(t, "Updated Description", result.Description)
	})

	t.Run("Delete Todo", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/todos/1", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		req = httptest.NewRequest(http.MethodGet, "/todos/1", nil)
		resp, _ = app.Test(req)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}