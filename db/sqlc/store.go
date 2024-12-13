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
func CommitOrRollback(transaction *sql.Tx, err error) error {
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

	log.Println("Commited transaction.")
	return nil
}

// Executes a function within a database transaction.
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("execTx: %w", err)
	}
	defer CommitOrRollback(tx, err)

	// We create a *Qeries object with a transaction instead of regular db.
	q := New(tx)
	err = fn(q)
	if err != nil {
		return fmt.Errorf("execTx: %w", err)
	}

	return nil
}
