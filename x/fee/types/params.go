package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewParams creates a new Params instance.
func NewParams(feeDenom string, burnRatio math.LegacyDec, bypassDenoms []string) Params {
	return Params{
		FeeDenom:     feeDenom,
		BurnRatio:    burnRatio,
		BypassDenoms: bypassDenoms,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(
		"token",
		math.LegacyNewDecWithPrec(50, 2),
		[]string{sdk.DefaultBondDenom},
	)
}

// Validate validates the set of params.
func (p Params) Validate() error {

	return nil
}
