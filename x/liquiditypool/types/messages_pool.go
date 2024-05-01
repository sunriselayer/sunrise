package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgCreatePool{}

func NewMsgCreatePool(creator string, baseDenom string, quoteDenom string) *MsgCreatePool {
	return &MsgCreatePool{
		Creator:    creator,
		BaseDenom:  baseDenom,
		QuoteDenom: quoteDenom,
	}
}

func (msg *MsgCreatePool) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgUpdatePool{}

func NewMsgUpdatePool(creator string, id uint64, baseDenom string, quoteDenom string) *MsgUpdatePool {
	return &MsgUpdatePool{
		Id:         id,
		Creator:    creator,
		BaseDenom:  baseDenom,
		QuoteDenom: quoteDenom,
	}
}

func (msg *MsgUpdatePool) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
