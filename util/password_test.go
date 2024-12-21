package u

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	password := RandomPassword()

	hash1, err := HashPassword(password)
	require.NoError(t, err)

	err = CheckPassword(password, hash1)
	require.NoError(t, err)
	require.NotEmpty(t, hash1)

	// Compare the hash with a wrong password's hash.
	wrongPassword := RandomPassword()
	err = CheckPassword(wrongPassword, hash1)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	// Now we'll hash the password again (will provide a different hash becaue of the salt) and check the password with the new hash.
	hash2, err := HashPassword(password)
	require.NoError(t, err)

	err = CheckPassword(password, hash2)
	require.NoError(t, err)
	require.NotEmpty(t, hash2)

	// Compare the two hashes.
	require.NotEqual(t, hash1, hash2)
}
