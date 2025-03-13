package types

func NewMsgSwapExactAmountIn(sender string) *MsgSwapExactAmountIn {
	return &MsgSwapExactAmountIn{
		Sender: sender,
	}
}
