package types

import (
	"cosmossdk.io/math"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sunriselayer/sunrise/app/consts"
)

// NewParams creates a new Params instance.
func NewParams(createPoolGas uint64, withdrawFeeRate math.LegacyDec, allowedQuoteDenoms []string, AllowedAddresses []string) Params {
	if len(AllowedAddresses) == 0 {
		AllowedAddresses = nil
	}
	return Params{
		CreatePoolGas:      createPoolGas,
		WithdrawFeeRate:    withdrawFeeRate.String(),
		AllowedQuoteDenoms: allowedQuoteDenoms,
		AllowedAddresses:   AllowedAddresses,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(
		1000000,
		math.LegacyNewDecWithPrec(1, 2),
		[]string{consts.MintDenom, consts.StableDenom},
		[]string{},
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

	for _, denom := range p.AllowedQuoteDenoms {
		if err := sdk.ValidateDenom(denom); err != nil {
			return errorsmod.Wrapf(err, "invalid allowed quote denom: %s", denom)
		}
	}

	for _, addr := range p.AllowedAddresses {
		if _, err := sdk.AccAddressFromBech32(addr); err != nil {
			return errorsmod.Wrapf(err, "invalid address: %s", addr)
		}
	}

	return nil
}
