package db

import (
	"context"
	"testing"
	"time"

	"github.com/ebaudet/simplebank/utils"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) (User, CreateUserParams) {
	hash, err := utils.HashPassword(utils.RandomPassword(6, 25))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:       utils.RandomOwner(),
		HashedPassword: hash,
		FullName:       utils.RandomFullName(),
		Email:          utils.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)

	return user, arg
}

func TestCreateUser(t *testing.T) {
	user, arg := createRandomUser(t)

	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.NotZero(t, user.CreatedAt)
}

func TestGetUser(t *testing.T) {
	user1, _ := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)

	require.Equal(t, user1, user2)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}
