package types

import (
	"cosmossdk.io/math"
)

func NewMsgConvert(sender string, minAmount math.Int, maxAmount math.Int) *MsgConvert {
	return &MsgConvert{
		Sender:    sender,
		MinAmount: minAmount,
		MaxAmount: maxAmount,
	}
}
