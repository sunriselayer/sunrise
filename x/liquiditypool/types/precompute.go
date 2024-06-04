package types

import (
	"cosmossdk.io/math"
)

var (
	MaxSpotPrice    = math.LegacyMustNewDecFromStr("100000000000000000000000000000000000000")
	MinSpotPrice    = math.LegacyNewDecWithPrec(1, 18)
	MaxSqrtPrice, _ = MaxSpotPrice.ApproxSqrt()
	MinSqrtPrice, _ = MinSpotPrice.ApproxSqrt()
)
