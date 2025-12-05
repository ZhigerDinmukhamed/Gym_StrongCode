package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func generateTestToken(secret []byte, claims jwt.MapClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(secret)
	return tokenString
}

func TestAuthMiddlewareValidToken(t *testing.T) {
	secret := []byte("test-secret-key")
	claims := jwt.MapClaims{
		"user_id":  float64(1),
		"email":    "test@example.com",
		"is_admin": true,
	}
	tokenString := generateTestToken(secret, claims)

	router := gin.New()
	router.Use(AuthMiddleware(secret))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestAuthMiddlewareMissingHeader(t *testing.T) {
	secret := []byte("test-secret-key")
	router := gin.New()
	router.Use(AuthMiddleware(secret))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestAuthMiddlewareInvalidFormat(t *testing.T) {
	secret := []byte("test-secret-key")
	router := gin.New()
	router.Use(AuthMiddleware(secret))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "InvalidToken")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestAdminOnlyWithAdmin(t *testing.T) {
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(IsAdminKey, true)
		c.Next()
	})
	router.Use(AdminOnly())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestAdminOnlyWithoutAdmin(t *testing.T) {
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(IsAdminKey, false)
		c.Next()
	})
	router.Use(AdminOnly())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

func TestGetUserID(t *testing.T) {
	c := &gin.Context{}
	c.Set(UserIDKey, 42)

	userID, exists := GetUserID(c)
	if !exists || userID != 42 {
		t.Errorf("expected userID 42, got %d", userID)
	}
}

func TestGetUserEmail(t *testing.T) {
	c := &gin.Context{}
	c.Set(UserEmailKey, "user@example.com")

	email, exists := GetUserEmail(c)
	if !exists || email != "user@example.com" {
		t.Errorf("expected email user@example.com, got %s", email)
	}
}

func TestIsAdmin(t *testing.T) {
	c := &gin.Context{}
	c.Set(IsAdminKey, true)

	if !IsAdmin(c) {
		t.Error("expected IsAdmin to return true")
	}
}
