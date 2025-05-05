package kzg

import (
	"bytes"
	"fmt"
	"math"

	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
)

func padTo32(input []byte) []byte {
	if len(input)%32 == 0 {
		return input
	}
	padding := 32 - (len(input) % 32)
	return append(input, bytes.Repeat([]byte{0}, padding)...)
}

func ConvertToElements(input []byte) ([]fr.Element, error) {
	padded32 := padTo32(input)
	if len(padded32) > MaxBlobSize {
		return nil, fmt.Errorf("input size must be less than or equal to %d: %d", MaxBlobSize, len(padded32))
	}
	numElems := len(padded32) / 32
	elements := make([]fr.Element, numElems)
	for i := range numElems {
		var tmp [32]byte
		copy(tmp[:], padded32[i*32:(i+1)*32])
		elements[i].SetBytes(tmp[:])
	}
	return elements, nil
}

// SquareLayout returns the width and height of the closest square-like rectangle
// that can contain n elements. The function tries to find dimensions that are
// as close as possible to a square while ensuring width * height >= n.
// It prefers smaller width and larger height for better visual layout.
func SquareLayout(n int) [2]int {
	if n <= 0 {
		return [2]int{0, 0}
	}

	// Start with the square root of n
	side := int(math.Ceil(math.Sqrt(float64(n))))

	// Find the closest dimensions that can contain n elements
	bestWidth := side
	bestHeight := side
	minDiff := math.MaxFloat64

	// Check dimensions around the square root
	for w := side; w > 0; w-- {
		h := int(math.Ceil(float64(n) / float64(w)))
		diff := math.Abs(float64(w - h))
		// Prefer smaller width when the difference is the same
		if diff <= minDiff {
			minDiff = diff
			bestWidth = w
			bestHeight = h
		}
	}

	return [2]int{bestWidth, bestHeight}
}

// ConvertToElementsSquare converts a flat slice of elements into a 2D square-like layout.
// The elements are packed in a way that minimizes empty spaces and maintains a square-like shape.
// The function allocates a square memory space using the larger of width and height.
func ConvertToElementsSquare(elements []fr.Element) [][]fr.Element {
	if len(elements) == 0 {
		return nil
	}

	layout := SquareLayout(len(elements))
	width, height := layout[0], layout[1]

	// Use the larger dimension for memory allocation
	dim := max(width, height)

	// Create a square 2D slice
	square := make([][]fr.Element, dim)
	for i := range dim {
		square[i] = make([]fr.Element, dim)
	}

	// Pack elements row by row according to width and height
	idx := 0
	for i := range height {
		for j := range width {
			if idx < len(elements) {
				square[i][j] = elements[idx]
				idx++
			}
		}
	}

	return square
}

// ExtendSquare performs 2D erasure coding on a square matrix using FFT.
// It extends both columns and rows to create a larger square matrix.
// Column extension is performed first, followed by row extension for better efficiency.
// Returns the extended square matrix and coefficients for both extensions.
// coeffs[0] contains coefficients for vertical extension (length = original dimension)
// coeffs[1] contains coefficients for horizontal extension (length = extended dimension)
func ExtendSquare(square [][]fr.Element) ([][]fr.Element, [][]fr.Element) {
	if len(square) == 0 {
		return nil, nil
	}

	// First, extend each column
	extendedColumns := make([][]fr.Element, len(square[0]))
	verticalCoeffs := make([]fr.Element, len(square)) // Store coefficients for vertical extension
	for j := range len(square[0]) {
		// Extract column
		column := make([]fr.Element, len(square))
		for i := range len(square) {
			column[i] = square[i][j]
		}

		// Calculate coefficients and extended evaluation points for the column
		coeffs, domainCardinality := CalculateCoefficients(column)
		extendedColumn := CalculateExtendedEvaluationPoints(coeffs, domainCardinality)

		// Store the extended column
		extendedColumns[j] = extendedColumn

		// Store coefficients for the first column only (they are the same for all columns)
		if j == 0 {
			copy(verticalCoeffs, coeffs[:len(square)])
		}
	}

	// Then, extend each row of the extended columns
	extendedDim := len(extendedColumns[0])
	extendedSquare := make([][]fr.Element, extendedDim)
	horizontalCoeffs := make([]fr.Element, extendedDim) // Store coefficients for horizontal extension
	for i := range extendedDim {
		// Extract row from extended columns
		row := make([]fr.Element, len(extendedColumns))
		for j := range len(extendedColumns) {
			row[j] = extendedColumns[j][i]
		}

		// Calculate coefficients and extended evaluation points for the row
		coeffs, domainCardinality := CalculateCoefficients(row)
		extendedRow := CalculateExtendedEvaluationPoints(coeffs, domainCardinality)

		// Store the extended row
		extendedSquare[i] = extendedRow

		// Store coefficients for the first row only (they are the same for all rows)
		if i == 0 {
			copy(horizontalCoeffs, coeffs[:extendedDim])
		}
	}

	return extendedSquare, [][]fr.Element{verticalCoeffs, horizontalCoeffs}
}
