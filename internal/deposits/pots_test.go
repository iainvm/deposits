package deposits_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/iainvm/deposits/internal/deposits"
	"github.com/stretchr/testify/require"
)

func TestNewPot(t *testing.T) {
	pot, err := deposits.NewPot("abcdefg")
	require.NoError(t, err)
	require.Equal(t, &deposits.Pot{
		Id:   pot.Id,
		Name: "abcdefg",
	}, pot)
}

func TestParsePot(t *testing.T) {
	id := uuid.NewString()
	pot, err := deposits.ParsePot(id, "abcdefg")
	require.NoError(t, err)
	require.Equal(t, &deposits.Pot{
		Id:   pot.Id,
		Name: "abcdefg",
	}, pot)
}

func TestAddAccount(t *testing.T) {

	id := uuid.NewString()
	pot, err := deposits.ParsePot(id, "Pot A")
	require.NoError(t, err)

	account, err := deposits.ParseAccount(uuid.NewString(), 1, 100, 0)
	require.NoError(t, err)

	err = pot.AddAccount(account)
	require.NoError(t, err)

	err = pot.AddAccount(account)
	require.ErrorIs(t, err, deposits.ErrWrapperTypeExistsInPot)
}
