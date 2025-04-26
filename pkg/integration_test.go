package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIntegration_SignUpAndSignIn(t *testing.T) {
	// Initialize AuthService and AuthHandler
	service := &AuthServiceImpl{SecretKey: "test-secret"}
	handler := &AuthHandler{Service: service}

	// Test SignUp
	signUpPayload := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	signUpBody, _ := json.Marshal(signUpPayload)
	signUpReq := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(signUpBody))
	signUpReq.Header.Set("Content-Type", "application/json")
	signUpRec := httptest.NewRecorder()
	handler.SignUp(signUpRec, signUpReq)

	if signUpRec.Code != http.StatusCreated {
		t.Fatalf("expected status %v, got %v", http.StatusCreated, signUpRec.Code)
	}

	// Test SignIn
	signInPayload := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	signInBody, _ := json.Marshal(signInPayload)
	signInReq := httptest.NewRequest(http.MethodPost, "/signin", bytes.NewBuffer(signInBody))
	signInReq.Header.Set("Content-Type", "application/json")
	signInRec := httptest.NewRecorder()
	handler.SignIn(signInRec, signInReq)

	if signInRec.Code != http.StatusOK {
		t.Fatalf("expected status %v, got %v", http.StatusOK, signInRec.Code)
	}

	var tokens TokenDetails
	if err := json.NewDecoder(signInRec.Body).Decode(&tokens); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if tokens.AccessToken == "" || tokens.RefreshToken == "" {
		t.Fatalf("expected non-empty tokens, got empty tokens")
	}
}