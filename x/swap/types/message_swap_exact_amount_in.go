package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	sdkmath "cosmossdk.io/math"
)

var _ sdk.Msg = &MsgSwapExactAmountIn{}

func NewMsgSwapExactAmountIn(
	sender string,
	interfaceProvider string,
	route Route,
	amountIn sdkmath.Int,
	minAmountOut sdkmath.Int,
) *MsgSwapExactAmountIn {
	return &MsgSwapExactAmountIn{
		Sender:            sender,
		InterfaceProvider: interfaceProvider,
		Route:             route,
		AmountIn:          amountIn,
		MinAmountOut:      minAmountOut,
	}
}

func (msg *MsgSwapExactAmountIn) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if msg.InterfaceProvider != "" {
		if _, err := sdk.AccAddressFromBech32(msg.InterfaceProvider); err != nil {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid interface provider address (%s)", err)
		}
	}

	if err := msg.Route.Validate(); err != nil {
		return errorsmod.Wrapf(ErrInvalidRoute, "invalid route: %s", err)
	}

	if !msg.AmountIn.IsPositive() {
		return errorsmod.Wrapf(ErrInvalidAmount, "amount in must be positive: %s", msg.AmountIn)
	}

	if !msg.MinAmountOut.IsPositive() {
		return errorsmod.Wrapf(ErrInvalidAmount, "min amount out must be positive: %s", msg.MinAmountOut)
	}

	return nil
}
