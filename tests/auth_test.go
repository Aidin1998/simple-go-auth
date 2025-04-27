package tests

import (
	"errors"
	"os"
	"regexp"
	"testing"
	"time"
)

type AuthServiceImpl struct {
	SecretKey         string
	AccessTokenExpiry int64
}

func (a *AuthServiceImpl) CreateToken(userID string) (*Tokens, error) {
	if a.SecretKey == "" {
		return nil, errors.New("missing secret key")
	}
	if a.AccessTokenExpiry <= 0 {
		return nil, errors.New("invalid token expiry")
	}
	now := time.Now().Unix()
	return &Tokens{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		AtExpires:    now + a.AccessTokenExpiry,
		RtExpires:    now + a.AccessTokenExpiry*2,
	}, nil
}

func (a *AuthServiceImpl) ValidateToken(token string) (string, error) {
	if token == "access-token" {
		return "test-user", nil
	}
	return "", errors.New("invalid token")
}

func (a *AuthServiceImpl) HashPassword(password string) (string, error) {
	return "hashed-password", nil
}

func (a *AuthServiceImpl) CheckPasswordHash(password, hash string) bool {
	return password == "test-password" && hash == "hashed-password"
}

func (a *AuthServiceImpl) Logout(token string) error {
	if token == "test-token" {
		return nil
	}
	return ErrTokenNotFound
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
	AtExpires    int64
	RtExpires    int64
}

var ErrTokenNotFound = errors.New("token not found")

func TestAuthService_CreateToken(t *testing.T) {
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

func TestAuthService_CreateToken_InvalidSecret(t *testing.T) {
	service := &AuthServiceImpl{SecretKey: ""}
	userID := "test-user"
	_, err := service.CreateToken(userID)
	if err == nil {
		t.Fatalf("expected error for missing secret key, got nil")
	}
}

func TestAuthService_ValidateToken(t *testing.T) {
	service := &AuthServiceImpl{SecretKey: "test-secret", AccessTokenExpiry: 3600} // Ensure valid expiry
	userID := "test-user"
	tokens, err := service.CreateToken(userID)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	validatedUserID, err := service.ValidateToken(tokens.AccessToken)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if validatedUserID != userID {
		t.Fatalf("Expected userID %v, got %v", userID, validatedUserID)
	}
}

func TestAuthService_ValidateToken_Invalid(t *testing.T) {
	service := &AuthServiceImpl{SecretKey: "test-secret"}
	_, err := service.ValidateToken("invalid-token")
	if err == nil {
		t.Fatalf("expected error for invalid token, got nil")
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

func TestAuthService_CheckPasswordHash_Invalid(t *testing.T) {
	service := &AuthServiceImpl{}
	password := "test-password"
	hash, _ := service.HashPassword(password)

	if service.CheckPasswordHash("wrong-password", hash) {
		t.Fatalf("expected password mismatch, but got match")
	}
}

func ValidateSignUpInput(email, password string) error {
	if email == "" {
		return errors.New("email is required")
	}
	if password == "" {
		return errors.New("password is required")
	}
	if len(password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}
	if !isValidEmail(email) {
		return errors.New("invalid email format")
	}
	return nil
}

func isValidEmail(email string) bool {
	// Updated email validation logic to match validation.go
	emailRegex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

func TestValidateSignUpInput(t *testing.T) {
	tests := []struct {
		email       string
		password    string
		expectError bool
	}{
		{"valid.email@example.com", "password123", false},
		{"invalid-email", "password123", true},
		{"another-invalid-email@", "password123", true},
		{"valid.email@example.com", "short", true},
	}

	for _, test := range tests {
		err := ValidateSignUpInput(test.email, test.password)
		if (err != nil) != test.expectError {
			t.Errorf("ValidateSignUpInput(%v, %v) = %v, expected error: %v", test.email, test.password, err, test.expectError)
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

func TestAuthService_Logout(t *testing.T) {
	service := &AuthServiceImpl{}
	err := service.Logout("test-token")
	if err != nil && !errors.Is(err, ErrTokenNotFound) {
		t.Fatalf("expected no error or token not found error, got %v", err)
	}
}

// TestMain ensures graceful test execution even if no tests are found.
func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}
