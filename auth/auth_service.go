package auth

import (
	"context"
	"errors"

	"my-go-project/aws"
)

// AuthTokens represents the tokens returned after authentication.
type AuthTokens struct {
	AccessToken  string `json:"access_token"`
	IdToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// AuthServiceImpl implements the authentication service using Cognito.
type AuthServiceImpl struct {
	CognitoClient *aws.CognitoClient
}

// NewAuthServiceImpl creates a new instance of AuthServiceImpl.
func NewAuthServiceImpl(client *aws.CognitoClient) *AuthServiceImpl {
	return &AuthServiceImpl{CognitoClient: client}
}

// SignUp registers a new user in Cognito.
func (s *AuthServiceImpl) SignUp(ctx context.Context, username, password, email string) error {
	return s.CognitoClient.SignUp(ctx, username, password, email)
}

// SignIn authenticates a user and returns their tokens.
func (s *AuthServiceImpl) SignIn(ctx context.Context, username, password string) (*AuthTokens, error) {
	out, err := s.CognitoClient.SignIn(ctx, username, password)
	if err != nil {
		return nil, err
	}
	ar := out.AuthenticationResult
	if ar == nil {
		return nil, errors.New("authentication failed: no result returned")
	}
	return &AuthTokens{
		AccessToken:  *ar.AccessToken,
		IdToken:      *ar.IdToken,
		RefreshToken: *ar.RefreshToken,
		ExpiresIn:    int64(ar.ExpiresIn),
		TokenType:    *ar.TokenType,
	}, nil
}

// SignOut revokes the user’s session in Cognito.
func (s *AuthServiceImpl) SignOut(ctx context.Context, accessToken string) error {
	return s.CognitoClient.SignOut(ctx, accessToken)
}

// ValidateToken checks the access token via Cognito GetUser API.
func (s *AuthServiceImpl) ValidateToken(ctx context.Context, token string) error {
	_, err := s.CognitoClient.GetUser(ctx, token)
	return err
}
