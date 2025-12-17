package integration

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"Gym-StrongCode/tests/testutils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGymHandler_CreateUpdateDelete(t *testing.T) {
	r, db := setupTestRouter(t)
	defer db.Close()

	// Регистрация и логин админа
	adminToken := registerAndLoginAdmin(t, r, db)

	// Создание зала
	gymData := map[string]string{
		"name":    "Test Gym",
		"address": "Test Address",
	}
	jsonData, _ := json.Marshal(gymData)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/admin/gyms", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var createdGym map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createdGym)
	gymID := int(createdGym["id"].(float64))

	// Обновление зала
	updateData := map[string]string{
		"name":    "Updated Gym",
		"address": "New Address",
	}
	jsonData, _ = json.Marshal(updateData)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/admin/gyms/1", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Удаление зала
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/admin/gyms/1", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	_ = gymID // используем переменную
}

func TestClassHandler_CreateUpdateDelete(t *testing.T) {
	r, db := setupTestRouter(t)
	defer db.Close()

	adminToken := registerAndLoginAdmin(t, r, db)

	// Создаем зал и тренера сначала
	gymID := createTestGymViaAPI(t, r, adminToken)
	trainerID := createTestTrainerViaAPI(t, r, adminToken)

	// Создание занятия
	classData := map[string]interface{}{
		"title":        "Yoga Class",
		"description":  "Beginner Yoga",
		"trainer_id":   trainerID,
		"gym_id":       gymID,
		"start_time":   "2025-12-25T10:00:00Z",
		"duration_min": 60,
		"capacity":     20,
	}
	jsonData, _ := json.Marshal(classData)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/admin/classes", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var createdClass map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createdClass)
	classID := int(createdClass["id"].(float64))

	// Удаление занятия
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/admin/classes/1", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	_ = classID // используем переменную
}

func TestBookingHandler_CreateAndCancel(t *testing.T) {
	r, db := setupTestRouter(t)
	defer db.Close()

	// Создаем пользователя и логинимся
	userToken := registerAndLoginUser(t, r, db, "user@test.com")
	adminToken := registerAndLoginAdmin(t, r, db)

	// Создаем зал, тренера и занятие
	gymID := createTestGymViaAPI(t, r, adminToken)
	trainerID := createTestTrainerViaAPI(t, r, adminToken)
	classID := createTestClassViaAPI(t, r, adminToken, gymID, trainerID)

	// Создание бронирования
	bookingData := map[string]interface{}{
		"class_id": classID,
	}
	jsonData, _ := json.Marshal(bookingData)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/bookings", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+userToken)
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	// Получение списка бронирований пользователя
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/bookings", nil)
	req.Header.Set("Authorization", "Bearer "+userToken)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Получение всех бронирований (админ)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/admin/bookings", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Отмена бронирования
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/bookings/1", nil)
	req.Header.Set("Authorization", "Bearer "+userToken)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestTrainerHandler_CreateUpdateDelete(t *testing.T) {
	r, db := setupTestRouter(t)
	defer db.Close()

	adminToken := registerAndLoginAdmin(t, r, db)

	// Создание тренера
	trainerData := map[string]string{
		"name": "John Trainer",
		"bio":  "Expert coach",
	}
	jsonData, _ := json.Marshal(trainerData)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/admin/trainers", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var createdTrainer map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createdTrainer)
	trainerID := int(createdTrainer["id"].(float64))

	// Обновление тренера
	updateData := map[string]string{
		"name": "John Pro Trainer",
		"bio":  "Master coach",
	}
	jsonData, _ = json.Marshal(updateData)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/admin/trainers/1", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Удаление тренера
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/admin/trainers/1", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	_ = trainerID // используем переменную
}

func TestMembershipHandler_CreateUpdateDelete(t *testing.T) {
	r, db := setupTestRouter(t)
	defer db.Close()

	adminToken := registerAndLoginAdmin(t, r, db)

	// Создание подписки
	membershipData := map[string]interface{}{
		"name":          "Premium",
		"duration_days": 30,
		"price_cents":   15000,
	}
	jsonData, _ := json.Marshal(membershipData)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/admin/memberships", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var createdMembership map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createdMembership)
	membershipID := int(createdMembership["id"].(float64))

	// Обновление подписки
	updateData := map[string]interface{}{
		"name":          "Premium Plus",
		"duration_days": 60,
		"price_cents":   25000,
	}
	jsonData, _ = json.Marshal(updateData)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/admin/memberships/1", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Удаление подписки
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/admin/memberships/1", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	_ = membershipID // используем переменную
}

// Вспомогательные функции
func registerAndLoginAdmin(t *testing.T, r *gin.Engine, db *sql.DB) string {
	// Создаем админа напрямую в базе (используя функцию из testutils)
	testutils.CreateTestUser(t, db, "admin@test.com", "password123", true)

	// Логин админа
	loginData := map[string]string{
		"email":    "admin@test.com",
		"password": "password123",
	}
	jsonData, _ := json.Marshal(loginData)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/users/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	var loginResp map[string]string
	json.Unmarshal(w.Body.Bytes(), &loginResp)
	return loginResp["token"]
}

func registerAndLoginUser(t *testing.T, r *gin.Engine, db *sql.DB, email string) string {
	// Регистрация
	registerData := map[string]string{
		"name":     "Test User",
		"email":    email,
		"password": "password123",
	}
	jsonData, _ := json.Marshal(registerData)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/users/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	// Логин
	loginData := map[string]string{
		"email":    email,
		"password": "password123",
	}
	jsonData, _ = json.Marshal(loginData)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/users/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	var loginResp map[string]string
	json.Unmarshal(w.Body.Bytes(), &loginResp)
	return loginResp["token"]
}

func createTestGymViaAPI(t *testing.T, r *gin.Engine, token string) int {
	require.NotEmpty(t, token, "token не должен быть пустым")
	gymData := map[string]string{
		"name":    "Test Gym",
		"address": "Test Address",
	}
	jsonData, _ := json.Marshal(gymData)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/admin/gyms", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	var gym map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &gym)
	return int(gym["id"].(float64))
}

func createTestTrainerViaAPI(t *testing.T, r *gin.Engine, token string) int {
	require.NotEmpty(t, token, "token не должен быть пустым")
	trainerData := map[string]string{
		"name": "Test Trainer",
		"bio":  "Test Bio",
	}
	jsonData, _ := json.Marshal(trainerData)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/admin/trainers", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	var trainer map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &trainer)
	return int(trainer["id"].(float64))
}

func createTestClassViaAPI(t *testing.T, r *gin.Engine, token string, gymID, trainerID int) int {
	require.NotEmpty(t, token, "token не должен быть пустым")
	classData := map[string]interface{}{
		"title":        "Test Class",
		"description":  "Test Description",
		"trainer_id":   trainerID,
		"gym_id":       gymID,
		"start_time":   "2025-12-25T10:00:00Z",
		"duration_min": 60,
		"capacity":     20,
	}
	jsonData, _ := json.Marshal(classData)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/admin/classes", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	var class map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &class)
	return int(class["id"].(float64))
}
