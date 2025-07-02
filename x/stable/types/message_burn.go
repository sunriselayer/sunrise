package types

import sdk "github.com/cosmos/cosmos-sdk/types"

func NewMsgBurn(sender string, amount sdk.Coins) *MsgBurn {
	return &MsgBurn{
		Sender: sender,
		Amount: amount,
	}
}
