package auth

import (
	"errors"
	"regexp"
	"strings"
)

func ValidateSignUpInput(email, password string) error {
	// Validate email format
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(emailRegex, email)
	if !matched {
		return errors.New("invalid email format")
	}
	if !strings.Contains(email, "@") {
		return errors.New("invalid email format")
	}
	// Validate password (example: minimum 8 characters)
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	return nil
}
