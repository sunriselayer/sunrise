package types

func NewMsgWithdrawSelfUnbonded(sender string) *MsgWithdrawSelfUnbonded {
	return &MsgSelfDelegate{
		Creator: sender,
	}
}
