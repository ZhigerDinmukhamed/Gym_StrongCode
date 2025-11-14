package main

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

type key int

const (
	ContextUserID key = iota
	ContextUserEmail
	ContextIsAdmin
)

var jwtSecret = []byte("replace-with-secure-secret") // замените перед деплоем!

// AuthMiddleware verifies JWT token in Authorization header: "Bearer <token>"
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error":"missing auth"})
			return
		}
		parts := strings.Fields(auth)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error":"invalid auth"})
			return
		}
		tokenStr := parts[1]
		claims := &jwt.MapClaims{}
		tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !tkn.Valid {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error":"invalid token"})
			return
		}
		// extract
		uidf := (*claims)["user_id"]
		email := (*claims)["email"].(string)
		isAdmin := (*claims)["is_admin"].(bool)
		// uid might be float64
		var uid int
		switch v := uidf.(type) {
		case float64:
			uid = int(v)
		case int:
			uid = v
		default:
			writeJSON(w,http.StatusUnauthorized,map[string]string{"error":"invalid token payload"})
			return
		}
		ctx := context.WithValue(r.Context(), ContextUserID, uid)
		ctx = context.WithValue(ctx, ContextUserEmail, email)
		ctx = context.WithValue(ctx, ContextIsAdmin, isAdmin)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AdminOnly ensures user is admin
func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isAdmin, ok := r.Context().Value(ContextIsAdmin).(bool)
		if !ok || !isAdmin {
			writeJSON(w, http.StatusForbidden, map[string]string{"error":"admin only"})
			return
		}
		next.ServeHTTP(w, r)
	})
}

