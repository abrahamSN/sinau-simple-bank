package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store privides all functions to execute db queries and transactions
type Store struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new Store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction
// The provided function fn receives a Store argument which
// contains all the store functions. Once the function fn
// returns an error, the transaction is rolled back. Otherwise,
// the transaction is committed.
// This function is used to execute multiple queries in a
// single transaction.
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	// txOpt is a transaction option that allows us to set the isolation level of the transaction
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// We create a new Queries object with the transaction object
	q := New(t x)
	// We call the function fn with the new Queries object
	err = fn(q)
	if err != nil {
		// If the function returns an error, we rollback the transaction
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction error: %v, rollback error: %v", err, rbErr)
		}
		return err
	}
	// If the function returns no error, we commit the transaction
	return tx.Commit()
}
