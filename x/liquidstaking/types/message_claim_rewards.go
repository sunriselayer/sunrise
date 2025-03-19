package types

func NewMsgClaimRewards(sender string) *MsgClaimRewards {
	return &MsgClaimRewards{
		Sender: sender,
	}
}
