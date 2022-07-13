package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/ebaudet/simplebank/utils"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T, from Account, to Account) (Transfer, CreateTransferParams) {
	arg := CreateTransferParams{
		FromAccountID: from.ID,
		ToAccountID:   to.ID,
		Amount:        utils.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)

	return transfer, arg
}

func TestCreateTransfer(t *testing.T) {
	from, _ := createRandomAccount(t)
	to, _ := createRandomAccount(t)
	transfer, arg := createRandomTransfer(t, from, to)

	require.NotEmpty(t, transfer)

	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)
}

func TestGetTransfer(t *testing.T) {
	from, _ := createRandomAccount(t)
	to, _ := createRandomAccount(t)
	transfer1, _ := createRandomTransfer(t, from, to)
	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)

	require.Equal(t, transfer1, transfer2)
	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, transfer1.Amount, transfer2.Amount)
	require.WithinDuration(t, transfer1.CreatedAt, transfer2.CreatedAt, time.Second)
}

func TestUpdateTransfer(t *testing.T) {
	from, _ := createRandomAccount(t)
	to, _ := createRandomAccount(t)
	transfer, _ := createRandomTransfer(t, from, to)

	arg := UpdateTransferParams{
		ID:     transfer.ID,
		Amount: utils.RandomMoney(),
	}

	transfer2, err := testQueries.UpdateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer.ID, transfer2.ID)
	require.Equal(t, transfer.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, arg.Amount, transfer2.Amount)
}

func TestDeleteTransfer(t *testing.T) {
	from, _ := createRandomAccount(t)
	to, _ := createRandomAccount(t)
	transfer1, _ := createRandomTransfer(t, from, to)

	err := testQueries.DeleteTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, transfer2)
}

func TestListTransfers(t *testing.T) {
	from, _ := createRandomAccount(t)
	to, _ := createRandomAccount(t)
	for i := 0; i < 10; i++ {
		_, _ = createRandomTransfer(t, from, to)
	}

	arg := ListTransfersParams{
		FromAccountID: from.ID,
		ToAccountID:   to.ID,
		Limit:         5,
		Offset:        5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)
		require.True(t, transfer.FromAccountID == from.ID || transfer.ToAccountID == to.ID)
	}
}
