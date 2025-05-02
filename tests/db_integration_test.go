package tests

import (
	"fmt"
	"os"
	"testing"
	"time"

	"my-go-project/config"
	"my-go-project/db"

	"github.com/stretchr/testify/assert"
)

func TestPostgresConnection(t *testing.T) {
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "password")
	os.Setenv("DB_NAME", "testdb")

	cfg, err := config.LoadConfig()
	assert.NoError(t, err)

	gormDB, err := db.InitDB(cfg)
	assert.NoError(t, err)

	// Clean up the database before running the test
	gormDB.Exec("DELETE FROM users")

	// Migrate models for test
	assert.NoError(t, gormDB.AutoMigrate(&db.User{}, &db.RefreshToken{}))

	// Use a unique username to avoid conflicts
	username := fmt.Sprintf("test_user_%d", time.Now().UnixNano())

	// Insert and query a user
	user := db.User{
		Username:  username,
		Email:     "t@e.com",
		Password:  "",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	assert.NoError(t, gormDB.Create(&user).Error)
	assert.NotZero(t, user.ID)
}
