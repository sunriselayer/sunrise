package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Redelegation struct {
	Delegator    sdk.AccAddress
	SrcValidator LiquidValidator
	DstValidator LiquidValidator
	Amount       sdkmath.Int
	Last         bool
	Error        error
}

// DivideByWeight divide the input value by the ratio of the param weight of the liquid validator and return it with crumb
// which is may occur while dividing according to the weight of active liquid validators by decimal error.
func DivideByWeight(avs ActiveLiquidValidators, input sdkmath.Int, whitelistedValsMap WhitelistedValsMap) (outputs []sdkmath.Int, crumb sdkmath.Int) {
	totalWeight := avs.TotalWeight(whitelistedValsMap)
	if !totalWeight.IsPositive() {
		return []sdkmath.Int{}, sdkmath.ZeroInt()
	}
	totalOutput := sdkmath.ZeroInt()
	unitInput := input.ToLegacyDec().QuoTruncate(totalWeight.ToLegacyDec())
	for _, val := range avs {
		output := unitInput.MulInt(val.GetWeight(whitelistedValsMap, true)).TruncateInt()
		totalOutput = totalOutput.Add(output)
		outputs = append(outputs, output)
	}
	return outputs, input.Sub(totalOutput)
}

// DivideByCurrentWeight divide the input value by the ratio of the weight of the liquid validator's liquid token and return it with crumb
// which is may occur while dividing according to the weight of liquid validators by decimal error, outputs is truncated decimal.
func DivideByCurrentWeight(lvs LiquidValidators, input sdkmath.LegacyDec, totalLiquidTokens sdkmath.Int, liquidTokenMap map[string]sdkmath.Int) (outputs []sdkmath.LegacyDec, crumb sdkmath.LegacyDec) {
	if !totalLiquidTokens.IsPositive() {
		return []sdkmath.LegacyDec{}, sdkmath.LegacyZeroDec()
	}
	totalOutput := sdkmath.LegacyZeroDec()
	unitInput := input.QuoTruncate(totalLiquidTokens.ToLegacyDec())
	for _, val := range lvs {
		output := unitInput.MulTruncate(liquidTokenMap[val.OperatorAddress].ToLegacyDec()).TruncateDec()
		totalOutput = totalOutput.Add(output)
		outputs = append(outputs, output)
	}
	return outputs, input.Sub(totalOutput)
}
