package types

import "cosmossdk.io/math"

func NewMsgMint(admin string, tokenCreator string, amount math.Int) *MsgMint {
	return &MsgMint{
		Admin:        admin,
		TokenCreator: tokenCreator,
		Amount:       amount,
	}
}
