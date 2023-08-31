package utility

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Compare the encrypted and the user provided passwords.
func VerifyPassword(hashedPassword string, checkedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(checkedPassword))
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("could not hash the password %w", err)
	}
	return string(hashedPassword), nil
}
