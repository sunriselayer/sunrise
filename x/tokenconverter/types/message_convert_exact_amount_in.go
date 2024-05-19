package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"cosmossdk.io/math"
)

var _ sdk.Msg = &MsgConvertExactAmountIn{}

func NewMsgConvertExactAmountIn(sender string, amountIn math.Int, minAmountOut math.Int) *MsgConvertExactAmountIn {
	return &MsgConvertExactAmountIn{
		Sender:       sender,
		AmountIn:     amountIn,
		MinAmountOut: minAmountOut,
	}
}

func (msg *MsgConvertExactAmountIn) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if !msg.AmountIn.IsPositive() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "amount in must be positive")
	}

	if !msg.MinAmountOut.IsPositive() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "min amount out must be positive")
	}

	return nil
}
