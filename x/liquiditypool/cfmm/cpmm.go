package cfmm

import (
	"cosmossdk.io/math"
)

var _ ConcentratedLiquidityMarketMaker = ConstantProductMarketMaker{}

func (cpmm ConstantProductMarketMaker) X(y math.LegacyDec, k math.LegacyDec) math.LegacyDec {
	return k.Quo(y)
}

func (cpmm ConstantProductMarketMaker) Y(x math.LegacyDec, k math.LegacyDec) math.LegacyDec {
	return k.Quo(x)
}

func (cpmm ConstantProductMarketMaker) K(x math.LegacyDec, y math.LegacyDec) math.LegacyDec {
	return x.Mul(y)
}

// (y + y_v) / x_v = p_max
// y_v / (x + x_v) = p_min
//
// x_v = (p_min x + y) / (p_max - p_min)
// y_v = p_min (x + x_v)
func (cpmm ConstantProductMarketMaker) VirtualLiquidity(a1Range math.Int, a2Range math.Int, priceRange PriceRange) (*math.LegacyDec, *math.LegacyDec, error) {
	if err := priceRange.Validate(); err != nil {
		return nil, nil, err
	}

	x := math.LegacyNewDecFromInt(a1Range)
	y := math.LegacyNewDecFromInt(a2Range)

	priceSub := priceRange.MaxPrice.Sub(priceRange.MinPrice)
	x_v := priceRange.MinPrice.Mul(x).Add(y).Quo(priceSub)
	y_v := priceRange.MinPrice.Mul(x.Add(x_v))

	return &x_v, &y_v, nil
}
