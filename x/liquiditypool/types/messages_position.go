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

var _ sdk.Msg = &MsgIncreaseLiquidity{}

func NewMsgIncreaseLiquidity(sender string, id uint64) *MsgIncreaseLiquidity {
	return &MsgIncreaseLiquidity{
		Id:     id,
		Sender: sender,
	}
}

func (msg *MsgIncreaseLiquidity) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgDecreaseLiquidity{}

func NewMsgDecreaseLiquidity(sender string, id uint64) *MsgDecreaseLiquidity {
	return &MsgDecreaseLiquidity{
		Id:     id,
		Sender: sender,
	}
}

func (msg *MsgDecreaseLiquidity) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	return nil
}
