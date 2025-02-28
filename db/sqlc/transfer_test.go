package db

import (
	"context"
	"testing"
	"time"

	"github.com/badermezzi/KubeGoBank/util"

	"github.com/stretchr/testify/require"
)

// Helper function to create random accounts for transfer testing
func createRandomTransferAccounts(t *testing.T) (account1 Account, account2 Account) {
	account1 = createRandomAccount(t)
	account2 = createRandomAccount(t)
	return
}

// Helper function to create a random transfer for testing
func createRandomTransfer(t *testing.T) Transfer {
	account1, account2 := createRandomTransferAccounts(t) // Create two random accounts
	amount := util.RandomMoney()

	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        amount,
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	createRandomTransfer(t)
}

func TestGetTransfer(t *testing.T) {
	transfer1 := createRandomTransfer(t)

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, transfer1.Amount, transfer2.Amount)
	require.WithinDuration(t, transfer1.CreatedAt, transfer2.CreatedAt, time.Second)
}

func TestListTransfers(t *testing.T) {
	n := 10
	account := createRandomAccount(t) // Create a single account to filter transfers for
	for i := 0; i < n; i++ {
		createRandomTransfer(t) // Create some random transfers (not necessarily involving 'account')
	}
	for i := 0; i < 5; i++ { // Create 5 transfers where 'account' is the sender
		createRandomTransferWithAccount(t, account, createRandomAccount(t))
	}
	for i := 0; i < 5; i++ { // Create 5 transfers where 'account' is the receiver
		createRandomTransferWithAccount(t, createRandomAccount(t), account)
	}

	arg := ListTransfersParams{
		Limit:         5,
		Offset:        5,
		FromAccountID: account.ID, // Filter by from_account_id
		ToAccountID:   account.ID, // Also filter by to_account_id (in OR condition of query)
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5) // Expecting 5 transfers based on limit

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.True(t, transfer.FromAccountID == account.ID || transfer.ToAccountID == account.ID) // Verify transfer involves the specified account
	}
}

// Helper function to create a random transfer with specific sender and receiver accounts
func createRandomTransferWithAccount(t *testing.T, fromAccount Account, toAccount Account) Transfer {
	amount := util.RandomMoney()

	arg := CreateTransferParams{
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
		Amount:        amount,
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}
