package auth

import (
	"errors"
	"regexp"
)

func ValidateSignUpInput(email, password string) error {
	// Stronger email validation
	emailRegex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)

	if !re.MatchString(email) {
		return errors.New("invalid email format")
	}

	// Validate password length
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	return nil
}
