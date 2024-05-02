package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgCreateRegistration{}

func NewMsgCreateRegistration(
	address string,
	proxyAddress string,

) *MsgCreateRegistration {
	return &MsgCreateRegistration{
		Address:      address,
		ProxyAddress: proxyAddress,
	}
}

func (msg *MsgCreateRegistration) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgUpdateRegistration{}

func NewMsgUpdateRegistration(
	address string,
	proxyAddress string,

) *MsgUpdateRegistration {
	return &MsgUpdateRegistration{
		Address:      address,
		ProxyAddress: proxyAddress,
	}
}

func (msg *MsgUpdateRegistration) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgDeleteRegistration{}

func NewMsgDeleteRegistration(
	address string,

) *MsgDeleteRegistration {
	return &MsgDeleteRegistration{
		Address: address,
	}
}

func (msg *MsgDeleteRegistration) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid address (%s)", err)
	}
	return nil
}
