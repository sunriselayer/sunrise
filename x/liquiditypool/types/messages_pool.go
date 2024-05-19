package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgCreatePool{}

func NewMsgCreatePool(sender string, lowerTick string, upperTick string) *MsgCreatePool {
	return &MsgCreatePool{
		Sender:    sender,
		LowerTick: lowerTick,
		UpperTick: upperTick,
	}
}

func (msg *MsgCreatePool) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgUpdatePool{}

func NewMsgUpdatePool(sender string, id uint64, lowerTick string, upperTick string) *MsgUpdatePool {
	return &MsgUpdatePool{
		Id:        id,
		Sender:    sender,
		LowerTick: lowerTick,
		UpperTick: upperTick,
	}
}

func (msg *MsgUpdatePool) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgDeletePool{}

func NewMsgDeletePool(sender string, id uint64) *MsgDeletePool {
	return &MsgDeletePool{
		Id:     id,
		Sender: sender,
	}
}

func (msg *MsgDeletePool) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	return nil
}
