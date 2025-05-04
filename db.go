package db

import (
	"database/sql"
	"simple-go-auth/internal/users/config"
	"time"
)

func InitializeDB(cfg *config.Config) (*sql.DB, error) {
	// ...existing code to initialize the database connection...
	db, err := sql.Open("postgres", "your_connection_string_here")
	if err != nil {
		return nil, err
	}

	// Set database connection pooling configurations
	db.SetMaxOpenConns(25)                  // Maximum number of open connections
	db.SetMaxIdleConns(10)                  // Maximum number of idle connections
	db.SetConnMaxLifetime(30 * time.Minute) // Maximum lifetime of a connection

	return db, nil
}
