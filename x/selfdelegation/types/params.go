package types

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewParams creates a new Params instance.
func NewParams(selfDelegationCap math.Int) Params {
	return Params{
		SelfDelegationCap: selfDelegationCap,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(
		math.NewInt(1_000_000).Mul(math.NewInt(1_000_000)),
	)
}

// Validate validates the set of params.
func (p Params) Validate() error {
	if p.SelfDelegationCap.IsNil() || !p.SelfDelegationCap.IsPositive() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "self delegation cap must be positive")
	}

	return nil
}
