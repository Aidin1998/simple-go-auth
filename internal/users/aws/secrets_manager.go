package aws

// SecretsManager defines how we fetch secrets.
type SecretsManager interface {
	GetSecret(name string) (string, error)
	GetJWTSecret() (string, error)
}
