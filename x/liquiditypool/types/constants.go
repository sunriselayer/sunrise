package types

import (
	"cosmossdk.io/math"
)

var (
	MultiplierSqrt         = math.LegacyNewDec(1000_000_000)
	Multiplier             = MultiplierSqrt.Mul(MultiplierSqrt)
	MaxSqrtPrice           = math.LegacyMustNewDecFromStr("10000000000000000000") // 10^19
	MaxMultipliedSpotPrice = Multiplier.Mul(MaxSqrtPrice).Mul(MaxSqrtPrice)
	MinSqrtPrice           = math.LegacyNewDecWithPrec(1, 18)
	MinMultipliedSpotPrice = Multiplier.Mul(MinSqrtPrice).Mul(MinSqrtPrice)
)
