package types

import (
	"github.com/sunriselayer/sunrise/app/consts"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewParams creates a new Params instance.
func NewParams(
	authorityContract string,
	stableDenom string,
) Params {
	return Params{
		AuthorityContract: authorityContract,
		StableDenom:       stableDenom,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams("", consts.StableDenom)
}

// Validate validates the set of params.
func (p Params) Validate() error {
	err := sdk.ValidateDenom(p.StableDenom)
	if err != nil {
		return err
	}

	return nil
}
