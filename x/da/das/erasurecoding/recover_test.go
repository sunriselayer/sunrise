package erasurecoding

import (
	"testing"

	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
	"github.com/stretchr/testify/require"
)

func TestRecoverCoefficients(t *testing.T) {
	tests := []struct {
		name              string
		extendedElements  map[uint64][]fr.Element
		domainCardinality uint64
		wantErr           bool
	}{
		{
			name: "valid case",
			extendedElements: map[uint64][]fr.Element{
				0: {fr.NewElement(1)},
				1: {fr.NewElement(2)},
				2: {fr.NewElement(3)},
				3: {fr.NewElement(4)},
			},
			domainCardinality: 4,
			wantErr:           false,
		},
		{
			name: "invalid domain cardinality (not power of 2)",
			extendedElements: map[uint64][]fr.Element{
				0: {fr.NewElement(1)},
				1: {fr.NewElement(2)},
				2: {fr.NewElement(3)},
			},
			domainCardinality: 3,
			wantErr:           true,
		},
		{
			name: "mismatched length",
			extendedElements: map[uint64][]fr.Element{
				0: {fr.NewElement(1)},
				1: {fr.NewElement(2)},
			},
			domainCardinality: 4,
			wantErr:           true,
		},
		{
			name:              "empty elements",
			extendedElements:  map[uint64][]fr.Element{},
			domainCardinality: 4,
			wantErr:           true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			coeffs, err := RecoverCoefficients(tt.extendedElements, tt.domainCardinality)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, int(tt.domainCardinality), len(coeffs))
		})
	}
}

func TestRecoverElements(t *testing.T) {
	tests := []struct {
		name     string
		coeffs   []fr.Element
		elemsLen uint64
		expected int
	}{
		{
			name:     "valid case",
			coeffs:   []fr.Element{fr.NewElement(1), fr.NewElement(2), fr.NewElement(3)},
			elemsLen: 3,
			expected: 3,
		},
		{
			name:     "more coefficients than elements",
			coeffs:   []fr.Element{fr.NewElement(1), fr.NewElement(2), fr.NewElement(3), fr.NewElement(4)},
			elemsLen: 2,
			expected: 2,
		},
		{
			name:     "empty coefficients",
			coeffs:   []fr.Element{},
			elemsLen: 0,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			elements := RecoverElements(tt.coeffs, tt.elemsLen)
			require.Equal(t, tt.expected, len(elements))
		})
	}
}

func TestRecoverCoefficientsAndElements(t *testing.T) {
	// Test the full recovery process
	originalElements := []fr.Element{
		fr.NewElement(1),
		fr.NewElement(3),
		fr.NewElement(37),
		fr.NewElement(28),
	}

	// Create extended elements map
	extendedElements := make(map[uint64][]fr.Element)
	for i, elem := range originalElements {
		extendedElements[uint64(i)] = []fr.Element{elem}
	}

	// Recover coefficients
	coeffs, err := RecoverCoefficients(extendedElements, 4)
	require.NoError(t, err)
	require.Equal(t, 4, len(coeffs))

	// Recover elements from coefficients
	recoveredElements := RecoverElements(coeffs, 4)
	require.Equal(t, len(originalElements), len(recoveredElements))

	// Compare original and recovered elements
	for i := range originalElements {
		require.Equal(t, originalElements[i], recoveredElements[i])
	}
}
