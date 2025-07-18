package types

import (
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/app/consts"
)

// NewParams creates a new Params instance.
func NewParams(nonTransferableDenom string, transferableDenom string, allowedAddresses []string) Params {
	if len(allowedAddresses) == 0 {
		allowedAddresses = nil
	}
	return Params{
		NonTransferableDenom: nonTransferableDenom,
		TransferableDenom:    transferableDenom,
		AllowedAddresses:     allowedAddresses,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(consts.BondDenom, consts.MintDenom, []string{})
}

// Validate validates the set of params.
func (p Params) Validate() error {
	if err := sdk.ValidateDenom(p.NonTransferableDenom); err != nil {
		return err
	}
	if err := sdk.ValidateDenom(p.TransferableDenom); err != nil {
		return err
	}

	for _, addr := range p.AllowedAddresses {
		if _, err := sdk.AccAddressFromBech32(addr); err != nil {
			return fmt.Errorf("invalid address: %w", err)
		}
	}
	return nil
}
