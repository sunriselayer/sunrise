package types

func NewMsgDecreaseLiquidity(creator string) *MsgDecreaseLiquidity {
	return &MsgDecreaseLiquidity{
		Creator: creator,
	}
}
