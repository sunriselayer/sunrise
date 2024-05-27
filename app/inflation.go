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
)

var (
	// initialInflationRate is the inflation rate that the network starts at.
	initialInflationRate = sdkmath.LegacyMustNewDecFromStr("0.08")
	// initialInflationRate is the max inflation rate at the first year based on bondedRatio
	initialInflationRateMax = sdkmath.LegacyMustNewDecFromStr("0.10")
	// initialInflationRate is the min inflation rate at the first year based on bondedRatio
	initialInflationRateMin = sdkmath.LegacyMustNewDecFromStr("0.06")

	// disinflationRate is the rate at which the inflation rate decreases each year.
	disinflationRate = sdkmath.LegacyMustNewDecFromStr("0.08")
	// targetInflationRate is the inflation rate that the network aims to
	// stabilize at. In practice, targetInflationRate acts as a minimum so that
	// the inflation rate doesn't decrease after reaching it.
	targetInflationRate = sdkmath.LegacyMustNewDecFromStr("0.02")
)

func InitialInflationRate() sdkmath.LegacyDec {
	return initialInflationRate
}

func InitialInflationRateMax() sdkmath.LegacyDec {
	return initialInflationRateMax
}

func InitialInflationRateMin() sdkmath.LegacyDec {
	return initialInflationRateMin
}

func DisinflationRate() sdkmath.LegacyDec {
	return disinflationRate
}

func TargetInflationRate() sdkmath.LegacyDec {
	return targetInflationRate
}

func InflationCalculationFn(ctx context.Context, minter minttypes.Minter, params minttypes.Params, bondedRatio sdkmath.LegacyDec) sdkmath.LegacyDec {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	return CalculateInflationRate(sdkCtx, time.Unix(genesisTime, 0), bondedRatio)
}

func ProvideInflationCalculatorFn() minttypes.InflationCalculationFn {
	return InflationCalculationFn
}

// CalculateInflationRate returns the inflation rate for the current year depending on
// the current block height in context. The inflation rate is expected to
// decrease every year according to the schedule specified in the README.
func CalculateInflationRate(ctx sdk.Context, genesis time.Time, bondedRatio sdkmath.LegacyDec) sdkmath.LegacyDec {
	// initialRate = initialMax - (initialMax-initialMin)*bondedRatio
	initialRate := initialInflationRateMax.Sub(
		initialInflationRateMax.Sub(initialInflationRateMin).Mul(bondedRatio),
	)

	// disinflatedRate = initialRate * (1 - disinflationRate)^((now - genesis).convertToYears())
	years := yearsSinceGenesis(genesis, ctx.BlockTime())
	disinflatedRate := initialRate.Mul(
		sdkmath.LegacyOneDec().Sub(disinflationRate).Power(uint64(years)),
	)

	// finalRate = max(disinflatedRate, convergenceRate)
	return sdkmath.LegacyMaxDec(disinflatedRate, TargetInflationRate())
}

// yearsSinceGenesis returns the number of years that have passed between
// genesis and current (rounded down).
func yearsSinceGenesis(genesis time.Time, current time.Time) (years int64) {
	if current.Before(genesis) {
		return 0
	}
	return current.Sub(genesis).Nanoseconds() / nanosecondsPerYear
}
