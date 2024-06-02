package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgSwapExactAmountIn{}

func NewMsgSwapExactAmountIn(sender string) *MsgSwapExactAmountIn {
	return &MsgSwapExactAmountIn{
		Sender: sender,
	}
}

func (msg *MsgSwapExactAmountIn) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	if err := msg.Route.Validate(); err != nil {
		return errorsmod.Wrapf(ErrInvalidRoute, "invalid route: %s", err)
	}

	return nil
}
