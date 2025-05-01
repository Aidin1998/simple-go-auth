package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"my-go-project/aws"

	"github.com/labstack/echo/v4"
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
func NewAuthServiceImpl(client *aws.CognitoClient) *AuthServiceImpl {
	return &AuthServiceImpl{CognitoClient: client}
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
		return nil, errors.New("authentication failed: no result returned")
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

type AuthService interface {
	ValidateToken(token string) (bool, error)
}

func NewMiddleware(authService AuthService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get("Authorization")
			if token == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing token"})
			}

			valid, err := authService.ValidateToken(token)
			if err != nil || !valid {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
			}

			return next(c)
		}
	}
}

// Duplicate declaration removed.

// Other methods of AuthServiceImpl

// ValidateToken validates a token and returns the result.
func (a *AuthServiceImpl) ValidateToken(token string) (bool, error) {
	// Implement token validation logic here
	return a.CognitoClient.ValidateToken(token)
}

// AuthServiceImpl is the implementation of the authentication service.
// Duplicate declaration removed.

// GenerateToken generates a token for the given userID.
func (s *AuthServiceImpl) GenerateToken(userID string) (string, error) {
	// Placeholder logic for token generation
	if userID == "" {
		return "", fmt.Errorf("userID cannot be empty")
	}
	return "mockTokenFor_" + userID, nil
}
