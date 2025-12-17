package integration

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"Gym-StrongCode/config"
	"Gym-StrongCode/internal/handler"
	"Gym-StrongCode/internal/middleware"
	"Gym-StrongCode/internal/repository"
	"Gym-StrongCode/internal/service"
	"Gym-StrongCode/tests/testutils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRouter(t *testing.T) (*gin.Engine, *sql.DB) {
	gin.SetMode(gin.TestMode)
	db := testutils.SetupTestDB(t)

	// Репозитории
	userRepo := repository.NewUserRepository(db)
	gymRepo := repository.NewGymRepository(db)
	membershipRepo := repository.NewMembershipRepository(db)
	trainerRepo := repository.NewTrainerRepository(db)
	classRepo := repository.NewClassRepository(db)
	bookingRepo := repository.NewBookingRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)

	// Сервисы
	cfg := &config.Config{
		JWTSecret: "test-secret-key",
		SMTPHost:  "",
	}
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	notificationService := service.NewNotificationService(cfg)
	gymService := service.NewGymService(gymRepo)
	membershipService := service.NewMembershipService(membershipRepo, paymentRepo, db, notificationService)
	trainerService := service.NewTrainerService(trainerRepo)
	classService := service.NewClassService(classRepo, trainerRepo, gymRepo)
	bookingService := service.NewBookingService(bookingRepo, classRepo, membershipRepo, notificationService)
	paymentService := service.NewPaymentService(paymentRepo)

	// Хендлеры
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userRepo)
	gymHandler := handler.NewGymHandler(gymService)
	membershipHandler := handler.NewMembershipHandler(membershipService)
	trainerHandler := handler.NewTrainerHandler(trainerService)
	classHandler := handler.NewClassHandler(classService)
	bookingHandler := handler.NewBookingHandler(bookingService, userRepo)
	paymentHandler := handler.NewPaymentHandler(paymentService)

	// Роутер
	r := gin.Default()

	api := r.Group("/api")
	{
		api.GET("/health", handler.HealthCheck)

		// Публичные
		api.POST("/users/register", authHandler.Register)
		api.POST("/users/login", authHandler.Login)
		api.GET("/classes", classHandler.List)
		api.GET("/gyms", gymHandler.List)
		api.GET("/memberships", membershipHandler.List)
		api.GET("/trainers", trainerHandler.List)

		// Авторизованные
		authorized := api.Group("")
		authorized.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			authorized.GET("/me", userHandler.GetCurrent)

			authorized.POST("/bookings", bookingHandler.Create)
			authorized.GET("/bookings", bookingHandler.ListUser)
			authorized.DELETE("/bookings/:id", bookingHandler.Cancel)

			authorized.POST("/memberships/buy", membershipHandler.Buy)

			authorized.POST("/payments", paymentHandler.CreateStandalone)
		}

		// Админские
		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		admin.Use(middleware.AdminOnly())
		{
			admin.GET("/users", userHandler.List)
			admin.DELETE("/users/:id", userHandler.Delete)

			admin.POST("/gyms", gymHandler.Create)
			admin.PUT("/gyms/:id", gymHandler.Update)
			admin.DELETE("/gyms/:id", gymHandler.Delete)

			admin.POST("/trainers", trainerHandler.Create)
			admin.PUT("/trainers/:id", trainerHandler.Update)
			admin.DELETE("/trainers/:id", trainerHandler.Delete)

			admin.POST("/classes", classHandler.Create)
			admin.PUT("/classes/:id", classHandler.Update)
			admin.DELETE("/classes/:id", classHandler.Delete)

			admin.GET("/bookings", bookingHandler.ListAll)

			admin.POST("/memberships", membershipHandler.Create)
			admin.PUT("/memberships/:id", membershipHandler.Update)
			admin.DELETE("/memberships/:id", membershipHandler.Delete)

			admin.GET("/payments", paymentHandler.ListAll)
		}
	}

	return r, db
}

func TestHealthCheck(t *testing.T) {
	r, _ := setupTestRouter(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/health", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "ok")
}

func TestRegisterAndLogin(t *testing.T) {
	r, _ := setupTestRouter(t)

	// 1. Регистрация
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

	require.Equal(t, http.StatusCreated, w.Code)

	var registerResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &registerResp)
	assert.NotZero(t, registerResp["id"])
	assert.Equal(t, "Test User", registerResp["name"])
	assert.Equal(t, "test@example.com", registerResp["email"])

	// 2. Логин
	loginBody := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	jsonBody, _ = json.Marshal(loginBody)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/users/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var loginResp map[string]string
	json.Unmarshal(w.Body.Bytes(), &loginResp)
	assert.NotEmpty(t, loginResp["token"])
}

func TestRegister_DuplicateEmail(t *testing.T) {
	r, _ := setupTestRouter(t)

	registerBody := map[string]string{
		"name":     "Test User",
		"email":    "duplicate@example.com",
		"password": "password123",
	}

	// Первая регистрация
	jsonBody, _ := json.Marshal(registerBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/users/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	// Повторная регистрация с тем же email
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/users/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLogin_InvalidCredentials(t *testing.T) {
	r, db := setupTestRouter(t)

	// Создаем пользователя напрямую
	testutils.CreateTestUser(t, db, "valid@example.com", "correctpassword", false)

	loginBody := map[string]string{
		"email":    "valid@example.com",
		"password": "wrongpassword",
	}

	jsonBody, _ := json.Marshal(loginBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/users/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetCurrentUser(t *testing.T) {
	r, db := setupTestRouter(t)

	// Создаем пользователя и получаем токен
	testutils.CreateTestUser(t, db, "user@example.com", "password123", false)

	// Логинимся
	loginBody := map[string]string{
		"email":    "user@example.com",
		"password": "password123",
	}
	jsonBody, _ := json.Marshal(loginBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/users/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	var loginResp map[string]string
	json.Unmarshal(w.Body.Bytes(), &loginResp)
	token := loginResp["token"]

	// Запрашиваем текущего пользователя
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var userResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &userResp)
	assert.Equal(t, "user@example.com", userResp["email"])
	assert.Equal(t, "Test User", userResp["name"])
}

func TestUnauthorizedAccess(t *testing.T) {
	r, _ := setupTestRouter(t)

	// Пытаемся получить доступ к защищенному эндпоинту без токена
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/me", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
