package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgCollectIncentives{}

func NewMsgCollectIncentives(sender string, positionIds []uint64) *MsgCollectIncentives {
	return &MsgCollectIncentives{
		Sender:      sender,
		PositionIds: positionIds,
	}
}

func (msg *MsgCollectIncentives) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if len(msg.PositionIds) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "position ids cannot be empty")
	}

	return nil
}
