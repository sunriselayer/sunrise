package types

func NewMsgSelfCancelUnbonding(sender string) *MsgSelfCancelUnbonding {
	return &MsgSelfCancelUnbonding{
		Creator: sender,
	}
}
