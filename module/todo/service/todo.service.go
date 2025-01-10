package todoService

import (
	"encoding/json"
	"errors"
	"fmt"
	"todolist-ilcs-api/database"
	todoModel "todolist-ilcs-api/module/todo/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GetPaginatedTodos fetches todos with pagination, filters, and search.
func GetPaginatedTodos(page, limit int, status, search string) ([]todoModel.Todo, int, error) {
	db := database.DB
	redisDB := database.RedisDB
	redisTTL := database.RedisCacheTTL

	offset := (page - 1) * limit
	cacheKey := fmt.Sprintf("todos:page:%d:limit:%d:status:%s:search:%s", page, limit, status, search)

	// Try fetching from cache
	cachedTodos, err := redisDB.Get(redisDB.Context(), cacheKey).Result()
	if err == nil {
		var todos []todoModel.Todo
		if json.Unmarshal([]byte(cachedTodos), &todos) == nil {
			return todos, len(todos), nil
		}
	}

	// Fetch from database
	var todos []todoModel.Todo
	query := db.Offset(offset).Limit(limit)

	if status != "" {
		isCompleted := status == "completed"
		query = query.Where("status = ?", isCompleted)
	}

	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("title ILIKE ? OR description ILIKE ?", searchTerm, searchTerm)
	}

	result := query.Find(&todos)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	// Cache results
	data, err := json.Marshal(todos)
	if err == nil {
		redisDB.Set(redisDB.Context(), cacheKey, data, redisTTL)
	}

	return todos, len(todos), nil
}

// GetTodoByID retrieves a todo by its ID.
func GetTodoByID(id string) (todoModel.Todo, error) {
	db := database.DB
	var todo todoModel.Todo

	if err := db.First(&todo, "id = ?", id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return todo, errors.New("todo not found")
	}

	return todo, nil
}

// CreateTodo creates a new todo.
func CreateTodo(todo *todoModel.Todo) error {
	db := database.DB
	todo.ID = uuid.New()

	if err := db.Create(todo).Error; err != nil {
		return err
	}

	return nil
}

// UpdateTodo updates an existing todo.
func UpdateTodo(id string, data map[string]interface{}) (todoModel.Todo, error) {
	db := database.DB
	var todo todoModel.Todo

	if err := db.First(&todo, "id = ?", id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return todo, errors.New("todo not found")
	}

	if err := db.Model(&todo).Updates(data).Error; err != nil {
		return todo, err
	}

	return todo, nil
}

// DeleteTodo deletes a todo by ID.
func DeleteTodo(id string) error {
	db := database.DB
	var todo todoModel.Todo

	if err := db.First(&todo, "id = ?", id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("todo not found")
	}

	if err := db.Delete(&todo).Error; err != nil {
		return err
	}

	return nil
}
