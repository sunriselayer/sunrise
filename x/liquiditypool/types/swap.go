package types

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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

func ValidateWeights(weights []PoolWeight) error {
	if len(weights) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "weights cannot be empty")
	}

	sum := math.LegacyZeroDec()

	for _, weight := range weights {
		if weight.Weight.LTE(math.LegacyZeroDec()) {
			return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "weight must be positive")
		}

		sum = sum.Add(weight.Weight)
	}

	if !sum.Equal(math.LegacyOneDec()) {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "sum of weights must be 1")
	}

	return nil
}
