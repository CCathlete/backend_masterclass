package u

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (hash string, err error) {
	bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		err = fmt.Errorf("failed to hash password: %w", err)
		return
	}
	hash = string(bcryptPassword)
	return
}

func CheckPassword(password, hash string) (err error) {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
