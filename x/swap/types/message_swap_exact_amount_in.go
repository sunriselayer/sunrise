package types

func NewMsgSwapExactAmountIn(creator string) *MsgSwapExactAmountIn {
	return &MsgSwapExactAmountIn{
		Creator: creator,
	}
}
