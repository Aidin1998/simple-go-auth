package tests

import (
	"os"
	"testing"

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

	// Migrate models for test
	assert.NoError(t, gormDB.AutoMigrate(&db.User{}, &db.RefreshToken{}))

	// Insert and query a user
	user := db.User{Username: "test", Email: "t@e.com"}
	assert.NoError(t, gormDB.Create(&user).Error)
	assert.NotZero(t, user.ID)
}
