package tests

import (
	"errors"
	"testing"
	"time"
)

// AuthServiceImpl represents a mock implementation of the authentication service.
// AuthServiceImpl is already defined in auth_test.go, so it is removed here.

// Removed duplicate CreateToken method to avoid conflict with auth_test.go.

// Removed duplicate ValidateToken method to avoid conflict with auth_test.go.

// Removed duplicate CheckPasswordHash method to avoid conflict with auth_test.go.

// Removed duplicate Logout method to avoid conflict with auth_test.go.

// Removed duplicate Tokens struct to avoid conflict with auth_test.go.

// Removed duplicate ErrTokenNotFound to avoid conflict with auth_test.go.

func TestAuthService_CreateToken_Success(t *testing.T) {
	service := &AuthServiceImpl{SecretKey: "test-secret", AccessTokenExpiry: 3600}
	userID := "test-user"

	tokens, err := service.CreateToken(userID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if tokens.AccessToken == "" || tokens.RefreshToken == "" {
		t.Fatalf("expected non-empty tokens, got empty tokens")
	}

	if tokens.AtExpires <= time.Now().Unix() || tokens.RtExpires <= time.Now().Unix() {
		t.Fatalf("expected future expiration times, got past times")
	}
}

func TestAuthService_ValidateToken_Success(t *testing.T) {
	service := &AuthServiceImpl{SecretKey: "test-secret", AccessTokenExpiry: 3600}
	userID := "test-user"

	tokens, _ := service.CreateToken(userID)
	validatedUserID, err := service.ValidateToken(tokens.AccessToken)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if validatedUserID != userID {
		t.Fatalf("expected userID %v, got %v", userID, validatedUserID)
	}
}

func TestAuthService_ValidateToken_Failure(t *testing.T) {
	service := &AuthServiceImpl{SecretKey: "test-secret"}
	_, err := service.ValidateToken("invalid-token")
	if err == nil {
		t.Fatalf("expected error for invalid token, got nil")
	}
}

func TestAuthService_HashPassword_Success(t *testing.T) {
	service := &AuthServiceImpl{}
	password := "test-password"

	hash, err := service.HashPassword(password)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if hash == "" {
		t.Fatalf("expected non-empty hash, got empty hash")
	}

	if !service.CheckPasswordHash(password, hash) {
		t.Fatalf("expected password to match hash, but it did not")
	}
}

func TestAuthService_CheckPasswordHash_Failure(t *testing.T) {
	service := &AuthServiceImpl{}
	password := "test-password"
	hash, _ := service.HashPassword(password)

	if service.CheckPasswordHash("wrong-password", hash) {
		t.Fatalf("expected password mismatch, but got match")
	}
}

func TestAuthService_Logout_Success(t *testing.T) {
	service := &AuthServiceImpl{}
	err := service.Logout("test-token")
	if err != nil && !errors.Is(err, ErrTokenNotFound) {
		t.Fatalf("expected no error or token not found error, got %v", err)
	}
}

// Removed duplicate ValidateSignUpInput and isValidEmail methods to avoid conflict with auth_test.go.

func TestValidateSignUpInput_Success(t *testing.T) {
	err := ValidateSignUpInput("valid.email@example.com", "password123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidateSignUpInput_Failure(t *testing.T) {
	tests := []struct {
		email       string
		password    string
		expectError bool
	}{
		{"", "password123", true},
		{"valid.email@example.com", "", true},
		{"invalid-email", "password123", true},
		{"valid.email@example.com", "short", true},
	}

	for _, test := range tests {
		err := ValidateSignUpInput(test.email, test.password)
		if (err != nil) != test.expectError {
			t.Errorf("ValidateSignUpInput(%v, %v) = %v, expected error: %v", test.email, test.password, err, test.expectError)
		}
	}
}
