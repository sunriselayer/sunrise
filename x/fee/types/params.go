package types

import (
	"cosmossdk.io/math"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewParams creates a new Params instance.
func NewParams(feeDenom string, burnRatio math.LegacyDec, bypassDenoms []string) Params {
	return Params{
		FeeDenom:     feeDenom,
		BurnRatio:    burnRatio,
		BypassDenoms: bypassDenoms,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(
		"fee",
		math.LegacyNewDecWithPrec(50, 2),
		[]string{"stake"},
	)
}

// Validate validates the set of params.
func (p Params) Validate() error {
	if err := sdk.ValidateDenom(p.FeeDenom); err != nil {
		return err
	}

	if p.BurnRatio.IsNegative() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "burn ratio must not be negative")
	}
	if p.BurnRatio.GT(math.LegacyOneDec()) {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "burn ratio must be less than 1")
	}

	for _, bypassDenom := range p.BypassDenoms {
		if err := sdk.ValidateDenom(bypassDenom); err != nil {
			return err
		}
	}

	return nil
}
