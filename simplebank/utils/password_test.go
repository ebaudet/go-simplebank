package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	password := RandomPassword(6, 25)
	hash1, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hash1)

	hash2, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hash2)
	require.NotEqual(t, hash1, hash2)

	err = CheckPassword(password, hash1)
	require.NoError(t, err)

	wrongPassword := RandomPassword(6, 25)
	err = CheckPassword(wrongPassword, hash1)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
}
