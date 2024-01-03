package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgLiquidUnstake{}

func NewMsgLiquidUnstake(
	liquidStaker sdk.AccAddress,
	amount sdk.Coin,
) *MsgLiquidUnstake {
	return &MsgLiquidUnstake{
		DelegatorAddress: liquidStaker.String(),
		Amount:           amount,
	}
}

func (msg *MsgLiquidUnstake) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid delegator address (%s)", err)
	}
	if ok := msg.Amount.IsZero(); ok {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "unstaking amount must not be zero")
	}
	if err := msg.Amount.Validate(); err != nil {
		return err
	}
	return nil
}
