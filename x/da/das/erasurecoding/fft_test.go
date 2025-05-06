package erasurecoding

import (
	"testing"

	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr/fft"
	"github.com/stretchr/testify/require"

	"github.com/sunriselayer/sunrise/x/da/das/consts"
)

func TestCalculateCoefficients(t *testing.T) {
	tests := []struct {
		name           string
		elements       []fr.Element
		expectedLen    int
		expectedDomain uint64
	}{
		{
			name: "power of 2 length",
			elements: []fr.Element{
				fr.NewElement(1),
				fr.NewElement(2),
				fr.NewElement(4),
				fr.NewElement(8),
			},
			expectedLen: 4,
		},
		/*
			{
				name: "non-power of 2 length",
				elements: []fr.Element{
					fr.NewElement(2),
					fr.NewElement(7),
					fr.NewElement(4),
				},
				expectedLen: 4,
			},*/
		{
			name:        "empty elements",
			elements:    []fr.Element{},
			expectedLen: 0,
		},
		/*
			{
				name:        "large input size",
				elements:    make([]fr.Element, 1024),
				expectedLen: 1024,
			},
				{
					name:        "boundary case",
					elements:    make([]fr.Element, consts.ElementsLenPerShard-1),
					expectedLen: consts.ElementsLenPerShard,
				},*/
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			coeffs := CalculateCoefficients(tt.elements)
			require.Equal(t, tt.expectedLen, len(coeffs))

			domain := fft.NewDomain(uint64(tt.expectedLen))
			recover := make([]fr.Element, tt.expectedLen)
			copy(recover, coeffs)
			if tt.expectedLen > 0 {
				fft.BitReverse(recover)
				domain.FFT(recover, fft.DIT)
			}

			for i := range tt.elements {
				require.Equal(t, tt.elements[i], recover[i])
			}
		})
	}
}

func TestCalculateExtendedEvaluationPoints(t *testing.T) {
	tests := []struct {
		name        string
		coeffs      []fr.Element
		expectedLen int
		wantErr     bool
	}{
		{
			name:        "normal case",
			coeffs:      []fr.Element{fr.NewElement(1), fr.NewElement(2)},
			expectedLen: 128, // Should be at least ElementsLenPerShard
			wantErr:     false,
		},
		{
			name:        "small domain cardinality",
			coeffs:      []fr.Element{fr.NewElement(1)},
			expectedLen: 128, // Should be at least ElementsLenPerShard
			wantErr:     false,
		},
		{
			name:        "empty coefficients",
			coeffs:      []fr.Element{},
			expectedLen: 0,
			wantErr:     true,
		},
		{
			name:        "large input size",
			coeffs:      make([]fr.Element, 1024),
			expectedLen: 2 * 1024,
			wantErr:     false,
		},
		{
			name:        "boundary case",
			coeffs:      make([]fr.Element, consts.ElementsLenPerShard-1),
			expectedLen: 2 * consts.ElementsLenPerShard,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			points, err := CalculateExtendedEvaluationPoints(tt.coeffs)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.expectedLen, len(points))
		})
	}
}

func TestEvaluationPoints(t *testing.T) {
	tests := []struct {
		name     string
		len      uint64
		expected int
	}{
		{
			name:     "power of 2",
			len:      4,
			expected: 4,
		},
		{
			name:     "non-power of 2",
			len:      3,
			expected: 3,
		},
		{
			name:     "zero length",
			len:      0,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			points := EvaluationPoints(tt.len)
			require.Equal(t, tt.expected, len(points))

			// Check that points are properly initialized
			if len(points) > 0 {
				require.True(t, points[0].IsOne())
				for i := 1; i < len(points); i++ {
					require.False(t, points[i].IsZero())
				}
			}
		})
	}
}

func TestFFTFullProcess(t *testing.T) {
	// Test the complete FFT process
	originalElements := []fr.Element{
		fr.NewElement(1),
		fr.NewElement(3),
		fr.NewElement(37),
		fr.NewElement(28),
	}

	// Calculate coefficients
	coeffs := CalculateCoefficients(originalElements)
	require.Equal(t, 4, len(coeffs))

	// Calculate extended evaluation points
	extendedPoints, err := CalculateExtendedEvaluationPoints(coeffs)
	require.NoError(t, err)
	require.Equal(t, 128, len(extendedPoints))

	// Verify that the extended points preserve original elements
	for i := range originalElements {
		require.Equal(t, originalElements[i], extendedPoints[i])
	}

	// Verify that the extended points are not zero
	for i := range extendedPoints {
		require.False(t, extendedPoints[i].IsZero())
	}
}
