package sqlc_test

import (
	"backend-masterclass/db/sqlc"
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	existed := make(map[int]bool)
	store := sqlc.NewStore(testDB)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">> before:", account1.Balance, account2.Balance)

	// Run a concurrent transfer transaction.
	// The concurrency here simulates a real life scenario of
	// a possible parallel access to the same resource.
	// Using an unbuffered channel will lock all goroutins except of the
	// first one that got to the results/ error channels.
	// Access to the DB would happen in parallel  but reading the result
	// would happen serially.
	// Why not a buffered channel?
	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan sqlc.TransferTxResult)
	var ctx context.Context

	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx %d", i+1)
		go func(txName string) {
			ctx = context.WithValue(context.Background(),
				sqlc.TxKey, txName)
			result, err := store.TransferTx(ctx,
				sqlc.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        amount,
				})

			errs <- err
			results <- result
			if err != nil {
				fmt.Println(txName, "sent its error to the channel.")
			} else {
				fmt.Println(txName, "sent its result to the channel.")
			}
		}(txName)
	}

	// Check results.
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// Check Transfer.
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, transfer.Amount, amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		// Validate presence in DB.
		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// Check entries.
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -(amount), fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		// Validate presence in DB.
		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, +(amount), toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		// Validate presence in DB.
		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// Check accounts.
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)
		//
		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		// Validate presence in DB.
		_, err = store.GetAccount(context.Background(), fromAccount.ID)
		require.NoError(t, err)
		//
		_, err = store.GetAccount(context.Background(), toAccount.ID)
		require.NoError(t, err)

		// Check accounts' balance.
		fmt.Printf(">> Tx: %d %d\n", fromAccount.Balance,
			toAccount.Balance)
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		// We perform n transactions, each time substracting/adding
		// by amount.
		require.True(t, diff1%amount == 0)

		// We want to make sure that the balances actually changes on
		// each read iteration.
		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		existed[k] = true // The diff must change with each transaction.
	}

	// Checking the final result after all transactions are done.
	// We want to check each of the accounts in this case.
	updatedAccount1, err := testQueries.GetAccount(context.Background(),
		account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(),
		account2.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)
	require.Equal(t, amount*int64(n),
		account1.Balance-updatedAccount1.Balance)
	require.Equal(t, amount*int64(n),
		updatedAccount2.Balance-account2.Balance)
}

// Here all we care about is if there are any deadlocks.
func TestTransferTxDeadlock(t *testing.T) {
	store := sqlc.NewStore(testDB)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	n := 10
	amount := int64(10)

	errs := make(chan error)
	var ctx context.Context

	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx %d", i+1)

		go func(txName string, i int) {
			var err error

			ctx = context.WithValue(context.Background(),
				sqlc.TxKey, txName)

			fromAccount := account1
			toAccount := account2

			if i%2 != 0 {
				fromAccount = account2
				toAccount = account1
			}

			_, err = store.TransferTx(ctx,
				sqlc.TransferTxParams{
					FromAccountID: fromAccount.ID,
					ToAccountID:   toAccount.ID,
					Amount:        amount,
				})

			errs <- err
			if err != nil {
				fmt.Println(txName, "sent its error to the channel.")
			}
		}(txName, i)
	}

	// Check errors.
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	// Checking the final result after all transactions are done.
	updatedAccount1, err := testQueries.GetAccount(context.Background(),
		account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(),
		account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)
}
