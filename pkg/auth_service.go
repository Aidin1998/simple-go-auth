package auth

import (
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"time"
	"errors"
)

// AuthServiceImpl implements the AuthService interface.
type AuthServiceImpl struct {
	SecretKey         string
	AccessTokenExpiry int // Token expiration in seconds
}

// CreateToken generates access and refresh tokens for a user.
func (a *AuthServiceImpl) CreateToken(userID string) (*TokenDetails, error) {
	td := &TokenDetails{}
	td.AtExpires = time.Now().Add(time.Duration(a.AccessTokenExpiry) * time.Second).Unix()
	td.AccessUUID = userID + "-access"
	td.RtExpires = time.Now().Add(7 * 24 * time.Hour).Unix()
	td.RefreshUUID = userID + "-refresh"

	var err error

	// Create Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUUID
	atClaims["user_id"] = userID
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(a.SecretKey))
	if err != nil {
		return nil, err
	}

	// Create Refresh Token
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUUID
	rtClaims["user_id"] = userID
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(a.SecretKey))
	if err != nil {
		return nil, err
	}

	return td, nil
}

// ValidateToken validates a JWT and returns the user ID if valid.
func (a *AuthServiceImpl) ValidateToken(token string) (string, error) {
	tokenObj, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(a.SecretKey), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := tokenObj.Claims.(jwt.MapClaims)
	if ok && tokenObj.Valid {
		userID, ok := claims["user_id"].(string)
		if !ok {
			return "", errors.New("invalid token claims")
		}
		return userID, nil
	}
	return "", errors.New("invalid token")
}

// HashPassword hashes a password using bcrypt.
func (a *AuthServiceImpl) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckPasswordHash compares a hashed password with its plaintext equivalent.
func (a *AuthServiceImpl) CheckPasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}