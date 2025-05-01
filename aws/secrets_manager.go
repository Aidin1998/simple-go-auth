package aws

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	sdkconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// SecretsManager defines the interface for fetching secrets.
type SecretsManager interface {
	GetSecret(name string) (string, error)
	GetJWTSecret() (string, error)
}

// AWSSecretsManager is the AWS-backed implementation.
type AWSSecretsManager struct {
	region string
}

// NewAWSSecretsManager creates a SecretsManager for the given AWS region.
func NewAWSSecretsManager(region string) *AWSSecretsManager {
	return &AWSSecretsManager{region: region}
}

// GetSecret retrieves a secret value from AWS Secrets Manager.
func (a *AWSSecretsManager) GetSecret(secretName string) (string, error) {
	// Load AWS SDK config with the specified region
	cfg, err := sdkconfig.LoadDefaultConfig(context.TODO(),
		sdkconfig.WithRegion(a.region),
	)
	if err != nil {
		return "", errors.New("unable to load AWS SDK config: " + err.Error())
	}

	// Create a Secrets Manager client
	client := secretsmanager.NewFromConfig(cfg)

	// Fetch the secret
	output, err := client.GetSecretValue(context.TODO(), &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	})
	if err != nil {
		return "", errors.New("unable to retrieve secret: " + err.Error())
	}
	if output.SecretString == nil {
		return "", errors.New("secret value is empty")
	}
	return *output.SecretString, nil
}

// GetJWTSecret retrieves the JWT secret key from AWS Secrets Manager.
func (a *AWSSecretsManager) GetJWTSecret() (string, error) {
	const jwtSecretName = "jwtSecretKey" // change this to your actual secret name
	return a.GetSecret(jwtSecretName)
}

// Compile-time check that AWSSecretsManager implements SecretsManager.
var _ SecretsManager = (*AWSSecretsManager)(nil)
