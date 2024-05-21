package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"cosmossdk.io/math"
)

var _ sdk.Msg = &MsgConvert{}

func NewMsgConvert(sender string, minAmount math.Int, maxAmount math.Int) *MsgConvert {
	return &MsgConvert{
		Sender:    sender,
		MinAmount: minAmount,
		MaxAmount: maxAmount,
	}
}

func (msg *MsgConvert) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if msg.MinAmount.IsNegative() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "min amount must not be negative")
	}

	if !msg.MaxAmount.IsPositive() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "max amount must be positive")
	}

	return nil
}
