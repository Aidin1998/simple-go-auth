package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFullAPIFlow(t *testing.T) {
	// Mock AWS Secrets Manager setup
	jwtSecretKey = "testsecret"

	// Step 1: Simulate user sign-up
	email := "test@example.com"
	password := "password123"
	err := ValidateSignUpInput(email, password)
	assert.NoError(t, err)

	hashedPassword, err := HashPassword(password)
	assert.NoError(t, err)

	// Step 2: Simulate user login and JWT generation
	userID := "12345"
	token, err := GenerateJWT(userID, time.Hour)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Step 3: Simulate API request with JWT validation
	req := httptest.NewRequest(http.MethodGet, "/protected-endpoint", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := ValidateJWT(r.Header.Get("Authorization")[7:])
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		assert.Equal(t, userID, claims["user_id"])
		w.WriteHeader(http.StatusOK)
	})

	handler.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Code)
}