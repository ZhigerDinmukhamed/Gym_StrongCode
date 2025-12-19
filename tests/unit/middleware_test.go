package unit

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"Gym_StrongCode/config"
	"Gym_StrongCode/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestAuthMiddleware_ValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	cfg := &config.Config{JWTSecret: "test-secret"}

	r.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	r.GET("/test", func(c *gin.Context) {
		userID := c.GetInt("user_id")
		c.JSON(http.StatusOK, gin.H{"user_id": userID})
	})

	// Создаем валидный токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  1,
		"email":    "test@test.com",
		"is_admin": false,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(cfg.JWTSecret))

	// Тестируем запрос с валидным токеном
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthMiddleware_NoToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	cfg := &config.Config{JWTSecret: "test-secret"}

	r.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	// Запрос без токена
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	cfg := &config.Config{JWTSecret: "test-secret"}

	r.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	// Запрос с невалидным токеном
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	cfg := &config.Config{JWTSecret: "test-secret"}

	r.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	// Создаем истекший токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  1,
		"email":    "test@test.com",
		"is_admin": false,
		"exp":      time.Now().Add(-24 * time.Hour).Unix(), // Истек вчера
	})
	tokenString, _ := token.SignedString([]byte(cfg.JWTSecret))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAdminOnly_AdminUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.Use(func(c *gin.Context) {
		// Симулируем что пользователь уже прошел auth middleware
		c.Set("is_admin", true)
		c.Next()
	})
	r.Use(middleware.AdminOnly())
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAdminOnly_RegularUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.Use(func(c *gin.Context) {
		// Симулируем обычного пользователя
		c.Set("is_admin", false)
		c.Next()
	})
	r.Use(middleware.AdminOnly())
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestRateLimit_AllowsRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.Use(middleware.RateLimitMiddleware()) // 10 запросов в секунду
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	// Делаем 5 запросов - должны пройти
	for i := 0; i < 5; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "127.0.0.1:12345"
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}
}

func TestRateLimit_BlocksExcessRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.Use(middleware.RateLimitMiddleware()) // Глобальный rate limiter
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	// Делаем много запросов быстро, чтобы превысить лимит
	statusOK := 0
	statusTooMany := 0

	for i := 0; i < 20; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "127.0.0.1:12345"
		r.ServeHTTP(w, req)

		switch w.Code {
		case http.StatusOK:
			statusOK++
		case http.StatusTooManyRequests:
			statusTooMany++
		}
	}

	// Проверяем что некоторые запросы были заблокированы
	assert.Greater(t, statusTooMany, 0, "должны быть заблокированные запросы")
}

func TestLoggingMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Создаем zap logger для middleware
	logger, _ := zap.NewDevelopment()
	r.Use(middleware.LoggingMiddleware(logger))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	// Просто проверяем что middleware не ломает запрос
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
