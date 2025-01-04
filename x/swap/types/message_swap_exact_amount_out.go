package types

func NewMsgSwapExactAmountOut(creator string) *MsgSwapExactAmountOut {
	return &MsgSwapExactAmountOut{
		Creator: creator,
	}
}
