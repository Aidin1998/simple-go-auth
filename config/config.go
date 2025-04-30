package config

import "github.com/spf13/viper"

type Config struct {
	Port string
}

// LoadConfig reads .env and environment variables.
// Defaults PORT to "80" if unset to match our AWS/domain.
func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	_ = viper.ReadInConfig()

	port := viper.GetString("PORT")
	if port == "" {
		port = "80" // Default to port 80 to match our AWS/domain.
	}

	return &Config{Port: port}, nil
}
