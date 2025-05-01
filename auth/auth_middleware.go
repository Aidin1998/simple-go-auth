package auth

import (
	"strings"

	"github.com/golang-jwt/jwt/v4"
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
			token, err := m.Service.ValidateToken(tokenString)
			if err != nil || !token.Valid {
				return c.JSON(401, map[string]string{"error": "Invalid or expired token"})
			}

			// Add user information to the context
			claims, ok := token.Claims.(jwt.MapClaims)
			if ok && token.Valid {
				c.Set("userID", claims["userID"])
			}

			return next(c)
		}
	}
}
