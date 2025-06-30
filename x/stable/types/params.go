package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/app/consts"
)

// NewParams creates a new Params instance.
func NewParams(
	stableDenom string,
	authorityAddresses []string,
	acceptedDenoms []string,
) Params {
	if len(authorityAddresses) == 0 {
		authorityAddresses = nil
	}
	if len(acceptedDenoms) == 0 {
		acceptedDenoms = nil
	}
	return Params{
		StableDenom:        stableDenom,
		AuthorityAddresses: authorityAddresses,
		AcceptedDenoms:     acceptedDenoms,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(consts.StableDenom, []string{}, []string{})
}

// Validate validates the set of params.
func (p Params) Validate() error {
	if p.StableDenom == "" {
		return fmt.Errorf("stable denom cannot be empty")
	}
	if err := sdk.ValidateDenom(p.StableDenom); err != nil {
		return fmt.Errorf("invalid stable denom: %w", err)
	}

	for _, addr := range p.AuthorityAddresses {
		if _, err := sdk.AccAddressFromBech32(addr); err != nil {
			return fmt.Errorf("invalid authority address: %w", err)
		}
	}

	denomSet := make(map[string]struct{})
	for _, denom := range p.AcceptedDenoms {
		if err := sdk.ValidateDenom(denom); err != nil {
			return fmt.Errorf("invalid accepted denom: %w", err)
		}
		if _, exists := denomSet[denom]; exists {
			return fmt.Errorf("duplicate accepted denom: %s", denom)
		}
		denomSet[denom] = struct{}{}
	}

	return nil
}
