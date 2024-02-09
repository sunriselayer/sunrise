package types

import (
	"cosmossdk.io/errors"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for blobstream module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams() Params {
	return Params{}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams()
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(ParamsStoreKeyDataCommitmentWindow, &p.DataCommitmentWindow, validateDataCommitmentWindow),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateDataCommitmentWindow(p.DataCommitmentWindow); err != nil {
		return errors.Wrap(err, "data commitment window")
	}
	return nil
}
