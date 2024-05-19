package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgCreatePosition{}

func NewMsgCreatePosition(sender string) *MsgCreatePosition {
	return &MsgCreatePosition{
		Sender: sender,
	}
}

func (msg *MsgCreatePosition) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgUpdatePosition{}

func NewMsgUpdatePosition(sender string, id uint64) *MsgUpdatePosition {
	return &MsgUpdatePosition{
		Id:     id,
		Sender: sender,
	}
}

func (msg *MsgUpdatePosition) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgDeletePosition{}

func NewMsgDeletePosition(sender string, id uint64) *MsgDeletePosition {
	return &MsgDeletePosition{
		Id:     id,
		Sender: sender,
	}
}

func (msg *MsgDeletePosition) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	return nil
}
