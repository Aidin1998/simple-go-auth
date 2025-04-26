package auth

import (
	"encoding/json"
	"net/http"
)

// AuthHandler provides HTTP handlers for authentication.
type AuthHandler struct {
	Service AuthService
}

// SignUpRequest represents the payload for the sign-up endpoint.
type SignUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// SignUp handles user registration.
func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var req SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	hashedPassword, err := h.Service.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	// Simulate user creation (replace with DB logic)
	user := User{
		ID:       "1", // Replace with generated ID
		Email:    req.Email,
		Password: hashedPassword,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// SignInRequest represents the payload for the sign-in endpoint.
type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// SignIn handles user login.
func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var req SignInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Simulate user retrieval (replace with DB logic)
	storedUser := User{
		ID:       "1",
		Email:    req.Email,
		Password: "$2a$10$7Q9h8J9J9J9J9J9J9J9J9u", // Replace with hashed password from DB
	}

	if !h.Service.CheckPasswordHash(req.Password, storedUser.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	tokens, err := h.Service.CreateToken(storedUser.ID)
	if err != nil {
		http.Error(w, "Error generating tokens", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tokens)
}

// SignOut handles user logout.
func (h *AuthHandler) SignOut(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User logged out"))
}

// RefreshToken handles token refreshing.
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Token refresh not implemented yet"))
}