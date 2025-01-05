package types

func NewMsgSelfUndelegate(sender string) *MsgSelfUndelegate {
	return &MsgSelfUndelegate{
		Creator: sender,
	}
}
