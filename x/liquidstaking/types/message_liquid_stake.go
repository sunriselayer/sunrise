package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgLiquidStake{}

func NewMsgLiquidStake(creator string) *MsgLiquidStake {
	return &MsgLiquidStake{
		Creator: creator,
	}
}

func (msg *MsgLiquidStake) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if ok := msg.Amount.IsZero(); ok {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "staking amount must not be zero")
	}
	if err := msg.Amount.Validate(); err != nil {
		return err
	}
	return nil
}
