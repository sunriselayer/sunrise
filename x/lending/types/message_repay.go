package types

import sdk "github.com/cosmos/cosmos-sdk/types"

func NewMsgRepay(sender string, borrowId uint64, amount sdk.Coin) *MsgRepay {
	return &MsgRepay{
		Sender:   sender,
		BorrowId: borrowId,
		Amount:   amount,
	}
}
