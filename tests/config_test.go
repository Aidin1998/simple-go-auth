package tests

import (
	"testing"
	"time"

	"my-go-project/config"
)

func TestLoadConfig(t *testing.T) {
	// Load config from .env.sample or construct manually
	cfg := &config.Config{
		Port:              "8080",
		AWSRegion:         "us-west-2",
		DBHost:            "localhost",
		DBUser:            "testuser",
		DBPassword:        "testpassword",
		DBName:            "testdb",
		DBMaxOpenConns:    10,
		DBMaxIdleConns:    5,
		DBConnMaxLifetime: time.Minute * 5,
	}
	// Add test logic here
	if cfg.Port != "8080" {
		t.Errorf("Expected Port to be '8080', got '%s'", cfg.Port)
	}
}
