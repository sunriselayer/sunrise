package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgCreateRegistration{}

func NewMsgCreateRegistration(liquidityProvider string, grantee string) *MsgCreateRegistration {
	return &MsgCreateRegistration{
		LiquidityProvider: liquidityProvider,
		Grantee:           grantee,
	}
}

func (msg *MsgCreateRegistration) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.LiquidityProvider)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid address (%s)", err)
	}
	_, err = sdk.AccAddressFromBech32(msg.Grantee)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgUpdateRegistration{}

func NewMsgUpdateRegistration(address string, grantee string) *MsgUpdateRegistration {
	return &MsgUpdateRegistration{
		LiquidityProvider: address,
		Grantee:           grantee,
	}
}

func (msg *MsgUpdateRegistration) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.LiquidityProvider)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid address (%s)", err)
	}
	_, err = sdk.AccAddressFromBech32(msg.Grantee)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgDeleteRegistration{}

func NewMsgDeleteRegistration(liquidityProvider string) *MsgDeleteRegistration {
	return &MsgDeleteRegistration{
		LiquidityProvider: liquidityProvider,
	}
}

func (msg *MsgDeleteRegistration) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.LiquidityProvider)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid address (%s)", err)
	}
	return nil
}
