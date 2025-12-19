package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRegistration_ValidationErrors(t *testing.T) {
	r, _ := setupTestRouter(t)

	tests := []struct {
		name       string
		body       map[string]interface{}
		wantStatus int
	}{
		{
			name:       "missing name",
			body:       map[string]interface{}{"email": "test@example.com", "password": "pass123"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing email",
			body:       map[string]interface{}{"name": "Test", "password": "pass123"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing password",
			body:       map[string]interface{}{"name": "Test", "email": "test@example.com"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "short password",
			body:       map[string]interface{}{"name": "Test", "email": "test@example.com", "password": "123"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid email format",
			body:       map[string]interface{}{"name": "Test", "email": "invalid-email", "password": "pass123"},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.body)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/users/register", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestUserFlow_CompleteScenario(t *testing.T) {
	r, _ := setupTestRouter(t)

	// 1. Регистрация
	registerBody := map[string]string{
		"name":     "Alice Johnson",
		"email":    "alice@example.com",
		"password": "securepass123",
	}
	jsonBody, _ := json.Marshal(registerBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/users/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	// 2. Логин
	loginBody := map[string]string{
		"email":    "alice@example.com",
		"password": "securepass123",
	}
	jsonBody, _ = json.Marshal(loginBody)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/users/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var loginResp map[string]string
	json.Unmarshal(w.Body.Bytes(), &loginResp)
	token := loginResp["token"]

	// 3. Получение профиля
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var userResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &userResp)
	assert.Equal(t, "Alice Johnson", userResp["name"])
	assert.Equal(t, "alice@example.com", userResp["email"])
	assert.False(t, userResp["is_admin"].(bool))
}
