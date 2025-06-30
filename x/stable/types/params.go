package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewParams creates a new Params instance.
func NewParams(
	authorityAddresses []string,
	acceptedDenoms []string,
	stableDenom string,
) Params {
	return Params{
		AuthorityAddresses: authorityAddresses,
		AcceptedDenoms:     acceptedDenoms,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams([]string{}, []string{}, "")
}

// Validate validates the set of params.
func (p Params) Validate() error {
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
