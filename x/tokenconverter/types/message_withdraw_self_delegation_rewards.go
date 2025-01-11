package types

func NewMsgWithdrawSelfDelegationRewards(sender string) *MsgWithdrawSelfDelegationRewards {
	return &MsgWithdrawSelfDelegationRewards{
		Sender: sender,
	}
}
