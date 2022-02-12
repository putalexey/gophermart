package utils

import (
	"crypto/sha256"
	"fmt"
)

func PasswordHash(password string) string {
	sha := sha256.New()
	sha.Write([]byte(password))
	return fmt.Sprintf("%x", sha.Sum(nil))
}

func PasswordCheck(password, passwordHash string) bool {
	newHash := PasswordHash(password)
	return newHash == passwordHash
}
