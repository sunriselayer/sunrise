package types

func NewMsgConvert(sender string) *MsgConvert {
	return &MsgConvert{
		Sender: sender,
	}
}
