package auth

import (
	"testing"
	"time"
)

func TestAuthService_CreateToken(t *testing.T) {
	service := &AuthServiceImpl{SecretKey: "test-secret"}
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

func TestAuthService_ValidateToken(t *testing.T) {
	service := &AuthServiceImpl{SecretKey: "test-secret"}
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

func TestAuthService_HashPassword(t *testing.T) {
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

func TestValidateSignUpInput(t *testing.T) {
	tests := []struct {
		email    string
		password string
		expectErr bool
	}{
		{"", "password123", true},
		{"test@example.com", "", true},
		{"invalid-email", "password123", true},
		{"test@example.com", "short", true},
		{"test@example.com", "validpassword", false},
	}

	for _, test := range tests {
		err := ValidateSignUpInput(test.email, test.password)
		if (err != nil) != test.expectErr {
			t.Errorf("ValidateSignUpInput(%v, %v) = %v, expected error: %v", test.email, test.password, err, test.expectErr)
		}
	}
}

func TestAuthService_CreateToken_Expired(t *testing.T) {
	service := &AuthServiceImpl{SecretKey: "test-secret", AccessTokenExpiry: -1}
	userID := "test-user"
	_, err := service.CreateToken(userID)
	if err == nil {
		t.Fatalf("expected error for expired token, got nil")
	}
}

func TestAuthService_CheckPasswordHash_Invalid(t *testing.T) {
	service := &AuthServiceImpl{}
	password := "test-password"
	hash, _ := service.HashPassword(password)

	if service.CheckPasswordHash("wrong-password", hash) {
		t.Fatalf("expected password mismatch, but got match")
	}
}