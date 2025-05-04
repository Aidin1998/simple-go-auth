package tests

import (
	"testing"

	"simple-go-auth/internal/users/config"
	"simple-go-auth/internal/users/db"

	"github.com/stretchr/testify/require"
)

func TestInitDB(t *testing.T) {
	// Load configuration
	cfg, err := config.LoadConfig()
	require.NoError(t, err)

	// Initialize database
	gormDB, err := db.InitDB(cfg)
	require.NoError(t, err)

	// Ping the database
	sqlDB, err := gormDB.DB()
	require.NoError(t, err)
	require.NoError(t, sqlDB.Ping())
}
