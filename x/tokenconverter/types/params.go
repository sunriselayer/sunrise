package types

import (
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"cosmossdk.io/math"
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(govDenom string, feeDenom string, maxSupplyFee math.Int) Params {
	return Params{
		GovDenom:     govDenom,
		FeeDenom:     feeDenom,
		MaxSupplyFee: maxSupplyFee,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		"stake",
		"fee",
		math.NewInt(1000_000_000_000_000),
	)
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if p.GovDenom == "" {
		return ErrEmptyGovDenom
	}

	if p.FeeDenom == "" {
		return ErrEmptyFeeDenom
	}

	if p.MaxSupplyFee.IsNegative() {
		return ErrNegativeMaxSupplyFee
	}

	return nil
}
