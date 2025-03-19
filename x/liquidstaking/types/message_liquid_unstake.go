package types

func NewMsgLiquidUnstake(sender string) *MsgLiquidUnstake {
	return &MsgLiquidUnstake{
		Sender: sender,
	}
}
