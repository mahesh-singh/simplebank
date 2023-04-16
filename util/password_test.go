package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	password := RandomString(10)
	hasedPassword, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hasedPassword)

	err = CheckPassword(hasedPassword, password)
	require.NoError(t, err)

	wrongPassword := RandomString(10)
	err = CheckPassword(hasedPassword, wrongPassword)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	hasedPassword1, err := HashPassword(password)
	require.NotEqual(t, hasedPassword, hasedPassword1)
}
