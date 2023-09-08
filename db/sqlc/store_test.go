package db

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(dbConnection)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">> before: ", account1.Balance, account2.Balance)
	n := 10
	amount := int64(5)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				account1.ID,
				account2.ID,
				amount,
			})
			errs <- err
			results <- result
		}()
	}
	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, transfer.FromAccountID, account1.ID)
		require.Equal(t, transfer.ToAccountID, account2.ID)
		require.Equal(t, transfer.Amount, amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		fromEntry := result.FromEntry
		require.Equal(t, fromEntry.Amount, -amount)
		require.Equal(t, fromEntry.AccountID, account1.ID)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.Equal(t, toEntry.Amount, amount)
		require.Equal(t, toEntry.AccountID, account2.ID)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check accounts balance

		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, fromAccount.ID, account1.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, toAccount.ID, account2.ID)

		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance

		require.Equal(t, diff2, diff1)
		require.True(t, diff2 > 0)
		require.True(t, diff2%amount == 0)

		k := int(diff2 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount1)

	updatedAccount2, err2 := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err2)
	require.NotEmpty(t, updatedAccount2)

	fmt.Println(">> after: ", updatedAccount1.Balance, updatedAccount2.Balance)

	require.Equal(t, updatedAccount1.Balance, account1.Balance-int64(n)*amount)
	require.Equal(t, updatedAccount2.Balance, account2.Balance+int64(n)*amount)
}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(dbConnection)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">> before: ", account1.Balance, account2.Balance)
	n := 10
	amount := int64(5)

	errs := make(chan error)

	for i := 0; i < n; i++ {
		i := i
		go func() {
			fromAccountId := account1.ID
			toAccountId := account2.ID

			if i%2 == 0 {
				fromAccountId = account2.ID
				toAccountId = account1.ID
			}
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				fromAccountId,
				toAccountId,
				amount,
			})
			errs <- err
		}()
	}
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount1)

	updatedAccount2, err2 := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err2)
	require.NotEmpty(t, updatedAccount2)

	fmt.Println(">> after: ", updatedAccount1.Balance, updatedAccount2.Balance)

	require.Equal(t, updatedAccount1.Balance, account1.Balance)
	require.Equal(t, updatedAccount2.Balance, account2.Balance)
}
