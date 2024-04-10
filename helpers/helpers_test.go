package helpers

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T)  {
	password := RandomString(6)

	hashedPassword, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	require.NoError(t, CheckPassword(password, hashedPassword))

	// wrong password
	wrongPassword := RandomString(6)
	require.Error(t, CheckPassword(wrongPassword, hashedPassword), bcrypt.ErrMismatchedHashAndPassword.Error())
	
	// same passwords should give different hash
	hashedPassword2, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword2)
	require.NotEqual(t, hashedPassword, hashedPassword2)
}