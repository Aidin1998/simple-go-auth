package auth

import (
	"encoding/json"

	"github.com/labstack/echo/v4"
)

// AuthHandler handles authentication-related HTTP requests.
type AuthHandler struct {
	Service *AuthServiceImpl
}

func NewHandler() *AuthHandler {
	return &AuthHandler{}
}

// Ping returns pong.
func (h *AuthHandler) Ping(c echo.Context) error {
	return c.String(200, "pong")
}

func (h *AuthHandler) SignUp(c echo.Context) error {
	// Placeholder for signup logic
	return c.String(200, "SignUp endpoint not implemented")
}

func (h *AuthHandler) SignIn(c echo.Context) error {
	var requestBody struct {
		UserID string `json:"userID"`
	}
	if err := json.NewDecoder(c.Request().Body).Decode(&requestBody); err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid request body"})
	}

	token, err := h.Service.GenerateToken(requestBody.UserID)
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to generate token"})
	}

	return c.JSON(200, map[string]string{"token": token})
}

func (h *AuthHandler) SignOut(c echo.Context) error {
	// Placeholder for signout logic
	return c.String(200, "SignOut endpoint not implemented")
}
