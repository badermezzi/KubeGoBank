package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/badermezzi/KubeGoBank/util"

	"github.com/stretchr/testify/require"
)

// Helper function to create a random account for testing
func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t) // Create a random user
	arg := CreateAccountParams{
		Owner:    user.Username, // Use user's username as account owner
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Currency, account.Currency)
	require.Equal(t, arg.Balance, account.Balance)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

// Helper function to create a random user for testing

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t) // Now TestCreateAccount just calls the helper
}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t) // Create an account to retrieve

	account2, err := testQueries.GetAccount(context.Background(), account1.ID) // Retrieve it using GetAccount
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID) // Ensure retrieved account matches the created one
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Currency, account2.Currency)
	require.Equal(t, account1.Balance, account2.Balance)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second) // Check CreatedAt within a second
}

func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t) // Create an account to update

	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: util.RandomMoney(), // Update balance with a new random amount
	}

	account2, err := testQueries.UpdateAccount(context.Background(), arg) // Update the account
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)                                     // IDs should be the same
	require.Equal(t, account1.Owner, account2.Owner)                               // Owner should remain the same
	require.Equal(t, account1.Currency, account2.Currency)                         // Currency should remain the same
	require.Equal(t, arg.Balance, account2.Balance)                                // Balance should be updated to the new random amount
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second) // CreatedAt should be roughly the same
}

func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t) // Create an account to delete

	err := testQueries.DeleteAccount(context.Background(), account1.ID) // Delete the account
	require.NoError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID) // Try to get the deleted account
	require.Error(t, err)                                                      // Expect an error (because it should be deleted)
	require.EqualError(t, err, sql.ErrNoRows.Error())                          // Specifically expect "sql: no rows in result" error
	require.Empty(t, account2)                                                 // Account should be empty (nil)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ { // Create 10 random accounts for listing
		createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  5, // Limit to 5 accounts per page
		Offset: 5, // Start from the 6th account (for the second "page")
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg) // List accounts with limit and offset
	require.NoError(t, err)
	require.Len(t, accounts, 5) // Expect to get 5 accounts back because of the limit

	for _, account := range accounts {
		require.NotEmpty(t, account) // Ensure each account in the list is not empty
	}
}
