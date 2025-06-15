package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewMsgMint(admin string, tokenCreator string, ybt sdk.Coin) *MsgMint {
	return &MsgMint{
		Admin:        admin,
		TokenCreator: tokenCreator,
		Ybt:          ybt,
	}
}
