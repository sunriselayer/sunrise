package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/app/consts"
)

// NewParams creates a new Params instance.
func NewParams(nonTransferableDenom string, transferableDenom string) Params {
	return Params{
		NonTransferableDenom: nonTransferableDenom,
		TransferableDenom:    transferableDenom,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(consts.BondDenom, consts.MintDenom)
}

// Validate validates the set of params.
func (p Params) Validate() error {
	if err := sdk.ValidateDenom(p.NonTransferableDenom); err != nil {
		return err
	}
	if err := sdk.ValidateDenom(p.TransferableDenom); err != nil {
		return err
	}
	return nil
}
