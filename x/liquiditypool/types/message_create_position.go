package types

func NewMsgCreatePosition(creator string) *MsgCreatePosition {
	return &MsgCreatePosition{
		Creator: creator,
	}
}
