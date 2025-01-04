package types

func NewMsgCollectVoteRewards(creator string) *MsgCollectVoteRewards {
	return &MsgCollectVoteRewards{
		Sender: creator,
	}
}
