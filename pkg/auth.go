package auth

import (
	"errors"
	"regexp"
)

// ValidateSignUpInput validates the user input for sign-up.
func ValidateSignUpInput(email, password string) error {
	if email == "" || password == "" {
		return errors.New("email and password cannot be empty")
	}

	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	if matched, _ := regexp.MatchString(emailRegex, email); !matched {
		return errors.New("invalid email format")
	}

	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	return nil
}