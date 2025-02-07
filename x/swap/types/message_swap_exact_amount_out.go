package types

func NewMsgSwapExactAmountOut(sender string) *MsgSwapExactAmountOut {
	return &MsgSwapExactAmountOut{
		Sender: sender,
	}
}
