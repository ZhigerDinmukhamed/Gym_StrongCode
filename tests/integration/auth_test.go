package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func getRouter() *gin.Engine {
	// TODO: Implement router setup for tests
	// For now, skip tests that require full setup
	return nil
}

func TestAuthFlow_Old(t *testing.T) {
	t.Skip("Legacy test - replaced by api_test.go TestRegisterAndLogin")
	r := getRouter()

	// Register
	registerBody := map[string]string{
		"name":     "Test User",
		"email":    "test@example.com",
		"password": "password123",
	}
	jsonBody, _ := json.Marshal(registerBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/users/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	// Login
	loginBody := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	jsonBody, _ = json.Marshal(loginBody)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/users/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NotEmpty(t, resp["token"])
}
