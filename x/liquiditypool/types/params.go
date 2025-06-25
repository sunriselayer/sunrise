package types

import (
	"cosmossdk.io/math"

	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewParams creates a new Params instance.
func NewParams(createPoolGas uint64, withdrawFeeRate math.LegacyDec) Params {
	return Params{
		CreatePoolGas:   createPoolGas,
		WithdrawFeeRate: withdrawFeeRate.String(),
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(
		1000000,
		math.LegacyNewDecWithPrec(1, 2),
	)
}

// Validate validates the set of params.
func (p Params) Validate() error {
	withdrawFeeRate, err := math.LegacyNewDecFromStr(p.WithdrawFeeRate)
	if err != nil {
		return err
	}
	if withdrawFeeRate.IsNegative() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "withdraw fee rate must not be negative")
	}
	if withdrawFeeRate.GT(math.LegacyOneDec()) {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "withdraw fee rate must be less than 1")
	}

	return nil
}
