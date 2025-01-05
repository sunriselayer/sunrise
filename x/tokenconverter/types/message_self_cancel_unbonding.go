package types

import (
	"cosmossdk.io/math"
)

func NewMsgSelfCancelUnbonding(sender string, amount math.Int, creationHeight int64) *MsgSelfCancelUnbonding {
	return &MsgSelfCancelUnbonding{
		Sender:         sender,
		Amount:         amount,
		CreationHeight: creationHeight,
	}
}
