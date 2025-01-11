package types

import (
	"cosmossdk.io/math"
)

func NewMsgWithdrawSelfUnbonded(sender string, amount math.Int) *MsgWithdrawSelfUnbonded {
	return &MsgWithdrawSelfUnbonded{
		Sender: sender,
	}
}
