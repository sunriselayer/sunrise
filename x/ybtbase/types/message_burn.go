package types

import "cosmossdk.io/math"

func NewMsgBurn(admin string, tokenCreator string, amount math.Int) *MsgBurn {
	return &MsgBurn{
		Admin:        admin,
		TokenCreator: tokenCreator,
		Amount:       amount,
	}
}
