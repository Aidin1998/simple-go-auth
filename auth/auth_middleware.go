package auth

import (
	"strings"

	"github.com/labstack/echo/v4"
)

type AuthMiddleware struct {
	Service *AuthServiceImpl
}

func (m *AuthMiddleware) Authenticate() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(401, map[string]string{"error": "Missing Authorization header"})
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			_, err := m.Service.ValidateToken(tokenString)
			if err != nil {
				return c.JSON(401, map[string]string{"error": "Invalid or expired token"})
			}

			// Additional validation for token claims can be added here if needed

			return next(c)
		}
	}
}
