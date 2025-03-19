package types

func NewMsgLiquidStake(sender string) *MsgLiquidStake {
	return &MsgLiquidStake{
		Sender: sender,
	}
}
