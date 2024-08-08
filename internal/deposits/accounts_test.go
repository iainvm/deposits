package deposits_test

import (
	"errors"
	"log/slog"
	"testing"

	"github.com/google/uuid"
	"github.com/iainvm/deposits/common/pointers"
	"github.com/iainvm/deposits/internal/deposits"
	"github.com/stretchr/testify/require"
)

func TestNewTotalAllocatedAmount(t *testing.T) {
	testCases := []struct {
		description   string
		input         int64
		expectedError error
		expectedValue *deposits.TotalAllocatedAmount
	}{
		{
			description:   "passes for 10",
			input:         10,
			expectedError: nil,
			expectedValue: pointers.New(deposits.TotalAllocatedAmount(10)),
		},
		{
			description:   "passes for 0",
			input:         0,
			expectedError: nil,
			expectedValue: pointers.New(deposits.TotalAllocatedAmount(0)),
		},
		{
			description:   "fails for -1",
			input:         -1,
			expectedError: deposits.ErrNegativeAmount,
			expectedValue: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {

			actualValue, actualError := deposits.NewTotalAllocatedAmount(testCase.input)

			if actualError != nil {
				require.ErrorIs(t, actualError, testCase.expectedError)
				return
			}

			require.Equal(t, actualValue, *testCase.expectedValue)
			require.Equal(t, actualValue.Int64(), testCase.expectedValue.Int64())

		})
	}
}

func TestParseAccountId(t *testing.T) {
	staticUUID := uuid.NewString()
	testCases := []struct {
		description   string
		input         string
		expectedError error
		expectedValue *deposits.AccountId
	}{
		{
			description:   "passes for uuid",
			input:         staticUUID,
			expectedError: nil,
			expectedValue: pointers.New(deposits.AccountId(staticUUID)),
		},
		{
			description:   "fail for random string",
			input:         "abcdef",
			expectedError: errors.New("invalid UUID length: 6"),
			expectedValue: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {

			actualValue, actualError := deposits.ParseAccountId(testCase.input)

			if actualError != nil {
				require.ErrorContains(t, actualError, testCase.expectedError.Error())
				return
			}

			require.Equal(t, actualValue, *testCase.expectedValue)
			require.Equal(t, actualValue.String(), testCase.expectedValue.String())

		})
	}
}

func TestNewAccount(t *testing.T) {
	t.Run("successful data", func(t *testing.T) {
		account, err := deposits.NewAccount(
			deposits.WrapperTypeISA,
			123456,
		)

		require.NoError(t, err)
		require.Equal(t, &deposits.Account{
			Id:                   account.Id,
			WrapperType:          deposits.WrapperTypeISA,
			NominalAmount:        123456,
			TotalAllocatedAmount: 0,
		}, account)
	})

	t.Run("invalid wrapper", func(t *testing.T) {
		_, err := deposits.NewAccount(
			0,
			123456,
		)

		require.ErrorIs(t, err, deposits.ErrInvalidWrapperType)
	})
}

func TestParseAccount(t *testing.T) {

	t.Run("successful data", func(t *testing.T) {
		account, err := deposits.ParseAccount(uuid.NewString(), 1, 10, 0)

		require.NoError(t, err)
		require.Equal(t, &deposits.Account{
			Id:                   account.Id,
			WrapperType:          1,
			NominalAmount:        10,
			TotalAllocatedAmount: 0,
		}, account)
	})

	t.Run("invalid id", func(t *testing.T) {
		_, err := deposits.ParseAccount("string", 1, 10, 0)

		require.ErrorContains(t, err, "invalid UUID length")
	})

	t.Run("invalid type", func(t *testing.T) {
		_, err := deposits.ParseAccount(uuid.NewString(), 0, 10, 0)

		require.ErrorIs(t, err, deposits.ErrInvalidWrapperType)
	})

	t.Run("invalid nominal amount", func(t *testing.T) {
		_, err := deposits.ParseAccount(uuid.NewString(), 1, -1, 0)

		require.ErrorIs(t, err, deposits.ErrNominalAmountNegative)
	})
}

func TestAddReceipt(t *testing.T) {

	accountUUID := uuid.NewString()
	account, err := deposits.ParseAccount(accountUUID, deposits.WrapperTypeSIPP.Int(), 100, 0)
	require.NoError(t, err)

	receiptUUID := uuid.NewString()
	receipt, err := deposits.ParseReceipt(receiptUUID, 50)
	require.NoError(t, err)

	err = account.AddReceipt(receipt)
	require.NoError(t, err)

	receiptUUID = uuid.NewString()
	receipt, err = deposits.ParseReceipt(receiptUUID, 50)
	require.NoError(t, err)

	err = account.AddReceipt(receipt)
	require.NoError(t, err)

	receiptUUID = uuid.NewString()
	receipt, err = deposits.ParseReceipt(receiptUUID, 50)
	require.NoError(t, err)

	err = account.AddReceipt(receipt)
	slog.Info("asdf", "account", account)
	require.ErrorIs(t, err, deposits.ErrNominalExceeded)
}
