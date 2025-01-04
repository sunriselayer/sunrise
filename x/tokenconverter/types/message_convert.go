package types

func NewMsgConvert(creator string) *MsgConvert {
	return &MsgConvert{
		Creator: creator,
	}
}
