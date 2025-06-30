package types

import (
	"cosmossdk.io/math"
)

func NewMsgBurn(authorityContract string, amount math.Int) *MsgBurn {
	return &MsgBurn{
		AuthorityContract: authorityContract,
		Amount:            amount,
	}
}
