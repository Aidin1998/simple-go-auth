package tests

import (
	"context"
	"os"
	"testing"

	"my-go-project/aws"

	"github.com/stretchr/testify/assert"
)

func TestCognitoLocalSignUpSignIn(t *testing.T) {
	// Arrange: point AWS SDK at Localstack
	os.Setenv("AWS_REGION", "us-east-1")
	os.Setenv("COGNITO_ENDPOINT", "http://localhost:4566") // Localstack endpoint
	client, err := aws.NewCognitoClient("us-east-1", "test-pool-id", "test-client-id")
	assert.NoError(t, err)

	// Act: SignUp & SignIn must succeed
	err = client.SignUp(context.TODO(), "user1", "P@ssw0rd!", "email@example.com")
	assert.NoError(t, err)

	authOut, err := client.SignIn(context.TODO(), "user1", "P@ssw0rd!")
	assert.NoError(t, err)
	assert.NotEmpty(t, authOut.AuthenticationResult.AccessToken)
}
