package types

import (
	"cosmossdk.io/math"
)

func NewMsgConvert(sender string, amount math.Int) *MsgConvert {
	return &MsgConvert{
		Sender: sender,
		Amount: amount,
	}
}
