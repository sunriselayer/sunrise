package types

func NewMsgSelfDelegate(creator string) *MsgSelfDelegate {
  return &MsgSelfDelegate{
		Creator: creator,
	}
}