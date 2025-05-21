package mint

import (
	"context"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func CalculateAnnualProvision(
	ctx context.Context,
	inflationRateCapInitial math.LegacyDec,
	inflationRateCapMinimum math.LegacyDec,
	disinflationRate math.LegacyDec,
	supplyCap math.Int,
	genesis time.Time,
	totalSupply math.Int,
) math.LegacyDec {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	years := yearsSinceGenesis(genesis, sdkCtx.BlockTime())

	inflationRateCap := inflationRateCapInitial.Mul(math.LegacyOneDec().Sub(disinflationRate).Power(years))
	if inflationRateCap.LT(inflationRateCapMinimum) {
		inflationRateCap = inflationRateCapMinimum
	}

	nextSupply := math.LegacyOneDec().Add(inflationRateCap).MulInt(totalSupply)

	supplyCapDec := math.LegacyNewDecFromInt(supplyCap)
	totalSupplyDec := math.LegacyNewDecFromInt(totalSupply)

	if nextSupply.GT(supplyCapDec) {
		nextSupply = supplyCapDec
	}
	if nextSupply.LT(totalSupplyDec) {
		nextSupply = totalSupplyDec
	}

	return nextSupply.Sub(totalSupplyDec)
}

func yearsSinceGenesis(genesis time.Time, current time.Time) (years uint64) {
	if current.Before(genesis) {
		return 0
	}
	const millisecondsPerYear = 31536000000

	return uint64(current.Sub(genesis).Milliseconds() / millisecondsPerYear)
}
