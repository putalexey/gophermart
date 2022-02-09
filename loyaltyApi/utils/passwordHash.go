package utils

import "crypto/sha256"

func PasswordHash(password string) string {
	sha := sha256.New()
	result := sha.Sum([]byte(password))
	return string(result)
}

func PasswordCheck(password, passwordHash string) bool {
	newHash := PasswordHash(password)
	return newHash == passwordHash
}
