package types

import (
	"cosmossdk.io/math"
)

func AmountsFromWeights(weights []PoolWeight, totalAmount math.Int) []math.Int {
	length := len(weights)
	amounts := make([]math.Int, length)
	sum := math.ZeroInt()

	for i := 0; i < length-1; i++ {
		amounts[i] = weights[i].Weight.MulInt(totalAmount).RoundInt()
		sum = sum.Add(amounts[i])
	}
	amounts[length-1] = totalAmount.Sub(sum)

	return amounts
}
