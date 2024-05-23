package swapstrategy

import (
	"fmt"

	"cosmossdk.io/math"
)

type spreadFactorOverOneMinusSpreadFactorGetter func() math.LegacyDec

func computeSpreadRewardChargePerSwapStepOutGivenIn(hasReachedTarget bool, amountIn, amountSpecifiedRemaining, spreadFactor math.LegacyDec, SfOveroneMinSf spreadFactorOverOneMinusSpreadFactorGetter) math.LegacyDec {
	if spreadFactor.IsZero() {
		return math.LegacyZeroDec()
	} else if spreadFactor.IsNegative() {
		panic(fmt.Errorf("spread factor must be non-negative, was (%s)", spreadFactor))
	}

	var spreadRewardChargeTotal math.LegacyDec
	if hasReachedTarget {
		spreadRewardChargeTotal = computeSpreadRewardChargeFromAmountIn(amountIn, SfOveroneMinSf())
	} else {
		spreadRewardChargeTotal = amountSpecifiedRemaining.Sub(amountIn)
	}

	if spreadRewardChargeTotal.IsNegative() {
		panic(fmt.Errorf("spread factor charge must be non-negative, was (%s)", spreadRewardChargeTotal))
	}

	return spreadRewardChargeTotal
}

func computeSpreadRewardChargeFromAmountIn(amountIn math.LegacyDec, spreadFactorOverOneMinusSpreadFactor math.LegacyDec) math.LegacyDec {
	return amountIn.MulRoundUp(spreadFactorOverOneMinusSpreadFactor)
}
