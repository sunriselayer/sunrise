package types

func NewMsgSelfDelegate(sender string) *MsgSelfDelegate {
	return &MsgSelfDelegate{
		Creator: sender,
	}
}
