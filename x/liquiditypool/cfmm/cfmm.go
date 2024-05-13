package cfmm

import (
	"cosmossdk.io/math"
)

type ConstantFunctionMarketMaker interface {
	X(y math.LegacyDec, k math.LegacyDec) math.LegacyDec
	Y(x math.LegacyDec, k math.LegacyDec) math.LegacyDec
	K(x math.LegacyDec, y math.LegacyDec) math.LegacyDec
}

type PriceRange struct {
	MinPrice math.LegacyDec
	MaxPrice math.LegacyDec
}

func NewPriceRange(minPrice math.LegacyDec, maxPrice math.LegacyDec) PriceRange {
	return PriceRange{
		MinPrice: minPrice,
		MaxPrice: maxPrice,
	}
}

func (pr PriceRange) Validate() error {
	if pr.MinPrice.IsZero() || pr.MinPrice.IsNegative() {
		return ErrInvalidMinPrice
	}

	if pr.MinPrice.GT(pr.MaxPrice) {
		return ErrMinPriceGreaterThanMaxPrice
	}

	return nil
}

type ConcentratedLiquidityMarketMaker interface {
	ConstantFunctionMarketMaker

	VirtualLiquidity(xRange math.Int, yRange math.Int, priceRange PriceRange) (*math.LegacyDec, *math.LegacyDec, error)
}
