package types

import (
	"cosmossdk.io/math"
)

func NewMsgSelfUndelegate(sender string, amount math.Int) *MsgSelfUndelegate {
	return &MsgSelfUndelegate{
		Sender: sender,
		Amount: amount,
	}
}
