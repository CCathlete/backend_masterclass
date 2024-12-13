package sqlc_test

// import (
// 	"backend-masterclass/db/sqlc"
// 	"backend-masterclass/util"
// 	"context"
// 	"database/sql"
// 	"testing"
// 	"time"

// 	"github.com/stretchr/testify/require"
// )

// // We pass in the test to run testify/require functions.
// // This function is a validation a preparatino for each
// // test written here.
// func createRandomAccount(t *testing.T) sqlc.Account {
// 	arg := sqlc.CreateAccountParams{
// 		Owner:    util.RandomOwner(),
// 		Balance:  util.RandomMoney(),
// 		Currency: util.RandCurrency(),
// 	}

// 	account, err := testQueries.CreateAccount(context.Background(), arg)
// 	require.NoError(t, err)
// 	require.NotEmpty(t, account)

// 	require.Equal(t, arg.Owner, account.Owner)
// 	require.Equal(t, arg.Balance, account.Balance)
// 	require.Equal(t, arg.Currency, account.Currency)

// 	require.NotZero(t, account.ID)
// 	require.NotZero(t, account.CreatedAt)

// 	return account
// }

// func TestCreateAccount(t *testing.T) {
// 	createRandomAccount(t)
// }

// func TestGetAccount(t *testing.T) {
// 	account := createRandomAccount(t)
// 	result, err := testQueries.GetAccount(context.Background(), account.ID)
// 	require.NoError(t, err)
// 	require.NotEmpty(t, result)

// 	require.Equal(t, account.ID, result.ID)
// 	require.Equal(t, account.Owner, result.Owner)
// 	require.Equal(t, account.Balance, result.Balance)
// 	require.Equal(t, account.Currency, result.Currency)
// 	// There might be a short delay from creation of the random account
// 	// to its storage in the DB and we don't this to fail the test.
// 	require.WithinDuration(t, account.CreatedAt, result.CreatedAt, time.Second)
// }

// func TestUpdateAccount(t *testing.T) {
// 	account := createRandomAccount(t)

// 	arg := sqlc.UpdateAccountParams{
// 		ID:      account.ID,
// 		Balance: util.RandomMoney(),
// 	}

// 	result, err := testQueries.UpdateAccount(context.Background(), arg)
// 	require.NoError(t, err)
// 	require.NotEmpty(t, result)

// 	require.Equal(t, account.ID, result.ID)
// 	require.Equal(t, account.Owner, result.Owner)
// 	// The balance should change to the new value.
// 	require.Equal(t, arg.Balance, result.Balance)
// 	require.Equal(t, account.Currency, result.Currency)
// 	// There might be a short delay from creation of the random account
// 	// to its storage in the DB and we don't this to fail the test.
// 	require.WithinDuration(t, account.CreatedAt, result.CreatedAt, time.Second)
// }

// func TestDeleteAccount(t *testing.T) {
// 	account := createRandomAccount(t)
// 	err := testQueries.DeleteAccount(context.Background(), account.ID)
// 	require.NoError(t, err)

// 	// Validation that the account was truly deleted
// 	valAccount, err := testQueries.GetAccount(context.Background(), account.ID)
// 	// We want an error to occur here.
// 	require.EqualError(t, err, sql.ErrNoRows.Error())
// 	require.Empty(t, valAccount)
// }

// func TestListAccounts(t *testing.T) {
// 	for i := 0; i < 10; i++ {
// 		createRandomAccount(t)
// 	}

// 	arg := sqlc.ListAccountParams{
// 		Limit:  5,
// 		Offset: 5, // Skips 5 matches before returning values.
// 	}
// 	results, err := testQueries.ListAccount(context.Background(), arg)
// 	require.NoError(t, err)

// 	for _, result := range results {
// 		require.NotEmpty(t, result)
// 	}
// }
