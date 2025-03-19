package types

func NewMsgWithdrawUnbonded(sender string) *MsgWithdrawUnbonded {
	return &MsgWithdrawUnbonded{
		Sender: sender,
	}
}
