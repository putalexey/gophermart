package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPasswordCheck(t *testing.T) {
	t.Run("checks generated password", func(t *testing.T) {
		password := "123456ASD"
		password2 := "123456asd"
		h, err := PasswordHash(password)
		assert.NoError(t, err)
		assert.True(t, PasswordCheck(password, h))
		assert.False(t, PasswordCheck(password2, h))
	})
}

func TestPasswordHash(t *testing.T) {
	t.Run("generates different hashes for same password", func(t *testing.T) {
		h1, err := PasswordHash("123456")
		assert.NoError(t, err)
		h2, err := PasswordHash("123456")
		assert.NoError(t, err)
		assert.NotEqual(t, h1, h2)
	})
}
