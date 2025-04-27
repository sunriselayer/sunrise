package types

import (
	"cosmossdk.io/math"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewParams creates a new Params instance.
func NewParams(feeDenom string, burnRatio math.LegacyDec) Params {
	return Params{
		FeeDenom:  feeDenom,
		BurnRatio: burnRatio.String(),
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(
		"token",
		math.LegacyNewDecWithPrec(50, 2),
	)
}

// Validate validates the set of params.
func (p Params) Validate() error {
	if err := sdk.ValidateDenom(p.FeeDenom); err != nil {
		return err
	}

	burnRatio, err := math.LegacyNewDecFromStr(p.BurnRatio)
	if err != nil {
		return err
	}
	if burnRatio.IsNegative() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "burn ratio must not be negative")
	}
	if burnRatio.GT(math.LegacyOneDec()) {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "burn ratio must be less than 1")
	}

	return nil
}
