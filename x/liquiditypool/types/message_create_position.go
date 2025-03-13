package types

func NewMsgCreatePosition(sender string) *MsgCreatePosition {
	return &MsgCreatePosition{
		Sender: sender,
	}
}
