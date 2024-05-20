package types

import (
	"errors"
	"fmt"

	"cosmossdk.io/math"
)

func TicksToSqrtPrice(lowerTick, upperTick int64) (math.LegacyDec, math.LegacyDec, error) {
	if lowerTick >= upperTick {
		return math.LegacyDec{}, math.LegacyDec{}, errors.New("tickLower should be less than tickUpper")
	}
	sqrtPriceUpperTick, err := TickToSqrtPrice(upperTick)
	if err != nil {
		return math.LegacyDec{}, math.LegacyDec{}, err
	}
	sqrtPriceLowerTick, err := TickToSqrtPrice(lowerTick)
	if err != nil {
		return math.LegacyDec{}, math.LegacyDec{}, err
	}
	return sqrtPriceLowerTick, sqrtPriceUpperTick, nil
}

func TickToSqrtPrice(tickIndex int64) (math.LegacyDec, error) {
	priceDec, err := TickToPrice(tickIndex)
	if err != nil {
		return math.LegacyDec{}, err
	}

	if tickIndex >= TICK_MIN {
		sqrtPrice, err := priceDec.ApproxSqrt()
		if err != nil {
			return math.LegacyDec{}, err
		}
		return sqrtPrice, nil
	}

	sqrtPrice, err := priceDec.ApproxSqrt()
	if err != nil {
		return math.LegacyDec{}, err
	}
	return sqrtPrice, nil
}

func TickToPrice(tickIndex int64) (math.LegacyDec, error) {
	if tickIndex == 0 {
		return math.LegacyOneDec(), nil
	}

	if tickIndex == TICK_MIN {
		return MinSpotPrice, nil
	}

	numAdditiveTicks, geometricExponentDelta, err := TickToAdditiveGeometricIndices(tickIndex)
	if err != nil {
		return math.LegacyDec{}, err
	}

	exponentAtCurrentTick := ExponentAtPriceOne + geometricExponentDelta
	var unscaledPrice int64 = 1_000_000
	if tickIndex < 0 {
		exponentAtCurrentTick = exponentAtCurrentTick - 1
		unscaledPrice *= 10
	}
	unscaledPrice += numAdditiveTicks
	price := PowTenInternal(exponentAtCurrentTick).MulInt64(unscaledPrice)

	if price.GT(MaxSpotPrice) || price.LT(MinSpotPrice) {
		return math.LegacyDec{}, ErrPriceOutOfBound
	}
	return price, nil
}

func TickToAdditiveGeometricIndices(tickIndex int64) (additiveTicks int64, geometricExponentDelta int64, err error) {
	if tickIndex == 0 {
		return 0, 0, nil
	}

	if tickIndex == TICK_MIN {
		return 0, -18, nil
	}

	if tickIndex < TICK_MIN {
		return 0, 0, ErrInvalidTickers
	}
	if tickIndex > TICK_MAX {
		return 0, 0, ErrInvalidTickers
	}

	geometricExponentDelta = tickIndex / geometricExponentIncrementDistanceInTicks

	numAdditiveTicks := tickIndex - (geometricExponentDelta * geometricExponentIncrementDistanceInTicks)
	return numAdditiveTicks, geometricExponentDelta, nil
}

func RoundDownTickToSpacing(tickIndex int64, tickSpacing int64) (int64, error) {
	tickIndexModulus := tickIndex % tickSpacing
	if tickIndexModulus < 0 {
		tickIndexModulus += tickSpacing
	}

	if tickIndexModulus != 0 {
		tickIndex = tickIndex - tickIndexModulus
	}

	if tickIndex > TICK_MAX || tickIndex < TICK_MIN {
		return 0, ErrInvalidTickers
	}

	return tickIndex, nil
}

func SqrtPriceToTickRoundDownSpacing(sqrtPrice math.LegacyDec, tickSpacing uint64) (int64, error) {
	tickIndex, err := CalculateSqrtPriceToTick(sqrtPrice)
	if err != nil {
		return 0, err
	}

	tickIndex, err = RoundDownTickToSpacing(tickIndex, int64(tickSpacing))
	if err != nil {
		return 0, err
	}

	return tickIndex, nil
}

func PowTenInternal(exponent int64) math.LegacyDec {
	if exponent >= 0 {
		return powersOfTen[exponent]
	}
	return negPowersOfTen[-exponent]
}

func CalculatePriceToTick(price math.LegacyDec) (tickIndex int64, err error) {
	if price.IsNegative() {
		return 0, fmt.Errorf("price must be greater than zero")
	}
	if price.GT(MaxSpotPrice) || price.LT(MinSpotPrice) {
		return 0, ErrPriceOutOfBound
	}
	if price.Equal(sdkOneDec) {
		return 0, nil
	}

	var geoSpacing *tickExpIndexData
	if price.GT(sdkOneDec) {
		index := 0
		geoSpacing = tickExp[int64(index)]
		for geoSpacing.maxPrice.LT(price) {
			index += 1
			geoSpacing = tickExp[int64(index)]
		}
	} else {
		index := -1
		geoSpacing = tickExp[int64(index)]
		for geoSpacing.initialPrice.GT(price) {
			index -= 1
			geoSpacing = tickExp[int64(index)]
		}
	}

	priceInThisExponent := price.Sub(geoSpacing.initialPrice)
	ticksFilledByCurrentSpacing := priceInThisExponent.QuoMut(geoSpacing.additiveIncrementPerTick)
	tickIndex = ticksFilledByCurrentSpacing.TruncateInt64() + geoSpacing.initialTick
	return tickIndex, nil
}

func CalculateSqrtPriceToTick(sqrtPrice math.LegacyDec) (tickIndex int64, err error) {
	price := sqrtPrice.Mul(sqrtPrice)

	tick, err := CalculatePriceToTick(price)
	if err != nil {
		return 0, err
	}

	if tick < TICK_MIN {
		return 0, ErrInvalidTickers
	}

	outOfBounds := false
	if tick <= TICK_MIN {
		tick = TICK_MIN + 1
		outOfBounds = true
	} else if tick >= TICK_MAX-1 {
		tick = TICK_MAX - 2
		outOfBounds = true
	}

	sqrtPriceTplus1, err := TickToSqrtPrice(tick + 1)
	if err != nil {
		return 0, ErrSqrtPriceToTick
	}

	if sqrtPrice.GTE(sqrtPriceTplus1) {
		sqrtPriceTplus2, err := TickToSqrtPrice(tick + 2)
		if err != nil {
			return 0, ErrSqrtPriceToTick
		}
		if (!outOfBounds && sqrtPrice.GTE(sqrtPriceTplus2)) || (outOfBounds && sqrtPrice.GT(sqrtPriceTplus2)) {
			return 0, ErrSqrtPriceToTick
		}

		if sqrtPrice.Equal(sqrtPriceTplus2) {
			return tick + 2, nil
		}
		return tick + 1, nil
	}

	sqrtPriceT, err := TickToSqrtPrice(tick)
	if err != nil {
		return 0, ErrSqrtPriceToTick
	}
	if sqrtPrice.GTE(sqrtPriceT) {
		return tick, nil
	}

	sqrtPriceTmin1, err := TickToSqrtPrice(tick - 1)
	if err != nil {
		return 0, ErrSqrtPriceToTick
	}
	if sqrtPrice.LT(sqrtPriceTmin1) {
		return 0, ErrSqrtPriceToTick
	}

	return tick - 1, nil
}
