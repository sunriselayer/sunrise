package types

func NewMsgSubmitProof(sender string) *MsgSubmitProof {
	return &MsgSubmitProof{
		Sender: sender,
	}
}
