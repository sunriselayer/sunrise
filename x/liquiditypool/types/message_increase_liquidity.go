package types

func NewMsgIncreaseLiquidity(sender string, id uint64) *MsgIncreaseLiquidity {
	return &MsgIncreaseLiquidity{
		Id:     id,
		Sender: sender,
	}
}
