package types

func NewMsgBurn(creator string) *MsgBurn {
  return &MsgBurn{
		Creator: creator,
	}
}