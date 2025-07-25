package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/app/consts"
)

// NewParams creates a new Params instance.
func NewParams(
	stableDenom string,
	AllowedAddresses []string,
) Params {
	if len(AllowedAddresses) == 0 {
		AllowedAddresses = nil
	}
	return Params{
		StableDenom:      stableDenom,
		AllowedAddresses: AllowedAddresses,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(consts.StableDenom, []string{})
}

// Validate validates the set of params.
func (p Params) Validate() error {
	if p.StableDenom == "" {
		return fmt.Errorf("stable denom cannot be empty")
	}
	if err := sdk.ValidateDenom(p.StableDenom); err != nil {
		return fmt.Errorf("invalid stable denom: %w", err)
	}

	for _, addr := range p.AllowedAddresses {
		if _, err := sdk.AccAddressFromBech32(addr); err != nil {
			return fmt.Errorf("invalid address: %w", err)
		}
	}

	return nil
}
