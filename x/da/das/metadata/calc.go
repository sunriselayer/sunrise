package metadata

import (
	"github.com/sunriselayer/sunrise/x/da/das/consts"
	"github.com/sunriselayer/sunrise/x/da/das/types"
)

func CalculateTotalSize(rows uint32, cols uint32) uint64 {
	return uint64(rows) * uint64(cols) * consts.ElementSize
}

func CalculateOpeningProofCount(rows uint32, cols uint32) uint32 {
	i := uint32(0)

	for cols&(cols-1) == 0 {
		cols >>= 1
		i++
	}

	return rows * i
}

func CalculateOpeningProofIndices(rows uint32, cols uint32) ([]types.OpeningProofIndex, error) {

}
