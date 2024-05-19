package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"cosmossdk.io/math"
)

var _ sdk.Msg = &MsgConvertExactAmountOut{}

func NewMsgConvertExactAmountOut(sender string, amountOut math.Int, maxAmountIn math.Int) *MsgConvertExactAmountOut {
	return &MsgConvertExactAmountOut{
		Sender:      sender,
		AmountOut:   amountOut,
		MaxAmountIn: maxAmountIn,
	}
}

func (msg *MsgConvertExactAmountOut) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if !msg.AmountOut.IsPositive() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "amount out must be positive")
	}

	if !msg.MaxAmountIn.IsPositive() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "max amount in must be positive")
	}

	return nil
}
