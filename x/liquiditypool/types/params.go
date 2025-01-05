package types

import (
	"cosmossdk.io/math"

	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewParams creates a new Params instance.
func NewParams(withdrawFeeRate, swapTreasuryTaxRate math.LegacyDec) Params {
	return Params{
		WithdrawFeeRate:     withdrawFeeRate,
		SwapTreasuryTaxRate: swapTreasuryTaxRate,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(
		math.LegacyNewDecWithPrec(1, 2),
		math.LegacyNewDecWithPrec(1, 2),
	)
}

// Validate validates the set of params.
func (p Params) Validate() error {
	if p.WithdrawFeeRate.IsNegative() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "withdraw fee rate must not be negative")
	}
	if p.WithdrawFeeRate.GT(math.LegacyOneDec()) {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "withdraw fee rate must be less than 1")
	}

	if p.SwapTreasuryTaxRate.IsNegative() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "swap treasury tax rate must not be negative")
	}
	if p.SwapTreasuryTaxRate.GT(math.LegacyOneDec()) {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "swap treasury tax rate must be less than 1")
	}

	return nil
}
