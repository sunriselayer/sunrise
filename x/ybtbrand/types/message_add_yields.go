package types

import sdk "github.com/cosmos/cosmos-sdk/types"

func NewMsgAddYields(admin string, tokenCreator string, amount sdk.Coins) *MsgAddYields {
	return &MsgAddYields{
		Admin:        admin,
		TokenCreator: tokenCreator,
		Amount:       amount,
	}
}
