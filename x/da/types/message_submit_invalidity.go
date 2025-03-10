package types

func NewMsgSubmitInvalidity(sender string) *MsgSubmitInvalidity {
	return &MsgSubmitInvalidity{
		Sender: sender,
	}
}
