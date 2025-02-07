package types

import (
	"errors"
	"fmt"

	"cosmossdk.io/math"
)

func CheckTicks(tickLower, tickUpper int64) error {
	if tickLower >= tickUpper {
		return errors.New("tickLower should be less than tickUpper")
	}
	if tickLower < TICK_MIN {
		return errors.New("tickLower is out of range")
	}

	if tickUpper > TICK_MAX {
		return errors.New("tickUpper is out of range")
	}
	return nil
}

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
	priceMultiplied, err := TickToMultipliedPrice(tickIndex, tickParams)
	if err != nil {
		return math.LegacyDec{}, err
	}

	sqrtPriceMultiplied, err := priceMultiplied.ApproxSqrt()
	if err != nil {
		return math.LegacyDec{}, err
	}
	sqrtPrice := sqrtPriceMultiplied.Quo(MultiplierSqrt)
	return sqrtPrice, nil
}

func TickToMultipliedPrice(tickIndex int64, tickParams TickParams) (math.LegacyDec, error) {
	offsetPrice := Pow(tickParams.PriceRatio, tickParams.BaseOffset)
	if tickIndex == 0 {
		return offsetPrice.Mul(Multiplier), nil
	}

	if tickIndex == TICK_MIN {
		return MinMultipliedSpotPrice, nil
	}

	var multipliedPrice math.LegacyDec
	if tickIndex > 0 {
		multipliedPrice = Multiplier.Mul(Pow(tickParams.PriceRatio, math.LegacyNewDec(tickIndex)))
	} else {
		multipliedPrice = Multiplier.Quo(Pow(tickParams.PriceRatio, math.LegacyNewDec(-tickIndex)))
	}

	multipliedPrice = multipliedPrice.Mul(offsetPrice)

	if multipliedPrice.GT(MaxMultipliedSpotPrice) || multipliedPrice.LT(MinMultipliedSpotPrice) {
		return math.LegacyDec{}, ErrPriceOutOfBound
	}
	return multipliedPrice, nil
}

func CalculateMultipliedPriceToTick(multipliedPrice math.LegacyDec, tickParams TickParams) (tickIndex int64, err error) {
	if multipliedPrice.IsNegative() {
		return 0, fmt.Errorf("price must be greater than zero")
	}
	if multipliedPrice.GT(MaxMultipliedSpotPrice) || multipliedPrice.LT(MinMultipliedSpotPrice) {
		return 0, ErrPriceOutOfBound
	}
	multipliedOffsetPrice := Multiplier.Mul(Pow(tickParams.PriceRatio, tickParams.BaseOffset))
	if multipliedPrice.Equal(multipliedOffsetPrice) {
		return 0, nil
	}

	tickIndex = 0
	if multipliedPrice.GT(multipliedOffsetPrice) {
		for multipliedPrice.GT(multipliedOffsetPrice) {
			multipliedPrice = multipliedPrice.Quo(tickParams.PriceRatio)
			tickIndex++
		}
	} else {
		for multipliedPrice.LT(multipliedOffsetPrice) {
			multipliedPrice = multipliedPrice.Mul(tickParams.PriceRatio)
			tickIndex--
		}
	}

	return tickIndex, nil
}

func CalculateSqrtPriceToTick(sqrtPrice math.LegacyDec, tickParams TickParams) (tickIndex int64, err error) {
	multipliedPrice := Multiplier.Mul(sqrtPrice).Mul(sqrtPrice)
	tick, err := CalculateMultipliedPriceToTick(multipliedPrice, tickParams)
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

func GetSqrtPriceFromQuoteBase(quoteAmount math.Int, baseAmount math.Int) (math.LegacyDec, error) {
	spotPriceMultiplied := quoteAmount.ToLegacyDec().Mul(Multiplier).Quo(baseAmount.ToLegacyDec())
	sqrtPriceMultiplied, err := spotPriceMultiplied.ApproxSqrt()
	if err != nil {
		return math.LegacyZeroDec(), err
	}
	sqrtPrice := sqrtPriceMultiplied.Quo(MultiplierSqrt)
	return sqrtPrice, err
}
