package types

import (
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	DefaultSecurityAddress = ""
	DefaultLimit           = uint64(5)
)

// NewParams creates a new Params instance.
func NewParams(securityAddress string, limit uint64) Params {
	return Params{
		SecurityAddress: securityAddress,
		Limit:           limit,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(DefaultSecurityAddress, DefaultLimit)
}

// Validate validates the set of params.
func (p Params) Validate() error {
	err := validateAddress(p.SecurityAddress)
	if err != nil {
		return fmt.Errorf("invalid security address: %w", err)
	}

	err = validateLimit(p.Limit)
	if err != nil {
		return fmt.Errorf("invalid limit: %w", err)
	}

	return nil
}

func validateAddress(address string) error {
	// address might be explicitly empty in test environments
	if len(address) == 0 {
		return nil
	}

	_, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return fmt.Errorf("invalid address: %w", err)
	}

	return nil
}

func validateLimit(limit uint64) error {
	if limit == 0 {
		return fmt.Errorf("limit cannot be zero")
	}

	return nil
}
