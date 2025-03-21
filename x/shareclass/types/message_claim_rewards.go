package types

func NewMsgClaimRewards(creator string) *MsgClaimRewards {
	return &MsgClaimRewards{
		Creator: creator,
	}
}
