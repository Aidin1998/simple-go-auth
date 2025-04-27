package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

type AuthMiddleware struct {
	Service *AuthServiceImpl
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := m.Service.ValidateToken(tokenString)
		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Add user information to the request context
		claims, ok := token.Claims.(jwt.MapClaims)
		if ok && token.Valid {
			ctx := context.WithValue(r.Context(), "userID", claims["userID"])
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}
