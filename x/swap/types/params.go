package types

import (
	"cosmossdk.io/math"

	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewParams creates a new Params instance.
func NewParams(interfaceFeeRate math.LegacyDec) Params {
	return Params{
		InterfaceFeeRate: interfaceFeeRate.String(),
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(math.LegacyNewDecWithPrec(1, 2)) // 1%
}

// Validate validates the set of params.
func (p Params) Validate() error {
	interfaceFeeRate, err := math.LegacyNewDecFromStr(p.InterfaceFeeRate)
	if err != nil {
		return err
	}
	if interfaceFeeRate.IsNegative() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "interface fee rate must not be negative")
	}
	if interfaceFeeRate.GT(math.LegacyOneDec()) {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "interface fee rate must be less than 1")
	}

	return nil
}
