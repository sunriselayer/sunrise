package types

import (
	"errors"
	"fmt"

	"cosmossdk.io/math"
)

func TicksToSqrtPrice(lowerTick, upperTick int64, tickParams TickParams) (math.LegacyDec, math.LegacyDec, error) {
	if lowerTick >= upperTick {
		return math.LegacyDec{}, math.LegacyDec{}, errors.New("tickLower should be less than tickUpper")
	}
	sqrtPriceUpperTick, err := TickToSqrtPrice(upperTick, tickParams)
	if err != nil {
		return math.LegacyDec{}, math.LegacyDec{}, err
	}
	sqrtPriceLowerTick, err := TickToSqrtPrice(lowerTick, tickParams)
	if err != nil {
		return math.LegacyDec{}, math.LegacyDec{}, err
	}
	return sqrtPriceLowerTick, sqrtPriceUpperTick, nil
}

func TickToSqrtPrice(tickIndex int64, tickParams TickParams) (math.LegacyDec, error) {
	priceDec, err := TickToPrice(tickIndex, tickParams)
	if err != nil {
		return math.LegacyDec{}, err
	}

	sqrtPrice, err := priceDec.ApproxSqrt()
	if err != nil {
		return math.LegacyDec{}, err
	}
	return sqrtPrice, nil
}

func TickToPrice(tickIndex int64, tickParams TickParams) (math.LegacyDec, error) {
	offsetPrice := Pow(tickParams.PriceRatio, tickParams.BaseOffset)
	if tickIndex == 0 {
		return offsetPrice, nil
	}

	if tickIndex == TICK_MIN {
		return MinSpotPrice, nil
	}

	var price math.LegacyDec
	if tickIndex > 0 {
		price = Pow(tickParams.PriceRatio, math.LegacyNewDec(tickIndex))
	} else {
		price = math.LegacyOneDec().Quo(Pow(tickParams.PriceRatio, math.LegacyNewDec(-tickIndex)))
	}

	price = price.Mul(offsetPrice)

	if price.GT(MaxSpotPrice) || price.LT(MinSpotPrice) {
		return math.LegacyDec{}, ErrPriceOutOfBound
	}
	return price, nil
}

func CalculatePriceToTick(price math.LegacyDec, tickParams TickParams) (tickIndex int64, err error) {
	if price.IsNegative() {
		return 0, fmt.Errorf("price must be greater than zero")
	}
	if price.GT(MaxSpotPrice) || price.LT(MinSpotPrice) {
		return 0, ErrPriceOutOfBound
	}
	offsetPrice := Pow(tickParams.PriceRatio, tickParams.BaseOffset)
	if price.Equal(offsetPrice) {
		return 0, nil
	}

	tickIndex = 0
	if price.GT(offsetPrice) {
		for price.GT(offsetPrice) {
			price = price.Quo(tickParams.PriceRatio)
			tickIndex++
		}
	} else {
		for price.LT(offsetPrice) {
			price = price.Mul(tickParams.PriceRatio)
			tickIndex--
		}
	}

	return tickIndex, nil
}

func CalculateSqrtPriceToTick(sqrtPrice math.LegacyDec, tickParams TickParams) (tickIndex int64, err error) {
	price := sqrtPrice.Mul(sqrtPrice)

	tick, err := CalculatePriceToTick(price, tickParams)
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

	sqrtPriceTplus1, err := TickToSqrtPrice(tick+1, tickParams)
	if err != nil {
		return 0, ErrSqrtPriceToTick
	}

	if sqrtPrice.GTE(sqrtPriceTplus1) {
		sqrtPriceTplus2, err := TickToSqrtPrice(tick+2, tickParams)
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

	sqrtPriceT, err := TickToSqrtPrice(tick, tickParams)
	if err != nil {
		return 0, ErrSqrtPriceToTick
	}
	if sqrtPrice.GTE(sqrtPriceT) {
		return tick, nil
	}

	sqrtPriceTmin1, err := TickToSqrtPrice(tick-1, tickParams)
	if err != nil {
		return 0, ErrSqrtPriceToTick
	}
	if sqrtPrice.LT(sqrtPriceTmin1) {
		return 0, ErrSqrtPriceToTick
	}

	return tick - 1, nil
}
