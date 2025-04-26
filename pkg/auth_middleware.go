package auth

import (
	"net/http"
	"strings"
)

// AuthMiddleware provides middleware for authentication.
type AuthMiddleware struct {
	Service AuthService
}

// Authenticate validates the JWT token from the Authorization header.
func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		userID, err := m.Service.ValidateToken(tokenParts[1])
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Add user ID to the request context for downstream handlers
		r = r.WithContext(context.WithValue(r.Context(), "userID", userID))
		next.ServeHTTP(w, r)
	})
}