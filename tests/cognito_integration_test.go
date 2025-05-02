package tests

import (
	"context"
	"testing"
	"time"

	"my-go-project/aws"
	"my-go-project/config"

	"github.com/stretchr/testify/require"
)

func TestCognitoFlows(t *testing.T) {
	// Load configuration and override endpoints for LocalStack
	cfg, _ := config.LoadConfig()
	cognitoClient, err := aws.NewCognitoClient(
		cfg.AWSRegion,
		"us-east-1_userpoolid",
		"localstack-appclient",
	)
	require.NoError(t, err)

	ctx := context.Background()
	username, pw, email := "testu", "Password1!", "t@e.com"

	// SignUp
	require.NoError(t, cognitoClient.SignUp(ctx, username, pw, email))
	time.Sleep(2 * time.Second) // LocalStack eventual consistency

	// ConfirmSignUp
	require.NoError(t, cognitoClient.ConfirmSignUp(ctx, username, "123456"))

	// SignIn
	out, err := cognitoClient.SignIn(ctx, username, pw)
	require.NoError(t, err)
	require.NotNil(t, out.AuthenticationResult)

	// SignOut
	accessToken := *out.AuthenticationResult.AccessToken
	require.NoError(t, cognitoClient.SignOut(ctx, accessToken))
}
