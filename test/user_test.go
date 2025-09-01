package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	userModel "runs-system-user-go/module/user/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserCRUD(t *testing.T) {
	app := SetupTestApp()

	t.Run("Create User", func(t *testing.T) {
		user := map[string]interface{}{"name": "Mikel Arteta", "email": "MikelArteta@pm.me"}
		body, _ := json.Marshal(user)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var result userModel.User
		json.NewDecoder(resp.Body).Decode(&result)
		assert.Equal(t, "Mikel Arteta", result.Name)
		assert.Equal(t, "MikelArteta@pm.me", result.Email)
	})

	t.Run("Get Users with Pagination, Search, and Filter", func(t *testing.T) {
		// Insert mock data into the database
		mockUsers := []userModel.User{
			{Name: "User 1", Email: "MikelArteta@pm.me"},
			{Name: "User 2", Email: "user2@pm.me"},
		}

		for _, user := range mockUsers {
			db.Create(&user) // Ensure the database is set up and accessible
		}

		// Test fetching with pagination
		req := httptest.NewRequest(http.MethodGet, "/api/v1/users?page=1&limit=2", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		data, _ := json.Marshal(response["data"]) // Extract "data" field
		var users []userModel.User
		json.Unmarshal(data, &users)

		assert.Len(t, users, 2) // Ensure pagination works (2 items per page)
		assert.Equal(t, "User 1", users[0].Name)
		assert.Equal(t, "User 2", users[1].Name)
	})

	t.Run("Get Single User", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/users/1", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var user userModel.User
		json.NewDecoder(resp.Body).Decode(&user)
		assert.Equal(t, "Mikel Arteta", user.Name)
	})

	t.Run("Update User", func(t *testing.T) {
		user := map[string]interface{}{"name": "Mikel Arteta", "email": "MikelArteta@pm.me"}
		body, _ := json.Marshal(user)

		req := httptest.NewRequest(http.MethodPatch, "/api/v1/users/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result userModel.User
		json.NewDecoder(resp.Body).Decode(&result)
		assert.Equal(t, "Mikel Arteta", result.Name)
		assert.Equal(t, "MikelArteta@pm.me", result.Email)
	})

	t.Run("Delete User", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/users/1", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		req = httptest.NewRequest(http.MethodGet, "/api/v1/users/1", nil)
		resp, _ = app.Test(req)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}
