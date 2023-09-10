package utility

import "golang.org/x/crypto/bcrypt"

// Compare the encrypted and the user provided passwords.
func VerifyPassword(hashedPassword string, checkedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(checkedPassword))
}
