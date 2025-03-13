package types

func NewMsgPublishData(sender string) *MsgPublishData {
	return &MsgPublishData{
		Sender: sender,
	}
}
