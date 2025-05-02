package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	sdkconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"
)

// CognitoClient wraps AWS Cognito SDK operations.
type CognitoClient struct {
	client      *cognitoidentityprovider.Client
	userPoolID  string
	appClientID string
}

// NewCognitoClient initializes a new CognitoClient with region, userPoolID, and appClientID.
func NewCognitoClient(region, userPoolID, appClientID string) (*CognitoClient, error) {
	cfg, err := sdkconfig.LoadDefaultConfig(context.TODO(), sdkconfig.WithRegion(region))
	if err != nil {
		return nil, err
	}

	// Add OpenTelemetry instrumentation to the AWS SDK client.
	otelaws.AppendMiddlewares(&cfg.APIOptions)

	return &CognitoClient{
		client:      cognitoidentityprovider.NewFromConfig(cfg),
		userPoolID:  userPoolID,
		appClientID: appClientID,
	}, nil
}

// SignUp registers a new user in Cognito.
func (c *CognitoClient) SignUp(ctx context.Context, username, password, email string) error {
	_, err := c.client.SignUp(ctx, &cognitoidentityprovider.SignUpInput{
		ClientId: aws.String(c.appClientID),
		Username: aws.String(username),
		Password: aws.String(password),
		UserAttributes: []types.AttributeType{
			{Name: aws.String("email"), Value: aws.String(email)},
		},
	})
	return err
}

// SignIn authenticates a user with Cognito and returns the raw SDK output.
func (c *CognitoClient) SignIn(ctx context.Context, username, password string) (*cognitoidentityprovider.InitiateAuthOutput, error) {
	return c.client.InitiateAuth(ctx, &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: types.AuthFlowTypeUserPasswordAuth,
		ClientId: aws.String(c.appClientID),
		AuthParameters: map[string]string{
			"USERNAME": username,
			"PASSWORD": password,
		},
	})
}

// SignOut revokes a user's session in Cognito.
func (c *CognitoClient) SignOut(ctx context.Context, accessToken string) error {
	_, err := c.client.GlobalSignOut(ctx, &cognitoidentityprovider.GlobalSignOutInput{
		AccessToken: aws.String(accessToken),
	})
	return err
}

// GetUser fetches the user associated with the access token—used for validating tokens.
func (c *CognitoClient) GetUser(ctx context.Context, accessToken string) (*cognitoidentityprovider.GetUserOutput, error) {
	return c.client.GetUser(ctx, &cognitoidentityprovider.GetUserInput{
		AccessToken: aws.String(accessToken),
	})
}

// ConfirmSignUp verifies a user’s email/SMS code.
func (c *CognitoClient) ConfirmSignUp(ctx context.Context, username, code string) error {
	_, err := c.client.ConfirmSignUp(ctx, &cognitoidentityprovider.ConfirmSignUpInput{
		ClientId:         aws.String(c.appClientID),
		Username:         aws.String(username),
		ConfirmationCode: aws.String(code),
	})
	return err
}

// RefreshAuth uses the REFRESH_TOKEN_AUTH flow to rotate tokens.
func (c *CognitoClient) RefreshAuth(ctx context.Context, refreshToken string) (*cognitoidentityprovider.InitiateAuthOutput, error) {
	return c.client.InitiateAuth(ctx, &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: types.AuthFlowTypeRefreshTokenAuth,
		ClientId: aws.String(c.appClientID),
		AuthParameters: map[string]string{
			"REFRESH_TOKEN": refreshToken,
		},
	})
}
