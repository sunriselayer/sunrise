package types

func NewMsgSubmitProof(creator string) *MsgSubmitProof {
	return &MsgSubmitProof{
		Creator: creator,
	}
}
