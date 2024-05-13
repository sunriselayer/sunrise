package cfmm

import (
	"cosmossdk.io/math"
)

var _ ConstantFunctionMarketMaker = StableMarketMaker{}

func (smm StableMarketMaker) X(y math.LegacyDec, k math.LegacyDec) math.LegacyDec {

}

func (smm StableMarketMaker) Y(x math.LegacyDec, k math.LegacyDec) math.LegacyDec {

}

// k = xy (x^2 + y^2)
func (smm StableMarketMaker) K(x math.LegacyDec, y math.LegacyDec) math.LegacyDec {
	xy := x.Mul(y)
	x2 := x.Mul(x)
	y2 := y.Mul(y)
	k := xy.Mul(x2.Add(y2))

	return k
}
