package types

import (
	"cosmossdk.io/math"
)

func NewMsgBurn(sender string, amount math.Int, outputDenom string) *MsgBurn {
	return &MsgBurn{
		Sender:      sender,
		Amount:      amount,
		OutputDenom: outputDenom,
	}
}
