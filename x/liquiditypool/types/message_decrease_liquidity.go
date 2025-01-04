package types

func NewMsgDecreaseLiquidity(sender string) *MsgDecreaseLiquidity {
	return &MsgDecreaseLiquidity{
		Sender: sender,
	}
}
