package auth

import (
	"context"
	"errors"
	"my-go-project/aws"
)

// AuthTokens holds the tokens returned by Cognito.
type AuthTokens struct {
	AccessToken  string `json:"access_token"`
	IdToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// AuthServiceImpl implements authentication via Cognito.
type AuthServiceImpl struct {
	CognitoClient *aws.CognitoClient
}

// NewAuthServiceImpl constructs a new AuthServiceImpl.
func NewAuthServiceImpl(c *aws.CognitoClient) *AuthServiceImpl {
	return &AuthServiceImpl{CognitoClient: c}
}

// SignUp registers a new user in Cognito.
func (s *AuthServiceImpl) SignUp(ctx context.Context, username, password, email string) error {
	return s.CognitoClient.SignUp(ctx, username, password, email)
}

// SignIn authenticates a user and returns their token set.
func (s *AuthServiceImpl) SignIn(ctx context.Context, username, password string) (*AuthTokens, error) {
	out, err := s.CognitoClient.SignIn(ctx, username, password)
	if err != nil {
		return nil, err
	}
	if out.AuthenticationResult == nil {
		return nil, errors.New("no auth result")
	}
	ar := out.AuthenticationResult
	return &AuthTokens{
		AccessToken:  *ar.AccessToken,
		IdToken:      *ar.IdToken,
		RefreshToken: *ar.RefreshToken,
		ExpiresIn:    int64(*ar.ExpiresIn),
		TokenType:    *ar.TokenType,
	}, nil
}

// SignOut revokes the user’s session in Cognito.
func (s *AuthServiceImpl) SignOut(ctx context.Context, accessToken string) error {
	return s.CognitoClient.SignOut(ctx, accessToken)
}

// ValidateToken validates the provided token.
func (svc *AuthServiceImpl) ValidateToken(ctx context.Context, token string) error {
	// Implement token validation logic here
	// Return nil if the token is valid, or an error if it is invalid
	return nil
}
