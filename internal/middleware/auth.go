// internal/middleware/auth.go
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

const (
	UserIDKey    = "user_id"
	UserEmailKey = "user_email"
	IsAdminKey   = "is_admin"
)

// AuthMiddleware проверяет JWT токен в заголовке Authorization
func AuthMiddleware(jwtSecret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			c.Abort()
			return
		}

		parts := strings.Fields(authHeader)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims := &jwt.MapClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		// Извлекаем данные из claims
		userIDFloat, ok := (*claims)["user_id"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token payload"})
			c.Abort()
			return
		}

		email, ok := (*claims)["email"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token payload"})
			c.Abort()
			return
		}

		isAdmin, ok := (*claims)["is_admin"].(bool)
		if !ok {
			isAdmin = false
		}

		// Сохраняем в контекст
		c.Set(UserIDKey, int(userIDFloat))
		c.Set(UserEmailKey, email)
		c.Set(IsAdminKey, isAdmin)

		c.Next()
	}
}

// AdminOnly проверяет, что пользователь является администратором
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		isAdmin, exists := c.Get(IsAdminKey)
		if !exists || !isAdmin.(bool) {
			c.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// GetUserID возвращает ID пользователя из контекста
func GetUserID(c *gin.Context) (int, bool) {
	userID, exists := c.Get(UserIDKey)
	if !exists {
		return 0, false
	}
	return userID.(int), true
}

// GetUserEmail возвращает email пользователя из контекста
func GetUserEmail(c *gin.Context) (string, bool) {
	email, exists := c.Get(UserEmailKey)
	if !exists {
		return "", false
	}
	return email.(string), true
}

// IsAdmin проверяет, является ли текущий пользователь администратором
func IsAdmin(c *gin.Context) bool {
	isAdmin, exists := c.Get(IsAdminKey)
	if !exists {
		return false
	}
	return isAdmin.(bool)
}
