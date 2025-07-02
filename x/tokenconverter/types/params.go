package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/app/consts"
)

// NewParams creates a new Params instance.
func NewParams(fromDenom string, toDenom string) Params {
	return Params{
		FromDenom: fromDenom,
		ToDenom:   toDenom,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(consts.BondDenom, consts.MintDenom)
}

// Validate validates the set of params.
func (p Params) Validate() error {
	if err := sdk.ValidateDenom(p.FromDenom); err != nil {
		return err
	}
	if err := sdk.ValidateDenom(p.ToDenom); err != nil {
		return err
	}
	return nil
}
