package auth

import (
	"net/http"
	"encoding/json"
)

// SignUpHandler handles user sign-up requests.
func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := ValidateSignUpInput(requestBody.Email, requestBody.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashedPassword, err := HashPassword(requestBody.Password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	// Simulate user creation (e.g., save to database)
	response := map[string]string{
		"message": "User created successfully",
		"email":   requestBody.Email,
		"password": hashedPassword, // In production, never return the password.
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}