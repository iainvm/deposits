package investor_test

import (
	"testing"

	"github.com/iainvm/deposits/internal/investor"
	"github.com/stretchr/testify/require"
)

func TestName(t *testing.T) {

	testCases := []struct {
		desc         string
		name         string
		expectedErr  error
		expectedName investor.Name
	}{
		{
			desc:         "Successful Name",
			name:         "Iain Majer",
			expectedErr:  nil,
			expectedName: investor.Name("Iain Majer"),
		},
		{
			desc:         "Blank Name",
			name:         "",
			expectedErr:  investor.ErrBlankName,
			expectedName: investor.Name(""),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.desc, func(t *testing.T) {
			name, err := investor.NewName(testCase.name)

			if testCase.expectedErr != nil {
				require.ErrorIs(t, err, testCase.expectedErr)
				return
			}

			require.Equal(t, testCase.expectedName, name)

		})
	}
}
