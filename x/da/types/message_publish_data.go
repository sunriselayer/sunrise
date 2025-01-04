package types

func NewMsgPublishData(creator string) *MsgPublishData {
	return &MsgPublishData{
		Creator: creator,
	}
}
