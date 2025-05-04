// config/config.go
package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Port               string // default "80"
	AWSRegion          string // default "ap-southeast-2"
	Env                string // add this
	DBHost             string
	DBUser             string
	DBPassword         string
	DBName             string
	AccessTokenExpiry  int           // default 3600
	DBMaxOpenConns     int           // max open DB connections
	DBMaxIdleConns     int           // max idle DB connections
	DBConnMaxLifetime  time.Duration // max connection lifetime
	JWTSecret          string        // populated at startup from AWS
	CognitoUserPoolID  string        // from COGNITO_USER_POOL_ID
	CognitoAppClientID string        // from COGNITO_APP_CLIENT_ID
	RecaptchaSecretKey string        // from RECAPTCHA_SECRET_KEY
	EchoReadTimeout    time.Duration // default '5s'
	EchoWriteTimeout   time.Duration // default '10s'
	MFAEnabled         bool          // default "false"
	SocialProviders    []string      // default ""
}

// LoadConfig reads .env and environment variables into Config.
func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	_ = viper.ReadInConfig()

	cfg := &Config{
		Port:               viper.GetString("PORT"),
		AWSRegion:          viper.GetString("AWS_REGION"),
		Env:                viper.GetString("ENV"), // add this
		DBHost:             viper.GetString("DB_HOST"),
		DBUser:             viper.GetString("DB_USER"),
		DBPassword:         viper.GetString("DB_PASSWORD"),
		DBName:             viper.GetString("DB_NAME"),
		AccessTokenExpiry:  viper.GetInt("ACCESS_TOKEN_EXPIRY"),
		DBMaxOpenConns:     viper.GetInt("DB_MAX_OPEN_CONNS"),
		DBMaxIdleConns:     viper.GetInt("DB_MAX_IDLE_CONNS"),
		DBConnMaxLifetime:  viper.GetDuration("DB_CONN_MAX_LIFETIME"),
		JWTSecret:          "", // will be fetched from AWS
		CognitoUserPoolID:  viper.GetString("COGNITO_USER_POOL_ID"),
		CognitoAppClientID: viper.GetString("COGNITO_APP_CLIENT_ID"),
		RecaptchaSecretKey: viper.GetString("RECAPTCHA_SECRET_KEY"),
		EchoReadTimeout:    viper.GetDuration("ECHO_READ_TIMEOUT"),
		EchoWriteTimeout:   viper.GetDuration("ECHO_WRITE_TIMEOUT"),
		MFAEnabled:         viper.GetBool("MFA_ENABLED"),
		SocialProviders:    viper.GetStringSlice("SOCIAL_PROVIDERS"),
	}

	// Fallback defaults
	if cfg.Port == "" {
		cfg.Port = "80"
	}
	if cfg.AWSRegion == "" {
		cfg.AWSRegion = "ap-southeast-2"
	}
	if cfg.Env == "" {
		cfg.Env = "production" // set default to "production"
	}
	if cfg.AccessTokenExpiry == 0 {
		cfg.AccessTokenExpiry = 3600
	}
	if cfg.CognitoUserPoolID == "" {
		cfg.CognitoUserPoolID = ""
	}
	if cfg.CognitoAppClientID == "" {
		cfg.CognitoAppClientID = ""
	}
	if cfg.EchoReadTimeout == 0 {
		cfg.EchoReadTimeout = 5 * time.Second
	}
	if cfg.EchoWriteTimeout == 0 {
		cfg.EchoWriteTimeout = 10 * time.Second
	}
	if !cfg.MFAEnabled {
		cfg.MFAEnabled = false
	}
	if len(cfg.SocialProviders) == 0 {
		cfg.SocialProviders = []string{}
	}

	return cfg, nil
}
