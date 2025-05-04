package tests

import (
	"encoding/json"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"simple-go-auth/internal/users/config"
	"simple-go-auth/internal/users/http"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestLoadConfig(t *testing.T) {
	// 1. Create a temporary .env file for this test
	envContent := "" +
		"PORT=1234\n" +
		"AWS_REGION=foo\n" +
		"ACCESS_TOKEN_EXPIRY=99\n" +
		"DB_HOST=bar\n" +
		"DB_USER=baz\n" +
		"DB_PASSWORD=qux\n" +
		"DB_NAME=mydb\n" +
		"DB_MAX_OPEN_CONNS=10\n" +
		"DB_MAX_IDLE_CONNS=5\n" +
		"DB_CONN_MAX_LIFETIME=2h\n"

	err := os.WriteFile(".env", []byte(envContent), 0644)
	assert.NoError(t, err)
	defer os.Remove(".env")

	// 2. Load the configuration
	cfg, err := config.LoadConfig()
	assert.NoError(t, err)

	// 3. Assert the loaded values
	assert.Equal(t, "1234", cfg.Port)
	assert.Equal(t, "foo", cfg.AWSRegion)
	assert.Equal(t, 99, cfg.AccessTokenExpiry)
	assert.Equal(t, "bar", cfg.DBHost)
	assert.Equal(t, "baz", cfg.DBUser)
	assert.Equal(t, "qux", cfg.DBPassword)
	assert.Equal(t, "mydb", cfg.DBName)
	assert.Equal(t, 10, cfg.DBMaxOpenConns)
	assert.Equal(t, 5, cfg.DBMaxIdleConns)
	assert.Equal(t, 2*time.Hour, cfg.DBConnMaxLifetime)
}

func TestConfigEndpoint(t *testing.T) {
	// Load a known test .env file via Viper
	cfg, err := config.LoadConfig()
	require.NoError(t, err)

	// Mount router
	router := http.SetupRouter(nil, nil, zap.NewNop(), cfg)
	req := httptest.NewRequest("GET", "/config", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// Decode and compare
	var got config.Config
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	require.Equal(t, cfg, &got)
}
