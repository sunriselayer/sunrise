package types

import (
	"cosmossdk.io/math"
)

func NewMsgMint(sender string, amount math.Int) *MsgMint {
	return &MsgMint{
		Sender: sender,
		Amount: amount,
	}
}
