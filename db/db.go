package db

import (
	"fmt"
	"time"

	"my-go-project/config"

	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitDB opens a Postgres connection with OpenTelemetry instrumentation
// and configures the connection pool.
func InitDB(cfg *config.Config) (*gorm.DB, error) {
	// Build the DSN
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	// Open a database/sql DB wrapped by otelsql.
	// It will auto-instrument all Exec/Query calls.
	sqlDB, err := otelsql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open instrumented SQL driver: %w", err)
	}

	// Pass the instrumented *sql.DB into GORM
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open gorm.DB: %w", err)
	}

	// Pull out the underlying *sql.DB to tune the pool
	underlying, err := gormDB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Configure connection pooling
	underlying.SetMaxOpenConns(cfg.DBMaxOpenConns)
	underlying.SetMaxIdleConns(cfg.DBMaxIdleConns)
	underlying.SetConnMaxLifetime(time.Duration(cfg.DBConnMaxLifetime) * time.Second)

	return gormDB, nil
}
