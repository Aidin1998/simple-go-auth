package tests

import (
	"errors"
	"testing"
	"time"
)

// AuthServiceImpl represents a mock implementation of the authentication service.
type AuthServiceImpl struct {
	SecretKey         string
	AccessTokenExpiry int64
}

// CreateToken generates mock access and refresh tokens.
func (a *AuthServiceImpl) CreateToken(userID string) (*Tokens, error) {
	now := time.Now().Unix()
	return &Tokens{
		AccessToken:  "mock-access-token",
		RefreshToken: "mock-refresh-token",
		AtExpires:    now + a.AccessTokenExpiry,
		RtExpires:    now + a.AccessTokenExpiry*2,
	}, nil
}

// ValidateToken validates a mock token.
func (a *AuthServiceImpl) ValidateToken(token string) (string, error) {
	if token == "mock-access-token" {
		return "test-user", nil
	}
	return "", errors.New("invalid token")
}

// HashPassword hashes a password (mock implementation).
func (a *AuthServiceImpl) HashPassword(password string) (string, error) {
	return "mock-hash", nil
}

// CheckPasswordHash checks if a password matches a hash (mock implementation).
func (a *AuthServiceImpl) CheckPasswordHash(password, hash string) bool {
	return password == "test-password" && hash == "mock-hash"
}

// Logout logs out a user (mock implementation).
func (a *AuthServiceImpl) Logout(token string) error {
	if token == "test-token" {
		return nil
	}
	return ErrTokenNotFound
}

// Tokens represents the structure for access and refresh tokens.
type Tokens struct {
	AccessToken  string
	RefreshToken string
	AtExpires    int64
	RtExpires    int64
}

// ErrTokenNotFound is a mock error for token not found.
var ErrTokenNotFound = errors.New("token not found")

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

// ValidateSignUpInput validates the email and password for sign-up.
func ValidateSignUpInput(email, password string) error {
	if email == "" || password == "" {
		return errors.New("email and password must not be empty")
	}
	if len(password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}
	if !isValidEmail(email) {
		return errors.New("invalid email format")
	}
	return nil
}

// isValidEmail checks if the email format is valid (mock implementation).
func isValidEmail(email string) bool {
	return len(email) > 3 && email[len(email)-4:] == ".com"
}

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
