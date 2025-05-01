package tests

import (
	"os"
	"testing"

	"my-go-project/aws"

	"github.com/stretchr/testify/assert"
)

func TestLocalSecretsManager(t *testing.T) {
	// Arrange: set env vars that LocalSecretsManager will read
	os.Setenv("JWT_SECRET", "dev_jwt_value")
	os.Setenv("MY_SECRET", "hello123")

	sm := aws.NewLocalSecretsManager()

	// Act & Assert: GetSecret
	val, err := sm.GetSecret("MY_SECRET")
	assert.NoError(t, err)
	assert.Equal(t, "hello123", val)

	// Act & Assert: GetJWTSecret
	jwt, err := sm.GetJWTSecret()
	assert.NoError(t, err)
	assert.Equal(t, "dev_jwt_value", jwt)
}
