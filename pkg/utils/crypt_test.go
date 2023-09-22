package utils_test

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"

	"go-api/pkg/utils"
)

func TestHashPassword(t *testing.T) {
	t.Run("Success hash password", func(t *testing.T) {
		password := "mysecurepassword"

		hashed, err := utils.HashPassword(password)
		assert.NoError(t, err)
		assert.NotEmpty(t, hashed)
		assert.NotEqual(t, hashed, password)
	})

	t.Run("Failed hash password", func(t *testing.T) {
		b := make([]byte, 73)
		_, err := rand.Read(b)
		assert.NoError(t, err)

		hashed, err := utils.HashPassword(string(b))
		assert.Error(t, err)
		assert.Empty(t, hashed)
	})
}

func TestComparePassword(t *testing.T) {
	password := "mysecurepassword"

	t.Run("Success match password", func(t *testing.T) {
		hashed, err := utils.HashPassword(password)
		assert.NoError(t, err)

		match := utils.ComparePassword(hashed, password)
		assert.True(t, match)
	})

	t.Run("Fail matching password", func(t *testing.T) {
		hashed, err := utils.HashPassword(password)
		assert.NoError(t, err)

		match := utils.ComparePassword(hashed, "incorrect_password")
		assert.False(t, match)
	})
}
