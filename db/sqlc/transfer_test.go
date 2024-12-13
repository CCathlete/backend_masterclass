package sqlc_test

import (
	"backend-masterclass/db/sqlc"
	"backend-masterclass/util"
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// We pass in the test to run testify/require functions.
// This function is a validation a preparatino for each
// test written here.
func createRandomTransfer(t *testing.T) sqlc.Transfer {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	arg := sqlc.CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomMoney(),
	}

	Transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, Transfer)

	require.Equal(t, account1.ID, Transfer.FromAccountID)
	require.Equal(t, account2.ID, Transfer.ToAccountID)
	require.Equal(t, arg.Amount, Transfer.Amount)

	require.NotZero(t, Transfer.ID)
	require.NotZero(t, Transfer.CreatedAt)

	return Transfer
}

func TestCreateTransfer(t *testing.T) {
	createRandomTransfer(t)
}

func TestGetTransfer(t *testing.T) {
	Transfer := createRandomTransfer(t)
	result, err := testQueries.GetTransfer(context.Background(), Transfer.ID)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.Equal(t, Transfer.ID, result.ID)
	require.Equal(t, Transfer.FromAccountID, result.FromAccountID)
	require.Equal(t, Transfer.ToAccountID, result.ToAccountID)
	require.Equal(t, Transfer.Amount, result.Amount)
	// There might be a short delay from creation of the random Transfer
	// to its storage in the DB and we don't this to fail the test.
	require.WithinDuration(t, Transfer.CreatedAt, result.CreatedAt, time.Second)
}

func TestUpdateTransfer(t *testing.T) {
	Transfer := createRandomTransfer(t)

	arg := sqlc.UpdateTransferParams{
		ID:     Transfer.ID,
		Amount: util.RandomMoney(),
	}

	result, err := testQueries.UpdateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.Equal(t, Transfer.ID, result.ID)
	require.Equal(t, Transfer.FromAccountID, result.FromAccountID)
	require.Equal(t, Transfer.ToAccountID, result.ToAccountID)
	// The amount should change to the new value.
	require.Equal(t, arg.Amount, result.Amount)
	// There might be a short delay from creation of the random Transfer
	// to its storage in the DB and we don't this to fail the test.
	require.WithinDuration(t, Transfer.CreatedAt, result.CreatedAt, time.Second)
}

func TestDeleteTransfer(t *testing.T) {
	Transfer := createRandomTransfer(t)
	err := testQueries.DeleteTransfer(context.Background(), Transfer.ID)
	require.NoError(t, err)

	// Validation that the Transfer was truly deleted
	valTransfer, err := testQueries.GetTransfer(context.Background(), Transfer.ID)
	// We want an error to occur here.
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, valTransfer)
}

func TestListTransfers(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomTransfer(t)
	}

	arg := sqlc.ListTransfersParams{
		Limit:  5,
		Offset: 5, // Skips 5 matches before returning values.
	}
	results, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)

	for _, result := range results {
		require.NotEmpty(t, result)
	}
}
