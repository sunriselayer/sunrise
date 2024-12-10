package types

import (
	"cosmossdk.io/math"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(interfaceFeeRate math.LegacyDec) Params {
	return Params{
		InterfaceFeeRate: interfaceFeeRate.String(),
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(math.LegacyNewDecWithPrec(1, 2)) // 1%
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{}
}

// Validate validates the set of params
func (p Params) Validate() error {
	_, err := math.LegacyNewDecFromStr(p.InterfaceFeeRate)
	if err != nil {
		return err
	}

	return nil
}
