package types

import sdk "github.com/cosmos/cosmos-sdk/types"

func NewMsgBurn(admin string, tokenCreator string, amount sdk.Coin) *MsgBurn {
	return &MsgBurn{
		Admin:        admin,
		TokenCreator: tokenCreator,
		Amount:       amount,
	}
}
