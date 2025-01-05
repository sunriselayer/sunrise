package types

import (
	"cosmossdk.io/math"

	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewParams creates a new Params instance.
func NewParams(interfaceFeeRate math.LegacyDec) Params {
	return Params{
		InterfaceFeeRate: interfaceFeeRate,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(math.LegacyNewDecWithPrec(1, 2)) // 1%
}

// Validate validates the set of params.
func (p Params) Validate() error {
	if p.InterfaceFeeRate.IsNegative() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "interface fee rate must not be negative")
	}
	if p.InterfaceFeeRate.GT(math.LegacyOneDec()) {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "interface fee rate must be less than 1")
	}

	return nil
}
