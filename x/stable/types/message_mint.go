package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewMsgMint(sender string, amount sdk.Coins) *MsgMint {
	return &MsgMint{
		Sender: sender,
		Amount: amount,
	}
}
