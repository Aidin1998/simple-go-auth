package auth

import (
	"github.com/labstack/echo/v4"
)

// AuthHandler handles authentication-related HTTP requests.
type AuthHandler struct {
	Service *AuthServiceImpl
}

// NewHandler creates a new AuthHandler with the provided service.
func NewHandler(svc *AuthServiceImpl) *AuthHandler {
	return &AuthHandler{Service: svc}
}

// Ping returns pong.
func (h *AuthHandler) Ping(c echo.Context) error {
	return c.String(200, "pong")
}

func (h *AuthHandler) SignUp(c echo.Context) error {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid request body"})
	}

	err := h.Service.SignUp(c.Request().Context(), req.Username, req.Password, req.Email)
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to create user"})
	}

	return c.JSON(201, map[string]string{"message": "user created"})
}

func (h *AuthHandler) SignIn(c echo.Context) error {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid request body"})
	}

	tokens, err := h.Service.SignIn(c.Request().Context(), req.Username, req.Password)
	if err != nil {
		return c.JSON(401, map[string]string{"error": "Invalid credentials"})
	}

	return c.JSON(200, map[string]string{
		"accessToken":  tokens.AccessToken,
		"idToken":      tokens.IdToken,
		"refreshToken": tokens.RefreshToken,
	})
}

func (h *AuthHandler) SignOut(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return c.JSON(400, map[string]string{"error": "Authorization token required"})
	}

	err := h.Service.SignOut(c.Request().Context(), token)
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to sign out"})
	}

	return c.NoContent(204)
}
