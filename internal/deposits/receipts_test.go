package deposits_test

import (
	"testing"

	"github.com/iainvm/deposits/internal/deposits"
	"github.com/stretchr/testify/require"
)

func TestNewReceipt(t *testing.T) {

	receipt, err := deposits.NewReceipt(100)
	require.NoError(t, err)
	require.Equal(t, &deposits.Receipt{
		Id:              receipt.Id,
		AllocatedAmount: 100,
	}, receipt)
}
