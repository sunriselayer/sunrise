package types

import sdk "github.com/cosmos/cosmos-sdk/types"

func NewMsgLiquidate(sender string, borrowId uint64, amount sdk.Coin) *MsgLiquidate {
	return &MsgLiquidate{
		Sender:   sender,
		BorrowId: borrowId,
		Amount:   amount,
	}
}
