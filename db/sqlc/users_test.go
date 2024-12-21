package sqlc

import (
	u "backend-masterclass/util"
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// We pass in the test to run testify/require functions.
// This function is a validation a preparatino for each
// test written here.
func createRandomUser(t *testing.T) User {
	arg := CreateUserParams{
		FullName: u.RandomFullName(),
		Username: u.RandomUsername(),
		Email:    u.RandomEmail(),
	}
	hashedPassword, err := u.HashPassword(u.RandomPassword())
	require.NoError(t, err)
	arg.HashedPassword = hashedPassword

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.PasswordChangedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user := createRandomUser(t)
	result, err := testQueries.GetUser(context.Background(), user.Username)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.Equal(t, user.Username, result.Username)
	require.Equal(t, user.FullName, result.FullName)
	require.Equal(t, user.HashedPassword, result.HashedPassword)
	require.Equal(t, user.FullName, result.FullName)
	// There might be a short delay from creation of the random user
	// to its storage in the DB and we don't this to fail the test.
	require.WithinDuration(t, user.PasswordChangedAt, result.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user.CreatedAt, result.CreatedAt, time.Second)
}

func TestUpdateUser(t *testing.T) {
	user := createRandomUser(t)
	newUsername := u.RandomUsername()

	arg := UpdateUserParams{
		Username:       user.Username,
		NewUsername:    newUsername,
		FullName:       user.FullName,
		HashedPassword: user.HashedPassword,
		Email:          user.Email,
	}

	result, err := testQueries.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.Equal(t, newUsername, result.Username)
	require.Equal(t, user.FullName, result.FullName)
	// The balance should change to the new value.
	require.Equal(t, arg.HashedPassword, result.HashedPassword)
	require.Equal(t, user.FullName, result.FullName)
	// NOTE: This is not true in case we change the password.
	require.WithinDuration(t, user.PasswordChangedAt, result.PasswordChangedAt, time.Second)
}

func TestDeleteUser(t *testing.T) {
	user := createRandomUser(t)
	err := testQueries.DeleteUser(context.Background(), user.Username)
	require.NoError(t, err)

	// Validation that the user was truly deleted
	valUser, err := testQueries.GetUser(context.Background(), user.Username)
	// We want an error to occur here.
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, valUser)
}

func TestListUsers(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomUser(t)
	}

	arg := ListUsersParams{
		Limit:  5,
		Offset: 5, // Skips 5 matches before returning values.
	}
	results, err := testQueries.ListUsers(context.Background(), arg)
	require.NoError(t, err)

	for _, result := range results {
		require.NotEmpty(t, result)
	}
}
