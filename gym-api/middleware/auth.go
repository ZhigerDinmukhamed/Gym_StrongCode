package middleware

import (
	"net/http"
	"strings"

	"gym-api/utils"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			http.Error(w, "Missing token", 401)
			return
		}

		token := strings.TrimPrefix(auth, "Bearer ")

		claims, err := utils.ParseToken(token)
		if err != nil {
			http.Error(w, "Invalid token", 401)
			return
		}

		// Add user info into context
		ctx := utils.ContextWithUser(r.Context(), claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
