package sqlc

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

// Store provides all functions to execute db queries and transactions.
type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// This will be the last thing executed (first defer).
// It gets the last error value of the function and the transaction
// and commits or rolls back according to the error.
func CommitOrRollback(ctx context.Context, transaction *sql.Tx, err error) error {
	txName := ctx.Value(TxKey)

	switch err {
	case nil:
		// Committing the transaction to the DB.
		// If we have multiple statements in the same transaction
		// we need to commit after all are executed successfully!
		err := transaction.Commit()
		if err != nil {
			log.Println(
				"There was a problem with committing the transaction.")
			return fmt.Errorf("commit err: %w", err)
		}
	default:
		err = transaction.Rollback()
		if err != nil {
			log.Println(
				"There was a problem with rolling the transaction back.")
			return fmt.Errorf("rollback err: %w", err)
		}
		return fmt.Errorf("rollback performed: %w", err)
	}

	fmt.Println(txName, "Committed.")
	return nil
}

// Executes a function within a database transaction.
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("execTx: %w", err)
	}
	defer CommitOrRollback(ctx, tx, err)

	// We create a *Qeries object with a transaction instead of regular db.
	q := New(tx)
	err = fn(q)
	if err != nil {
		return fmt.Errorf("execTx: %w", err)
	}

	return nil
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

var TxKey *string

// Performs a mony transfer from one account to another.
// It creates a transfer record, add account entries, and
// updates accounts' balance within a single DB ransaction.
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	txName := ctx.Value(TxKey)

	err := store.execTx(ctx,
		func(q *Queries) error {
			var err error

			// The params types have the same fields so we can do a simple
			// conversion.
			result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams(arg))
			if err != nil {
				return fmt.Errorf("%s: TransferTx: %w", txName, err)
			}

			result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
				AccountID: arg.FromAccountID,
				Amount:    -(arg.Amount), // Currency exits this account.
			})
			if err != nil {
				return fmt.Errorf("%s: TransferTx: %w", txName, err)
			}

			result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
				AccountID: arg.ToAccountID,
				Amount:    +(arg.Amount), // Currency enters this account.
			})
			if err != nil {
				return fmt.Errorf("%s: TransferTx: %w", txName, err)
			}

			// The updating of accounts must be made in a consistent way
			// regarding the account ids. Meaning, that we can't have an
			// account being updated first when it's a from account but
			// updated second when its to account - CREATES DEADLOCK.
			if arg.FromAccountID < arg.ToAccountID {

				result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID,
					arg.ToAccountID, -arg.Amount, +arg.Amount)
				if err != nil {
					return fmt.Errorf("%s: TransferTx: %w", txName, err)
				}

			} else {

				result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID,
					arg.FromAccountID, +arg.Amount, -arg.Amount)
				if err != nil {
					return fmt.Errorf("%s: TransferTx: %w", txName, err)
				}

			}

			return nil
		})
	if err != nil {
		return TransferTxResult{}, fmt.Errorf("%s: TransferTx: %w", txName, err)
	}

	return result, nil
}

func addMoney(
	ctx context.Context,
	q *Queries,
	aID1,
	aID2,
	amount1,
	amount2 int64,
) (account1 Account, account2 Account, err error) {

	account1, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
		Amount: amount1,
		ID:     aID1,
	})
	if err != nil {
		return // because of named return values, returns empty accounts and err.
	}

	account2, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
		Amount: amount2,
		ID:     aID2,
	})
	if err != nil {
		return
	}

	return
}
