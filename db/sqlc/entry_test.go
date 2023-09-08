package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"simpleBank/util"
	"testing"
	"time"
)

func createRandomEntryForAccount(t *testing.T, accountId int64) Entry {
	arg := CreateEntryParams{
		accountId,
		util.RandomMoney(),
	}
	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)
	return entry
}

func createRandomEntry(t *testing.T) Entry {
	account := createRandomAccount(t)
	return createRandomEntryForAccount(t, account.ID)
}

func TestCreateEntry(t *testing.T) {
	createRandomAccount(t)
}

func TestGetEntry(t *testing.T) {
	entry1 := createRandomEntry(t)
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, entry2)
	require.Equal(t, entry2.AccountID, entry1.AccountID)
	require.Equal(t, entry2.ID, entry2.ID)
	require.Equal(t, entry2.Amount, entry1.Amount)
	require.WithinDuration(t, entry2.CreatedAt, entry1.CreatedAt, time.Second)
}

func TestListOfEntries(t *testing.T) {
	account := createRandomAccount(t)
	for i := 0; i < 10; i++ {
		createRandomEntryForAccount(t, account.ID)
	}
	arg := ListEntriesForAccountParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    5,
	}
	entries, err := testQueries.ListEntriesForAccount(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)
	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}
