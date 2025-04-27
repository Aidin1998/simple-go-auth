package auth

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// GetSecret retrieves a secret value from AWS Secrets Manager.
func GetSecret(secretName string) (string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return "", errors.New("unable to load AWS SDK config: " + err.Error())
	}

	client := secretsmanager.NewFromConfig(cfg)
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	}

	result, err := client.GetSecretValue(context.TODO(), input)
	if err != nil {
		return "", errors.New("unable to retrieve secret: " + err.Error())
	}

	if result.SecretString == nil {
		return "", errors.New("secret value is empty")
	}

	return *result.SecretString, nil
}

// GetJWTSecret retrieves the JWT secret key from AWS Secrets Manager.
func GetJWTSecret() (string, error) {
	secretName := "jwtSecretKey" // Replace with your actual secret name
	return GetSecret(secretName)
}
