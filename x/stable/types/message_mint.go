package types

import (
	"cosmossdk.io/math"
)

func NewMsgMint(authorityContract string, amount math.Int) *MsgMint {
	return &MsgMint{
		AuthorityContract: authorityContract,
		Amount:            amount,
	}
}
