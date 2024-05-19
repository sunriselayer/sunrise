package types

import (
	"cosmossdk.io/math"
)

func ratio(supplyFeeToken math.Int, supplyGovToken math.Int) math.LegacyDec {
	feeDec := math.LegacyNewDecFromInt(supplyFeeToken)
	govDec := math.LegacyNewDecFromInt(supplyGovToken)
	ratio := feeDec.Quo(govDec)

	if ratio.GT(math.LegacyOneDec()) {
		ratio = math.LegacyOneDec()
	}

	return ratio

}

func CalculateAmountOutFeeToken(supplyFeeToken math.Int, supplyGovToken math.Int, amountInGovToken math.Int) math.Int {
	ratio := ratio(supplyFeeToken, supplyGovToken)

	return ratio.MulInt(amountInGovToken).TruncateInt()
}

func CalculateAmountInGovToken(supplyFeeToken math.Int, supplyGovToken math.Int, amountOutFeeToken math.Int) math.Int {
	ratio := ratio(supplyFeeToken, supplyGovToken)

	return math.LegacyOneDec().Quo(ratio).MulInt(amountOutFeeToken).TruncateInt()
}
