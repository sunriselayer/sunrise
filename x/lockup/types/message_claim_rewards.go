package types

func NewMsgClaimRewards(owner string, validator string) *MsgClaimRewards {
	return &MsgClaimRewards{
		Owner:     owner,
		Validator: validator,
	}
}
