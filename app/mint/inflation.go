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
) math.Int {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	years := yearsSinceGenesis(genesis, sdkCtx.BlockTime())

	inflationRateCap := inflationRateCapInitial.Mul(math.LegacyOneDec().Sub(disinflationRate).Power(years))
	if inflationRateCap.LT(inflationRateCapMinimum) {
		inflationRateCap = inflationRateCapMinimum
	}

	nextSupply := math.LegacyOneDec().Add(inflationRateCap).MulInt(totalSupply).TruncateInt()

	if nextSupply.GT(supplyCap) {
		nextSupply = supplyCap
	}
	if nextSupply.LT(totalSupply) {
		nextSupply = totalSupply
	}

	return nextSupply.Sub(totalSupply)
}

func yearsSinceGenesis(genesis time.Time, current time.Time) (years uint64) {
	if current.Before(genesis) {
		return 0
	}
	const millisecondsPerYear = 31556952000

	return uint64(current.Sub(genesis).Milliseconds() / millisecondsPerYear)
}
