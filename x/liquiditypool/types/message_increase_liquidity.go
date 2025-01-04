package types

func NewMsgIncreaseLiquidity(creator string) *MsgIncreaseLiquidity {
	return &MsgIncreaseLiquidity{
		Creator: creator,
	}
}
