package auth

import (
	"strings"

	"github.com/labstack/echo/v4"
)

// NewMiddleware creates a middleware function for token validation.
func NewMiddleware(svc *AuthServiceImpl) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get("Authorization")
			if token == "" {
				return c.NoContent(401)
			}

			// Optionally strip "Bearer " prefix
			token = strings.TrimPrefix(token, "Bearer ")

			if err := svc.ValidateToken(c.Request().Context(), token); err != nil {
				return c.NoContent(401)
			}

			return next(c)
		}
	}
}
