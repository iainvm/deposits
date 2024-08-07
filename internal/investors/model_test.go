package investors_test

import (
	"testing"

	"github.com/iainvm/deposits/internal/investors"
	"github.com/stretchr/testify/require"
)

func TestName(t *testing.T) {

	testCases := []struct {
		desc         string
		name         string
		expectedErr  error
		expectedName investors.Name
	}{
		{
			desc:         "Successful Name",
			name:         "Iain Majer",
			expectedErr:  nil,
			expectedName: investors.Name("Iain Majer"),
		},
		{
			desc:         "Blank Name",
			name:         "",
			expectedErr:  investors.ErrBlankName,
			expectedName: investors.Name(""),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.desc, func(t *testing.T) {
			name, err := investors.NewName(testCase.name)

			if testCase.expectedErr != nil {
				require.ErrorIs(t, err, testCase.expectedErr)
				return
			}

			require.Equal(t, testCase.expectedName, name)

		})
	}
}
