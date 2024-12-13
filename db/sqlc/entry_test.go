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
func createRandomEntry(t *testing.T) sqlc.Entry {
	arg := sqlc.CreateEntryParams{
		AccountID: createRandomAccount(t).ID,
		Amount:    util.RandomMoney(),
	}

	Entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, Entry)

	require.Equal(t, arg.AccountID, Entry.AccountID)
	require.Equal(t, arg.Amount, Entry.Amount)

	require.NotZero(t, Entry.ID)
	require.NotZero(t, Entry.CreatedAt)

	return Entry
}

func TestCreateEntry(t *testing.T) {
	createRandomEntry(t)
}

func TestGetEntry(t *testing.T) {
	Entry := createRandomEntry(t)
	result, err := testQueries.GetEntry(context.Background(), Entry.ID)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.Equal(t, Entry.ID, result.ID)
	require.Equal(t, Entry.AccountID, result.AccountID)
	require.Equal(t, Entry.Amount, result.Amount)
	// There might be a short delay from creation of the random Entry
	// to its storage in the DB and we don't this to fail the test.
	require.WithinDuration(t, Entry.CreatedAt, result.CreatedAt, time.Second)
}

func TestUpdateEntry(t *testing.T) {
	Entry := createRandomEntry(t)

	arg := sqlc.UpdateEntryParams{
		ID:     Entry.ID,
		Amount: util.RandomMoney(),
	}

	result, err := testQueries.UpdateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.Equal(t, Entry.ID, result.ID)
	require.Equal(t, Entry.AccountID, result.AccountID)
	// The balance should change to the new value.
	require.Equal(t, arg.Amount, result.Amount)
	// There might be a short delay from creation of the random Entry
	// to its storage in the DB and we don't this to fail the test.
	require.WithinDuration(t, Entry.CreatedAt, result.CreatedAt, time.Second)
}

func TestUpdateEntryByAccount(t *testing.T) {
	account := createRandomAccount(t)

	createArg := sqlc.CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), createArg)
	require.NoError(t, err)

	updateArg := sqlc.UpdateEntryByAccountParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	}

	result, err := testQueries.UpdateEntryByAccount(context.Background(), updateArg)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.Equal(t, entry.ID, result.ID)
	require.Equal(t, entry.AccountID, result.AccountID)
	// The balance should change to the new value.
	require.Equal(t, updateArg.Amount, result.Amount)
	// There might be a short delay from creation of the random entry
	// to its storage in the DB and we don't this to fail the test.
	require.WithinDuration(t, entry.CreatedAt, result.CreatedAt, time.Second)
}

func TestDeleteEntry(t *testing.T) {
	Entry := createRandomEntry(t)
	err := testQueries.DeleteEntry(context.Background(), Entry.ID)
	require.NoError(t, err)

	// Validation that the Entry was truly deleted
	valEntry, err := testQueries.GetEntry(context.Background(), Entry.ID)
	// We want an error to occur here.
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, valEntry)
}

func TestListEntries(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomEntry(t)
	}

	arg := sqlc.ListEntriesParams{
		Limit:  5,
		Offset: 5, // Skips 5 matches before returning values.
	}
	results, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)

	for _, result := range results {
		require.NotEmpty(t, result)
	}
}

func TestGetAccountEntries(t *testing.T) {
	account := createRandomAccount(t)

	for i := 0; i < 10; i++ {
		createArg := sqlc.CreateEntryParams{
			AccountID: account.ID,
			Amount:    util.RandomMoney(),
		}

		_, err := testQueries.CreateEntry(context.Background(), createArg)
		require.NoError(t, err)
	}

	results, err := testQueries.GetAccountEntries(context.Background(),
		account.ID)
	require.NoError(t, err)

	for _, result := range results {
		require.NotEmpty(t, result)
	}
}
