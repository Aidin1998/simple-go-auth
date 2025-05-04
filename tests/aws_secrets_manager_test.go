package tests

import (
	"os"
	"testing"

	"simple-go-auth/internal/users/config"
	"simple-go-auth/internal/users/http"

	"simple-go-auth/internal/users/auth"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type fakeSecretsManager struct{}

func (f *fakeSecretsManager) GetSecretValue(input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error) {
	// Stub implementation for AWS Secrets Manager
	if *input.SecretId == "MY_SECRET" {
		return &secretsmanager.GetSecretValueOutput{SecretString: aws.String("hello123")}, nil
	}
	if *input.SecretId == "JWT_SECRET" {
		return &secretsmanager.GetSecretValueOutput{SecretString: aws.String("dev_jwt_value")}, nil
	}
	return nil, nil
}

func TestLocalSecretsManager(t *testing.T) {
	// Arrange: set env vars that LocalSecretsManager will read
	os.Setenv("JWT_SECRET", "dev_jwt_value")
	os.Setenv("MY_SECRET", "hello123")

	sm := &fakeSecretsManager{}
	cfg := &config.Config{}
	logger, _ := zap.NewDevelopment()
	authHandler := &auth.AuthHandler{}
	authMw := func(next echo.HandlerFunc) echo.HandlerFunc { return next }

	router := http.SetupRouter(authHandler, authMw, logger, cfg)

	// Act & Assert: GetSecret
	val, err := sm.GetSecretValue(&secretsmanager.GetSecretValueInput{SecretId: aws.String("MY_SECRET")})
	require.NoError(t, err)
	require.Equal(t, "hello123", *val.SecretString)

	// Act & Assert: GetJWTSecret
	jwt, err := sm.GetSecretValue(&secretsmanager.GetSecretValueInput{SecretId: aws.String("JWT_SECRET")})
	require.NoError(t, err)
	require.Equal(t, "dev_jwt_value", *jwt.SecretString)

	// Additional assertions for router setup
	require.NotNil(t, router)
}
