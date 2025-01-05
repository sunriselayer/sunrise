package types

import (
	"cosmossdk.io/math"
)

// NewParams creates a new Params instance.
func NewParams(interfaceFeeRate math.LegacyDec) Params {
	return Params{
		InterfaceFeeRate: interfaceFeeRate,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(math.LegacyNewDecWithPrec(1, 2)) // 1%
}

// Validate validates the set of params.
func (p Params) Validate() error {

	return nil
}
