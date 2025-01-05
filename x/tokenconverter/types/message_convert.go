package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"cosmossdk.io/math"
)

var _ sdk.Msg = &MsgConvert{}

func NewMsgConvert(sender string, amount math.Int) *MsgConvert {
	return &MsgConvert{
		Sender: sender,
		Amount: amount,
	}
}

func (msg *MsgConvert) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if !msg.Amount.IsPositive() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "amount must be positive")
	}

	return nil
}
