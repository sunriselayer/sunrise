package types

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewParams creates a new Params instance.
func NewParams(bondDenom string, feeDenom string, selfDelegationCap math.Int) Params {
	return Params{
		BondDenom:         bondDenom,
		FeeDenom:          feeDenom,
		SelfDelegationCap: selfDelegationCap,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(
		sdk.DefaultBondDenom,
		"token",
		math.NewInt(1_000_000).Mul(math.NewInt(1_000_000)),
	)
}

// Validate validates the set of params.
func (p Params) Validate() error {
	if err := sdk.ValidateDenom(p.BondDenom); err != nil {
		return err
	}

	if err := sdk.ValidateDenom(p.FeeDenom); err != nil {
		return err
	}

	if p.SelfDelegationCap.IsNil() || !p.SelfDelegationCap.IsPositive() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "self delegation cap must be positive")
	}

	return nil
}
