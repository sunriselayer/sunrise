package app

import (
	"context"
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

var _ minttypes.InflationCalculationFn = InflationCalculationFn

const (
	nanosecondsPerSecond = 1_000_000_000
	secondsPerMinute     = 60
	minutesPerHour       = 60
	hoursPerDay          = 24
	// daysPerYear is the mean length of the Gregorian calendar year. Note this
	// value isn't 365 because 97 out of 400 years are leap years. See
	// https://en.wikipedia.org/wiki/Year
	daysPerYear        = 365.2425
	secondsPerYear     = int64(secondsPerMinute * minutesPerHour * hoursPerDay * daysPerYear) // 31,556,952
	nanosecondsPerYear = nanosecondsPerSecond * secondsPerYear                                // 31,556,952,000,000,000

	// genesis time - Tuesday, May 7, 2024 2:42:49 AM UTC
	genesisTime = int64(1715049769)

	// initialInflationRate is the inflation rate that the network starts at.
	initialInflationRate = 0.08
	// initialInflationRate is the max inflation rate at the first year based on bondedRatio
	initialInflationRateMax = 0.1
	// initialInflationRate is the min inflation rate at the first year based on bondedRatio
	initialInflationRateMin = 0.06

	// disinflationRate is the rate at which the inflation rate decreases each year.
	disinflationRate = 0.08
	// targetInflationRate is the inflation rate that the network aims to
	// stabilize at. In practice, targetInflationRate acts as a minimum so that
	// the inflation rate doesn't decrease after reaching it.
	targetInflationRate = 0.02
)

var (
	initialInflationRateAsDec    = sdkmath.LegacyNewDecWithPrec(initialInflationRate*1000, 3)
	initialInflationRateMaxAsDec = sdkmath.LegacyNewDecWithPrec(initialInflationRateMax*1000, 3)
	initialInflationRateMinAsDec = sdkmath.LegacyNewDecWithPrec(initialInflationRateMin*1000, 3)
	disinflationRateAsDec        = sdkmath.LegacyNewDecWithPrec(disinflationRate*1000, 3)
	targetInflationRateAsDec     = sdkmath.LegacyNewDecWithPrec(targetInflationRate*1000, 3)
)

func InitialInflationRateAsDec() sdkmath.LegacyDec {
	return initialInflationRateAsDec
}

func InitialInflationRateMaxAsDec() sdkmath.LegacyDec {
	return initialInflationRateMaxAsDec
}

func InitialInflationRateMinAsDec() sdkmath.LegacyDec {
	return initialInflationRateMinAsDec
}

func DisinflationRateAsDec() sdkmath.LegacyDec {
	return disinflationRateAsDec
}

func TargetInflationRateAsDec() sdkmath.LegacyDec {
	return targetInflationRateAsDec
}

func InflationCalculationFn(ctx context.Context, minter minttypes.Minter, params minttypes.Params, bondedRatio sdkmath.LegacyDec) sdkmath.LegacyDec {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	return CalculateInflationRate(sdkCtx, time.Unix(genesisTime, 0), bondedRatio)
}

// CalculateInflationRate returns the inflation rate for the current year depending on
// the current block height in context. The inflation rate is expected to
// decrease every year according to the schedule specified in the README.
func CalculateInflationRate(ctx sdk.Context, genesis time.Time, bondedRatio sdkmath.LegacyDec) sdkmath.LegacyDec {
	// initialRate = initialMax - (initialMax-initialMin)*bondedRatio
	initialRate := initialInflationRateMaxAsDec.Sub(
		initialInflationRateMaxAsDec.Sub(initialInflationRateMinAsDec).Mul(bondedRatio),
	)

	// disinflatedRate = initialRate * (1 - disinflationRate)^((now - genesis).convertToYears())
	years := yearsSinceGenesis(genesis, ctx.BlockTime())
	disinflatedRate := initialRate.Mul(
		sdkmath.LegacyOneDec().Sub(disinflationRateAsDec).Power(uint64(years)),
	)

	// finalRate = max(disinflatedRate, convergenceRate)
	return sdkmath.LegacyMaxDec(disinflatedRate, TargetInflationRateAsDec())
}

// yearsSinceGenesis returns the number of years that have passed between
// genesis and current (rounded down).
func yearsSinceGenesis(genesis time.Time, current time.Time) (years int64) {
	if current.Before(genesis) {
		return 0
	}
	return current.Sub(genesis).Nanoseconds() / nanosecondsPerYear
}
