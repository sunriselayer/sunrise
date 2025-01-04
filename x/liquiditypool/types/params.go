package types

import (
	"cosmossdk.io/math"
)

// NewParams creates a new Params instance.
func NewParams(withdrawFeeRate, swapTreasuryTaxRate math.LegacyDec) Params {
	return Params{
		WithdrawFeeRate:     withdrawFeeRate,
		SwapTreasuryTaxRate: swapTreasuryTaxRate,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(
		math.LegacyNewDecWithPrec(1, 2),
		math.LegacyNewDecWithPrec(1, 2),
	)
}

// Validate validates the set of params.
func (p Params) Validate() error {

	return nil
}
