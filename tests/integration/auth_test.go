package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"Gym-StrongCode/cmd"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func getRouter() *gin.Engine {
	// Запускаем main без сервера, только роутер
	// Здесь можно использовать рефакторинг main в отдельную функцию initRouter()
	// Для простоты — создаём новый роутер с теми же настройками
	// (или вынесите роутер в отдельную функцию)
	return nil // Замените на реальный роутер или используйте тестовой БД
}

func TestRegisterAndLogin(t *testing.T) {
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
	assert.Contains(t, resp, "token")
}