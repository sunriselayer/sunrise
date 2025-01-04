package types

func NewMsgCollectVoteRewards(creator string) *MsgCollectVoteRewards {
	return &MsgCollectVoteRewards{
		Creator: creator,
	}
}
