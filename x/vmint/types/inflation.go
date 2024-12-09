package types

import (
	"context"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	DefaultInitialSupply    = math.NewInt(500_000_000).Mul(math.NewInt(1_000_000))
	DefaultInflationRateCap = math.LegacyMustNewDecFromStr("0.1")
	DefaultDisinflationRate = math.LegacyMustNewDecFromStr("0.08")
	DefaultSupplyCap        = math.NewInt(1_000_000_000).Mul(math.NewInt(1_000_000))
)

func InflationRate(
	ctx context.Context,
	genesis time.Time,
	inflationRateCapInitial math.LegacyDec,
	inflationRateCapMinimum math.LegacyDec,
	disinflationRate math.LegacyDec,
	supplyCap math.Int,
	totalSupply math.Int,
) math.LegacyDec {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	years := yearsSinceGenesis(genesis, sdkCtx.BlockTime())

	inflationRateCap := inflationRateCapInitial.Mul(math.LegacyOneDec().Sub(disinflationRate).Power(years))
	if inflationRateCap.LT(inflationRateCapMinimum) {
		inflationRateCap = inflationRateCapMinimum
	}

	nextSupply := math.LegacyOneDec().Add(inflationRateCap).Power(years).MulInt(totalSupply).TruncateInt()

	var inflationRate math.LegacyDec
	if nextSupply.GT(supplyCap) {
		nextSupply = supplyCap
		// 1 + π = s'/s
		inflationRate = math.LegacyNewDecFromInt(nextSupply).QuoInt(totalSupply).Sub(math.LegacyOneDec())
	} else {
		inflationRate = inflationRateCap
	}

	return inflationRate
}

func yearsSinceGenesis(genesis time.Time, current time.Time) (years uint64) {
	if current.Before(genesis) {
		return 0
	}
	const millisecondsPerYear = 31556952000

	return uint64(current.Sub(genesis).Milliseconds() / millisecondsPerYear)
}
