package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgCreateRegistration{}

func NewMsgCreateRegistration(
    creator string,
    address string,
    proxyAddress string,
    
) *MsgCreateRegistration {
  return &MsgCreateRegistration{
		Creator : creator,
		Address: address,
		ProxyAddress: proxyAddress,
        
	}
}

func (msg *MsgCreateRegistration) ValidateBasic() error {
  _, err := sdk.AccAddressFromBech32(msg.Creator)
  	if err != nil {
  		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
  	}
  return nil
}

var _ sdk.Msg = &MsgUpdateRegistration{}

func NewMsgUpdateRegistration(
    creator string,
    address string,
    proxyAddress string,
    
) *MsgUpdateRegistration {
  return &MsgUpdateRegistration{
		Creator: creator,
        Address: address,
        ProxyAddress: proxyAddress,
        
	}
}

func (msg *MsgUpdateRegistration) ValidateBasic() error {
  _, err := sdk.AccAddressFromBech32(msg.Creator)
  if err != nil {
    return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
  }
   return nil
}

var _ sdk.Msg = &MsgDeleteRegistration{}

func NewMsgDeleteRegistration(
    creator string,
    address string,
    
) *MsgDeleteRegistration {
  return &MsgDeleteRegistration{
		Creator: creator,
		Address: address,
        
	}
}

func (msg *MsgDeleteRegistration) ValidateBasic() error {
  _, err := sdk.AccAddressFromBech32(msg.Creator)
  if err != nil {
    return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
  }
  return nil
}
