package types

func NewMsgSubmitValidityProof(sender string) *MsgSubmitValidityProof {
	return &MsgSubmitValidityProof{
		Sender: sender,
	}
}
