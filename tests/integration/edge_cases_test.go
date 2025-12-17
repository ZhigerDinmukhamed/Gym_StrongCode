package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"Gym-StrongCode/tests/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConcurrentUserRegistration(t *testing.T) {
	r, _ := setupTestRouter(t)

	// Тест на race condition при одновременной регистрации
	done := make(chan bool, 2)

	registerUser := func(email string) {
		defer func() { done <- true }()

		body := map[string]string{
			"name":     "Test User",
			"email":    email,
			"password": "password123",
		}
		jsonBody, _ := json.Marshal(body)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/users/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
	}

	go registerUser("concurrent1@example.com")
	go registerUser("concurrent2@example.com")

	<-done
	<-done
	// Если нет паники - тест пройден
}

func TestGymManagement_AdminOnly(t *testing.T) {
	r, db := setupTestRouter(t)

	// Создаем обычного пользователя
	testutils.CreateTestUser(t, db, "regular@example.com", "password123", false)

	// Логинимся
	loginBody := map[string]string{
		"email":    "regular@example.com",
		"password": "password123",
	}
	jsonBody, _ := json.Marshal(loginBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/users/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	var loginResp map[string]string
	json.Unmarshal(w.Body.Bytes(), &loginResp)
	regularToken := loginResp["token"]

	// Пытаемся создать зал (только для админов)
	gymBody := map[string]string{
		"name":    "New Gym",
		"address": "Test Address",
	}
	jsonBody, _ = json.Marshal(gymBody)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/admin/gyms", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+regularToken)
	r.ServeHTTP(w, req)

	// Должен быть запрещен доступ
	assert.Equal(t, http.StatusForbidden, w.Code)

	// Теперь с админом
	testutils.CreateTestUser(t, db, "admin@example.com", "password123", true)
	loginBody = map[string]string{
		"email":    "admin@example.com",
		"password": "password123",
	}
	jsonBody, _ = json.Marshal(loginBody)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/users/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	json.Unmarshal(w.Body.Bytes(), &loginResp)
	adminToken := loginResp["token"]

	// Создаем зал с админским токеном
	gymBody2 := map[string]string{
		"name":    "Admin Gym",
		"address": "Admin Address",
	}
	jsonBody2, _ := json.Marshal(gymBody2)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/admin/gyms", bytes.NewBuffer(jsonBody2))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestPublicEndpoints_NoAuthRequired(t *testing.T) {
	r, db := setupTestRouter(t)

	// Создаем тестовые данные
	gymID := testutils.CreateTestGym(t, db, "Public Gym", "Address 1")
	trainerID := testutils.CreateTestTrainer(t, db, "Trainer", "Bio")
	testutils.CreateTestClass(t, db, "Public Class", trainerID, gymID, 20)
	testutils.CreateTestMembership(t, db, "Basic", 30, 1000000)

	// Тестируем публичные эндпоинты без токена
	tests := []struct {
		name     string
		endpoint string
	}{
		{"List Gyms", "/api/gyms"},
		{"List Classes", "/api/classes"},
		{"List Memberships", "/api/memberships"},
		{"List Trainers", "/api/trainers"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", tt.endpoint, nil)
			r.ServeHTTP(w, req)

			require.Equal(t, http.StatusOK, w.Code)

			var result []interface{}
			err := json.Unmarshal(w.Body.Bytes(), &result)
			require.NoError(t, err)
			assert.NotEmpty(t, result)
		})
	}
}

func TestInvalidToken(t *testing.T) {
	r, _ := setupTestRouter(t)

	tests := []struct {
		name  string
		token string
	}{
		{"No Bearer prefix", "invalidtoken"},
		{"Invalid JWT", "Bearer invalid.jwt.token"},
		{"Expired token", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDk0NTkyMDB9.fake"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/me", nil)
			req.Header.Set("Authorization", tt.token)
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusUnauthorized, w.Code)
		})
	}
}

func TestSQLInjectionProtection(t *testing.T) {
	r, _ := setupTestRouter(t)

	// Попытка SQL-инъекции через email
	maliciousBody := map[string]string{
		"email":    "admin@example.com' OR '1'='1",
		"password": "anything",
	}

	jsonBody, _ := json.Marshal(maliciousBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/users/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	// Должна быть ошибка, не успешный логин
	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestLargePayload(t *testing.T) {
	r, _ := setupTestRouter(t)

	// Большой payload
	largeString := string(make([]byte, 1024*1024)) // 1MB
	body := map[string]string{
		"name":     largeString,
		"email":    "test@example.com",
		"password": "password123",
	}

	jsonBody, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/users/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	// Должна быть обработана корректно (или отклонена по размеру)
	assert.NotEqual(t, http.StatusInternalServerError, w.Code)
}
