package types

func NewMsgSend(owner string) *MsgSend {
	return &MsgSend{
		Owner: owner,
	}
}
