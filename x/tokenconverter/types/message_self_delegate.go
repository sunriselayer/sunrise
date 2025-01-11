package types

import (
	"cosmossdk.io/math"
)

func NewMsgSelfDelegate(sender string, amount math.Int) *MsgSelfDelegate {
	return &MsgSelfDelegate{
		Sender: sender,
		Amount: amount,
	}
}
