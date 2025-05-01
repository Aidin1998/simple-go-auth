package auth

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// NewMiddleware creates an Echo middleware for token validation.
func NewMiddleware(svc *AuthServiceImpl) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			header := c.Request().Header.Get("Authorization")
			if header == "" {
				return c.NoContent(http.StatusUnauthorized)
			}
			// Strip "Bearer " if present
			token := strings.TrimPrefix(header, "Bearer ")
			if err := svc.ValidateToken(c.Request().Context(), token); err != nil {
				return c.NoContent(http.StatusUnauthorized)
			}
			return next(c)
		}
	}
}
