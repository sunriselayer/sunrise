package types

import sdk "github.com/cosmos/cosmos-sdk/types"

func NewMsgSupply(sender string, amount sdk.Coin) *MsgSupply {
	return &MsgSupply{
		Sender: sender,
		Amount: amount,
	}
}
