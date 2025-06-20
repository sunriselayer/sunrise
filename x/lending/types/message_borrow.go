package types

func NewMsgBorrow(sender string, borrowDenom string, collateralPoolId uint64, collateralPositionId uint64) *MsgBorrow {
	return &MsgBorrow{
		Sender:               sender,
		BorrowDenom:          borrowDenom,
		CollateralPoolId:     collateralPoolId,
		CollateralPositionId: collateralPositionId,
	}
}
