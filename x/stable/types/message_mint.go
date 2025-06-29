package types

func NewMsgMint(creator string) *MsgMint {
  return &MsgMint{
		Creator: creator,
	}
}