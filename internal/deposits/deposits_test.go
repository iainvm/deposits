package deposits_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/iainvm/deposits/internal/deposits"
	"github.com/stretchr/testify/require"
)

func TestParseDeposits(t *testing.T) {
	_, err := deposits.ParseDeposit(uuid.NewString())
	require.NoError(t, err)
}

func TestNewDeposits(t *testing.T) {
	_, err := deposits.New()
	require.NoError(t, err)
}
