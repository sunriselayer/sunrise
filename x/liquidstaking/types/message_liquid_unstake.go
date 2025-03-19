package types

func NewMsgLiquidUnstake(creator string) *MsgLiquidUnstake {
	return &MsgLiquidUnstake{
		Creator: creator,
	}
}
