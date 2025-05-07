package erasurecoding

import (
	"bytes"

	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr/fft"

	"github.com/sunriselayer/sunrise/x/da/das/consts"
)

func padTo32(input []byte) []byte {
	if len(input)%32 == 0 {
		return input
	}
	padding := 32 - (len(input) % 32)
	return append(input, bytes.Repeat([]byte{0}, padding)...)
}

func ConvertToElements(input []byte) []fr.Element {
	padded32 := padTo32(input)
	numElems := len(padded32) / 32
	elements := make([]fr.Element, numElems)
	for i := range numElems {
		var tmp [32]byte
		copy(tmp[:], padded32[i*32:(i+1)*32])
		elements[i].SetBytes(tmp[:])
	}
	return elements
}

func SplitElementsIntoRows(elements []fr.Element) ([][]fr.Element, uint32, uint32) {
	numElems := len(elements)
	domain := fft.NewDomain(uint64(numElems))

	if domain.Cardinality > consts.SrsLen {
		// Calculate how many rows we need
		numRows := 2
		for {
			rowSize := (numElems + numRows - 1) / numRows
			if uint64(rowSize) <= consts.SrsLen {
				break
			}
			numRows *= 2
		}

		// Calculate row size and remainder
		rowSize := numElems / numRows
		remainder := numElems % numRows

		// Create rows with equal size where possible
		// Extra elements (if any) are added to the first rows in order
		rows := make([][]fr.Element, numRows)
		start := 0
		maxRowSize := rowSize
		for rowIndex := range numRows {
			// Add one extra element to the first 'remainder' rows
			currentRowSize := rowSize
			if rowIndex < remainder {
				currentRowSize++

			}
			rows[rowIndex] = elements[start : start+currentRowSize]
			start += currentRowSize

			maxRowSize = max(maxRowSize, currentRowSize)
		}
		return rows, uint32(numRows), uint32(maxRowSize)
	}

	// If we don't need to split, return the original elements as a single row
	return [][]fr.Element{elements}, 1, uint32(numElems)
}
