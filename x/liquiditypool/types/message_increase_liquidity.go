package types

func NewMsgIncreaseLiquidity(sender string) *MsgIncreaseLiquidity {
	return &MsgIncreaseLiquidity{
		Sender: sender,
	}
}
