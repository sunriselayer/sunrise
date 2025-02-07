package types

func NewMsgDecreaseLiquidity(sender string, id uint64) *MsgDecreaseLiquidity {
	return &MsgDecreaseLiquidity{
		Id:     id,
		Sender: sender,
	}
}
