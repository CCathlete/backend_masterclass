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
func createRandomTransfer(t *testing.T, fromID, toID int64) sqlc.Transfer {
	arg := sqlc.CreateTransferParams{
		FromAccountID: fromID,
		ToAccountID:   toID,
	}

	Transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, Transfer)

	require.Equal(t, arg.Owner, Transfer.Owner)
	require.Equal(t, arg.Balance, Transfer.Balance)
	require.Equal(t, arg.Currency, Transfer.Currency)

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
	require.Equal(t, Transfer.Owner, result.Owner)
	require.Equal(t, Transfer.Balance, result.Balance)
	require.Equal(t, Transfer.Currency, result.Currency)
	// There might be a short delay from creation of the random Transfer
	// to its storage in the DB and we don't this to fail the test.
	require.WithinDuration(t, Transfer.CreatedAt, result.CreatedAt, time.Second)
}

func TestUpdateTransfer(t *testing.T) {
	Transfer := createRandomTransfer(t)

	arg := sqlc.UpdateTransferParams{
		ID:      Transfer.ID,
		Balance: util.RandomMoney(),
	}

	result, err := testQueries.UpdateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.Equal(t, Transfer.ID, result.ID)
	require.Equal(t, Transfer.Owner, result.Owner)
	// The balance should change to the new value.
	require.Equal(t, arg.Balance, result.Balance)
	require.Equal(t, Transfer.Currency, result.Currency)
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
