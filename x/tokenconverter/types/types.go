package types

import (
	"cosmossdk.io/math"
)

func ratio(maxSupplyFeeToken math.Int, supplyGovToken math.Int) math.LegacyDec {
	feeDec := math.LegacyNewDecFromInt(maxSupplyFeeToken)
	govDec := math.LegacyNewDecFromInt(supplyGovToken)
	ratio := feeDec.Quo(govDec)

	if ratio.GT(math.LegacyOneDec()) {
		ratio = math.LegacyOneDec()
	}

	return ratio

}

func CalculateAmountOutFeeToken(maxSupplyFeeToken math.Int, supplyGovToken math.Int, amountInGovToken math.Int) math.Int {
	ratio := ratio(maxSupplyFeeToken, supplyGovToken)

	return ratio.MulInt(amountInGovToken).TruncateInt()
}

func CalculateAmountInGovToken(maxSupplyFeeToken math.Int, supplyGovToken math.Int, amountOutFeeToken math.Int) math.Int {
	ratio := ratio(maxSupplyFeeToken, supplyGovToken)

	return math.LegacyOneDec().Quo(ratio).MulInt(amountOutFeeToken).TruncateInt()
}
