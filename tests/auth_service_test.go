package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"github.com/stretchr/testify/assert"
)

func TestGenerateJWT(t *testing.T) {
	jwtSecretKey = "testsecret"
	userID := "12345"
	expiration := time.Hour
	token, err := GenerateJWT(userID, expiration)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := ValidateJWT(token)
	assert.NoError(t, err)
	assert.Equal(t, userID, claims["user_id"])
}

func TestValidateJWT_InvalidToken(t *testing.T) {
	jwtSecretKey = "testsecret"
	invalidToken := "invalid.token.string"

	_, err := ValidateJWT(invalidToken)
	assert.Error(t, err)
}

func TestValidateSignUpInput(t *testing.T) {
	err := ValidateSignUpInput("test@example.com", "password123")
	assert.NoError(t, err)

	err = ValidateSignUpInput("invalid-email", "password123")
	assert.Error(t, err)

	err = ValidateSignUpInput("test@example.com", "short")
	assert.Error(t, err)
}

func TestHashPassword(t *testing.T) {
	password := "password123"
	hashedPassword, err := HashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)

	err = CheckPassword(hashedPassword, password)
	assert.NoError(t, err)

	err = CheckPassword(hashedPassword, "wrongpassword")
	assert.Error(t, err)
}