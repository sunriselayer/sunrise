package erasurecoding

import (
	"testing"

	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
	"github.com/stretchr/testify/require"
)

func TestPadTo32(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []byte
	}{
		{
			name:     "already padded",
			input:    make([]byte, 32),
			expected: make([]byte, 32),
		},
		{
			name:     "needs padding",
			input:    []byte{1, 2, 3},
			expected: append([]byte{1, 2, 3}, make([]byte, 29)...),
		},
		{
			name:     "empty input",
			input:    []byte{},
			expected: []byte{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := padTo32(tt.input)
			require.Equal(t, tt.expected, result)
			if len(result) > 0 {
				require.Equal(t, 0, len(result)%32)
			}
		})
	}
}

func TestConvertToElements(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected int // expected number of elements
	}{
		{
			name:     "32 bytes",
			input:    make([]byte, 32),
			expected: 1,
		},
		{
			name:     "64 bytes",
			input:    make([]byte, 64),
			expected: 2,
		},
		{
			name:     "31 bytes (needs padding)",
			input:    make([]byte, 31),
			expected: 1,
		},
		{
			name:     "empty input",
			input:    []byte{},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertToElements(tt.input)
			require.Equal(t, tt.expected, len(result))
		})
	}
}

func TestSplitElementsIntoRows(t *testing.T) {
	tests := []struct {
		name          string
		elements      []fr.Element
		expectedRows  int
		checkRowSizes bool
	}{
		{
			name:          "single row (small)",
			elements:      make([]fr.Element, 100),
			expectedRows:  1,
			checkRowSizes: true,
		},
		{
			name:          "multiple rows (large)",
			elements:      make([]fr.Element, 5000),
			expectedRows:  2,
			checkRowSizes: true,
		},
		{
			name:          "empty elements",
			elements:      []fr.Element{},
			expectedRows:  1,
			checkRowSizes: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rows := SplitElementsIntoRows(tt.elements)
			require.Equal(t, tt.expectedRows, len(rows))

			if tt.checkRowSizes {
				// Check that all rows except possibly the last one have equal size
				if len(rows) > 1 {
					firstRowSize := len(rows[0])
					for i := 1; i < len(rows)-1; i++ {
						require.Equal(t, firstRowSize, len(rows[i]))
					}
				}

				// Check that the total number of elements is preserved
				totalElements := 0
				for _, row := range rows {
					totalElements += len(row)
				}
				require.Equal(t, len(tt.elements), totalElements)
			}
		})
	}
}
