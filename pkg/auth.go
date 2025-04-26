package auth

import (
	"time"
)

// TokenDetails holds the details of access and refresh tokens.
type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUUID   string
	RefreshUUID  string
	AtExpires    int64
	RtExpires    int64
}

// User represents a user in the system.
type User struct {
	ID       string
	Email    string
	Password string
}

// AuthService defines the interface for authentication services.
type AuthService interface {
	CreateToken(userID string) (*TokenDetails, error)
	ValidateToken(token string) (string, error)
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
}