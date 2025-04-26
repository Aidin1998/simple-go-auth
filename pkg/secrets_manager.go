package auth

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// SecretsManager provides methods to fetch secrets from AWS Secrets Manager.
type SecretsManager struct {
	client *secretsmanager.Client
}

// NewSecretsManager initializes a new SecretsManager instance.
func NewSecretsManager() (*SecretsManager, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	return &SecretsManager{
		client: secretsmanager.NewFromConfig(cfg),
	}, nil
}

// GetSecret retrieves a secret value by its name.
func (sm *SecretsManager) GetSecret(secretName string) (map[string]string, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId: &secretName,
	}

	result, err := sm.client.GetSecretValue(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	if result.SecretString == nil {
		return nil, errors.New("secret string is nil")
	}

	var secretMap map[string]string
	if err := json.Unmarshal([]byte(*result.SecretString), &secretMap); err != nil {
		return nil, err
	}

	return secretMap, nil
}