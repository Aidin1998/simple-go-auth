package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type AuthServiceImpl struct {
	SecretKey         string
	AccessTokenExpiry int
}

func (s *AuthServiceImpl) GenerateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(time.Duration(s.AccessTokenExpiry) * time.Second).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.SecretKey))
}

func (s *AuthServiceImpl) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.SecretKey), nil
	})
}
