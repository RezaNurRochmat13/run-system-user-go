package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuth(t *testing.T) {
	app := SetupTestApp()

	t.Run("Register User", func(t *testing.T) {
		user := map[string]interface{}{"name": "Mikel Arteta", "email": "MikelArteta@pm.me", "password": "password"}
		body, _ := json.Marshal(user)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	t.Run("Login User", func(t *testing.T) {
		user := map[string]interface{}{"email": "MikelArteta@pm.me", "password": "password"}
		body, _ := json.Marshal(user)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

}
