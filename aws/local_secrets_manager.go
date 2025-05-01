package aws

import (
	"fmt"

	"github.com/spf13/viper"
)

// LocalSecretsManager reads secrets from Viper (env/.env).
type LocalSecretsManager struct{}

// NewLocalSecretsManager creates a new LocalSecretsManager.
func NewLocalSecretsManager() *LocalSecretsManager {
	return &LocalSecretsManager{}
}

// GetSecret retrieves a secret value from Viper.
func (l *LocalSecretsManager) GetSecret(name string) (string, error) {
	val := viper.GetString(name)
	if val == "" {
		return "", fmt.Errorf("secret %s not set", name)
	}
	return val, nil
}

// GetJWTSecret retrieves the JWT secret key from Viper.
func (l *LocalSecretsManager) GetJWTSecret() (string, error) {
	return l.GetSecret("JWT_SECRET")
}

// Compile-time check that LocalSecretsManager implements SecretsManager.
var _ SecretsManager = (*LocalSecretsManager)(nil)
