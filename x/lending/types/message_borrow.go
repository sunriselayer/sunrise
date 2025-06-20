package types

import sdk "github.com/cosmos/cosmos-sdk/types"

func NewMsgBorrow(sender string, amount sdk.Coin, collateralPoolId uint64, collateralPositionId uint64) *MsgBorrow {
	return &MsgBorrow{
		Sender:               sender,
		Amount:               amount,
		CollateralPoolId:     collateralPoolId,
		CollateralPositionId: collateralPositionId,
	}
}
