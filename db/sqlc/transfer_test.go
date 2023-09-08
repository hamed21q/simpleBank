package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"simpleBank/util"
	"testing"
	"time"
)

func createRandomTransferForAccounts(t *testing.T, fromAccountId, toAccountId int64) Transfer {
	arg := CreateTransferParams{
		FromAccountID: fromAccountId,
		ToAccountID:   toAccountId,
		Amount:        util.RandomMoney(),
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

func createRandomTransfer(t *testing.T) Transfer {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	return createRandomTransferForAccounts(t, account1.ID, account2.ID)
}

func TestCreateTransfer(t *testing.T) {
	createRandomTransfer(t)
}

func TestGetTransfer(t *testing.T) {
	transfer1 := createRandomTransfer(t)
	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, transfer2)
	require.Equal(t, transfer2.ToAccountID, transfer1.ToAccountID)
	require.Equal(t, transfer2.FromAccountID, transfer1.FromAccountID)
	require.Equal(t, transfer2.ID, transfer1.ID)
	require.Equal(t, transfer2.Amount, transfer1.Amount)
	require.WithinDuration(t, transfer2.CreatedAt, transfer1.CreatedAt, time.Second)
}

func TestListTransfersToAccount(t *testing.T) {
	account := createRandomAccount(t)
	for i := 0; i < 10; i++ {
		fromAccount := createRandomAccount(t)
		createRandomTransferForAccounts(t, fromAccount.ID, account.ID)
	}
	arg := ListTransfersToAccountParams{
		ToAccountID: account.ID,
		Limit:       5,
		Offset:      5,
	}
	transfers, err := testQueries.ListTransfersToAccount(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)
	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}
}

func TestListTransfersFromAccount(t *testing.T) {
	account := createRandomAccount(t)
	for i := 0; i < 10; i++ {
		toAccount := createRandomAccount(t)
		createRandomTransferForAccounts(t, account.ID, toAccount.ID)
	}
	arg := ListTransfersFromAccountParams{
		FromAccountID: account.ID,
		Limit:         5,
		Offset:        5,
	}
	transfers, err := testQueries.ListTransfersFromAccount(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)
	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}
}
