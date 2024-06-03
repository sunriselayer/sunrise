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
func NewParams(feeDenom string, burnRatio math.LegacyDec, bypassDenoms []string) Params {
	return Params{
		FeeDenom:     feeDenom,
		BurnRatio:    burnRatio,
		BypassDenoms: bypassDenoms,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		"fee",
		math.LegacyMustNewDecFromStr("0.5"),
		[]string{"stake"},
	)
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if p.FeeDenom == "" {
		return ErrEmptyFeeDenom
	}

	if p.BurnRatio.IsNegative() || p.BurnRatio.GTE(math.LegacyOneDec()) {
		return ErrInvalidBurnRatio
	}

	for _, denom := range p.BypassDenoms {
		if denom == "" {
			return ErrEmptyBypassDenom
		}
	}

	return nil
}
