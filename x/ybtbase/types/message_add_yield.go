package types

import "cosmossdk.io/math"

func NewMsgAddYield(admin string, tokenCreator string, amount math.Int) *MsgAddYield {
	return &MsgAddYield{
		Admin:        admin,
		TokenCreator: tokenCreator,
		Amount:       amount,
	}
}
