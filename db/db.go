package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// User represents the user model for the database.
type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique;not null"`
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
}

// InitDB initializes the database connection and migrates the User model.
func InitDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("auth.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto-migrate the User model
	if err := db.AutoMigrate(&User{}); err != nil {
		return nil, err
	}

	return db, nil
}
