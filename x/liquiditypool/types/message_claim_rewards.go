package types

func NewMsgClaimRewards(sender string, positionIds []uint64) *MsgClaimRewards {
	return &MsgClaimRewards{
		Sender:      sender,
		PositionIds: positionIds,
	}
}
