package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier // Embed Querier interface
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

// SQLStore provides all functions to execute SQL queries and transactions
type SQLStore struct {
	*Queries         // Embed Queries to access individual query functions
	db       *sql.DB // Hold the database connection
}

// NewStore creates a new store
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		Queries: New(db), // Initialize embedded Queries with the db connection
		db:      db,
	}
}

// executeTx executes a function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil) // Start a new transaction
	if err != nil {
		return err
	}

	q := New(tx) // Create new Queries object using the transaction
	err = fn(q)  // Execute the provided function within the transaction
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr) // Return error if rollback also fails
		}
		return err // Return original error if rollback succeeds
	}

	return tx.Commit() // Commit transaction if function execution is successful
}

// TransferTxParams contains the parameters for the transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
	FromAccount Account  `json:"from_account"` // Add FromAccount to result
	ToAccount   Account  `json:"to_account"`   // Add ToAccount to result
}

// TransferTx performs a money transfer from one account to another.
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams(arg))
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount, // Debit from account
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount, // Credit to account
		})
		if err != nil {
			return err
		}

		// Update account balances using addMoney for consistent locking.
		result.FromAccount, result.ToAccount, err = store.addMoney(
			ctx,
			q,
			arg.FromAccountID,
			-arg.Amount,
			arg.ToAccountID,
			arg.Amount,
		)

		return err // Return any error from addMoney or Create operations
	})

	return result, err
}

// addMoney adds or subtracts money from an account.  It ensures consistent locking
// by always updating the account with the lower ID first.
func (store *SQLStore) addMoney(ctx context.Context, q *Queries, accountID1 int64, amount1 int64, accountID2 int64, amount2 int64) (account1 Account, account2 Account, err error) {
	if accountID1 < accountID2 {
		account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     accountID1,
			Amount: amount1,
		})
		if err != nil {
			return // account1 and account2 will be their zero values, err is populated
		}
		account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     accountID2,
			Amount: amount2,
		})
		if err != nil {
			return //account1 is properly initialized, account2 is zero value
		}
	} else {
		account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     accountID2,
			Amount: amount2,
		})
		if err != nil {
			return // account1 and account2 will be their zero values, err is populated
		}

		account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     accountID1,
			Amount: amount1,
		})
		if err != nil {
			return //account2 is properly initialized, account1 is zero value
		}
	}
	return // Both accounts, and possibly an error, are returned.
}
