package types

import (
	"cosmossdk.io/math"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sunriselayer/sunrise/app/consts"
)

// NewParams creates a new Params instance.
func NewParams(feeDenom string, burnDenom string, burnRatio math.LegacyDec, burnPoolId uint64, burnEnabled bool) Params {
	return Params{
		FeeDenom:    feeDenom,
		BurnDenom:   burnDenom,
		BurnRatio:   burnRatio.String(),
		BurnPoolId:  burnPoolId,
		BurnEnabled: burnEnabled,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(
		consts.StableDenom,
		consts.MintDenom,
		math.LegacyNewDecWithPrec(50, 2),
		0,
		false,
	)
}

// Validate validates the set of params.
func (p Params) Validate() error {
	if err := sdk.ValidateDenom(p.FeeDenom); err != nil {
		return err
	}
	if err := sdk.ValidateDenom(p.BurnDenom); err != nil {
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
