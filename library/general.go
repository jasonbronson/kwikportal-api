package library

import (
	"golang.org/x/crypto/bcrypt"
)

// GeneratePassword generates a bcrypt hash for the given password.
func GeneratePassword(password string) string {
	// Hash the password before storing it
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return ""
	}
	return string(hashedPassword)
}
