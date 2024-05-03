package types

import (
	"cosmossdk.io/math"
)

type ConstantFunctionMarketMaker interface {
	x(y math.Int, k math.LegacyDec) math.Int
	y(x math.Int, k math.LegacyDec) math.Int
	k(x math.Int, y math.Int) math.LegacyDec
}
