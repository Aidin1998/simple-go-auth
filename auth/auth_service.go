// auth/auth_service.go
package auth

import (
	"context"
	"errors"
	"time"

	"my-go-project/aws"
	"my-go-project/db"

	"gorm.io/gorm"
)

// AuthTokens represents the tokens returned after authentication.
type AuthTokens struct {
	AccessToken  string `json:"access_token"`
	IdToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// AuthServiceImpl implements the authentication service using Cognito and GORM.
type AuthServiceImpl struct {
	CognitoClient *aws.CognitoClient
	DB            *gorm.DB
}

// NewAuthServiceImpl creates a new instance of AuthServiceImpl.
func NewAuthServiceImpl(client *aws.CognitoClient, db *gorm.DB) *AuthServiceImpl {
	return &AuthServiceImpl{
		CognitoClient: client,
		DB:            db,
	}
}

// SignUp registers a new user in Cognito and persists a User record.
func (s *AuthServiceImpl) SignUp(ctx context.Context, username, password, email string) error {
	// 1) Create in Cognito
	if err := s.CognitoClient.SignUp(ctx, username, password, email); err != nil {
		return err
	}
	// 2) Persist in local DB
	user := &db.User{Username: username, Email: email}
	if err := s.DB.Create(user).Error; err != nil {
		return err
	}
	return nil
}

// SignIn authenticates a user via Cognito, persists the refresh token, and returns tokens.
func (s *AuthServiceImpl) SignIn(ctx context.Context, username, password string) (*AuthTokens, error) {
	// 1) Cognito auth
	out, err := s.CognitoClient.SignIn(ctx, username, password)
	if err != nil {
		return nil, err
	}
	ar := out.AuthenticationResult
	if ar == nil {
		return nil, errors.New("authentication failed: no result returned")
	}

	// 2) Build AuthTokens
	tokens := &AuthTokens{
		AccessToken:  *ar.AccessToken,
		IdToken:      *ar.IdToken,
		RefreshToken: *ar.RefreshToken,
		ExpiresIn:    int64(ar.ExpiresIn),
		TokenType:    *ar.TokenType,
	}

	// 3) Persist refresh token
	var user db.User
	if err := s.DB.
		Where("username = ?", username).
		First(&user).Error; err != nil {
		return nil, err
	}
	rt := &db.RefreshToken{
		UserID:    user.ID,
		Token:     tokens.RefreshToken,
		ExpiresAt: time.Now().Add(time.Duration(tokens.ExpiresIn) * time.Second),
	}
	if err := s.DB.Create(rt).Error; err != nil {
		return nil, err
	}

	return tokens, nil
}

// SignOut revokes the session in Cognito and marks the refresh token revoked in DB.
func (s *AuthServiceImpl) SignOut(ctx context.Context, accessToken string) error {
	// 1) Cognito global sign-out
	if err := s.CognitoClient.SignOut(ctx, accessToken); err != nil {
		return err
	}
	// 2) Revoke in local DB
	return s.DB.
		Model(&db.RefreshToken{}).
		Where("token = ?", accessToken).
		Update("revoked", true).
		Error
}

// ValidateToken checks the access token via Cognito GetUser API.
func (s *AuthServiceImpl) ValidateToken(ctx context.Context, token string) error {
	_, err := s.CognitoClient.GetUser(ctx, token)
	return err
}

// ConfirmSignUp confirms a user's sign-up using a verification code.
func (s *AuthServiceImpl) ConfirmSignUp(ctx context.Context, username, code string) error {
	return s.CognitoClient.ConfirmSignUp(ctx, username, code)
}

// RefreshTokens handles refresh‐token rotation and revocation.
func (s *AuthServiceImpl) RefreshTokens(ctx context.Context, refreshToken string) (*AuthTokens, error) {
	// 1) Lookup existing, non‐revoked token
	var rt db.RefreshToken
	if err := s.DB.
		Where("token = ? AND revoked = false", refreshToken).
		First(&rt).Error; err != nil {
		return nil, err
	}

	// 2) Rotate via Cognito
	out, err := s.CognitoClient.RefreshAuth(ctx, refreshToken)
	if err != nil {
		return nil, err
	}
	ar := out.AuthenticationResult
	if ar == nil {
		return nil, errors.New("authentication failed: no result returned")
	}

	// 3) Revoke old and record previous token
	if err := s.DB.Model(&rt).Updates(db.RefreshToken{
		Revoked:       true,
		PreviousToken: rt.Token,
	}).Error; err != nil {
		return nil, err
	}

	// 4) Insert new refresh token record
	newRT := &db.RefreshToken{
		UserID:        rt.UserID,
		Token:         *ar.RefreshToken,
		ExpiresAt:     time.Now().Add(time.Duration(ar.ExpiresIn) * time.Second),
		PreviousToken: rt.Token,
	}
	if err := s.DB.Create(newRT).Error; err != nil {
		return nil, err
	}

	// 5) Return the refreshed tokens
	return &AuthTokens{
		AccessToken:  *ar.AccessToken,
		IdToken:      *ar.IdToken,
		RefreshToken: *ar.RefreshToken,
		ExpiresIn:    int64(ar.ExpiresIn),
		TokenType:    *ar.TokenType,
	}, nil
}
