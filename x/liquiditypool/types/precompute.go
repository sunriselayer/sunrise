package types

import (
	"cosmossdk.io/math"
)

var (
	ExponentAtPriceOne int64 = -6
	MaxSpotPrice             = math.LegacyMustNewDecFromStr("100000000000000000000000000000000000000")
	MinSpotPrice             = math.LegacyNewDecWithPrec(1, 18)

	sdkOneDec      = math.LegacyOneDec()
	sdkTenDec      = math.LegacyNewDec(10)
	powersOfTen    []math.LegacyDec
	negPowersOfTen []math.LegacyDec

	geometricExponentIncrementDistanceInTicks = 9 * math.LegacyNewDec(10).PowerMut(uint64(-ExponentAtPriceOne)).TruncateInt64()
)

type tickExpIndexData struct {
	// if price < initialPrice, we are not in this exponent range.
	initialPrice math.LegacyDec
	// if price >= maxPrice, we are not in this exponent range.
	maxPrice math.LegacyDec
	// additive increment per tick here.
	additiveIncrementPerTick math.LegacyDec
	// the tick that corresponds to initial price
	initialTick int64
}

var tickExp map[int64]*tickExpIndexData = make(map[int64]*tickExpIndexData)

func buildTickExp() {
	// build positive indices first
	maxPrice := sdkOneDec
	curExpIndex := int64(0)
	for maxPrice.LT(MaxSpotPrice) {
		tickExp[curExpIndex] = &tickExpIndexData{
			// price range 10^curExpIndex to 10^(curExpIndex + 1). (10, 100)
			initialPrice:             sdkTenDec.Power(uint64(curExpIndex)),
			maxPrice:                 sdkTenDec.Power(uint64(curExpIndex + 1)),
			additiveIncrementPerTick: PowTenInternal(ExponentAtPriceOne + curExpIndex),
			initialTick:              geometricExponentIncrementDistanceInTicks * curExpIndex,
		}
		maxPrice = tickExp[curExpIndex].maxPrice
		curExpIndex += 1
	}

	minPrice := sdkOneDec
	curExpIndex = -1
	for minPrice.GT(math.LegacyNewDecWithPrec(1, 18)) {
		tickExp[curExpIndex] = &tickExpIndexData{
			// price range 10^curExpIndex to 10^(curExpIndex + 1). (0.001, 0.01)
			initialPrice:             PowTenInternal(curExpIndex),
			maxPrice:                 PowTenInternal(curExpIndex + 1),
			additiveIncrementPerTick: PowTenInternal(ExponentAtPriceOne + curExpIndex),
			initialTick:              geometricExponentIncrementDistanceInTicks * curExpIndex,
		}
		minPrice = tickExp[curExpIndex].initialPrice
		curExpIndex -= 1
	}
}

// Set precision multipliers
func init() {
	negPowersOfTen = make([]math.LegacyDec, math.LegacyPrecision+1)
	for i := 0; i <= math.LegacyPrecision; i++ {
		negPowersOfTen[i] = sdkOneDec.Quo(sdkTenDec.Power(uint64(i)))
	}
	// 10^77 < math.MaxInt < 10^78
	powersOfTen = make([]math.LegacyDec, 77)
	for i := 0; i <= 76; i++ {
		powersOfTen[i] = sdkTenDec.Power(uint64(i))
	}

	buildTickExp()
}
