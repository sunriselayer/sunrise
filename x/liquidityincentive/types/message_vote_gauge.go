package types

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgVoteGauge{}

func NewMsgVoteGauge(sender string) *MsgVoteGauge {
	return &MsgVoteGauge{
		Sender: sender,
	}
}

func (msg *MsgVoteGauge) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	totalWeight := math.LegacyZeroDec()
	for _, poolWeight := range msg.PoolWeights {
		weight, _ := math.LegacyNewDecFromStr(poolWeight.Weight)
		totalWeight = totalWeight.Add(weight)
	}
	if totalWeight.GT(math.LegacyOneDec()) {
		return errorsmod.Wrapf(ErrTotalWeightGTOne, "total weight: %s", totalWeight.String())
	}
	return nil
}
