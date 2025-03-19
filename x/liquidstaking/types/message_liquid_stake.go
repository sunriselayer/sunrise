package types

func NewMsgLiquidStake(creator string) *MsgLiquidStake {
	return &MsgLiquidStake{
		Creator: creator,
	}
}
