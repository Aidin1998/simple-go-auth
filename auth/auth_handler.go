package auth

import (
	"encoding/json"
	"net/http"

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

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	// Placeholder for signup logic
	w.Write([]byte("SignUp endpoint not implemented"))
}

func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		UserID string `json:"userID"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token, err := h.Service.GenerateToken(requestBody.UserID)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *AuthHandler) SignOut(w http.ResponseWriter, r *http.Request) {
	// Placeholder for signout logic
	w.Write([]byte("SignOut endpoint not implemented"))
}
