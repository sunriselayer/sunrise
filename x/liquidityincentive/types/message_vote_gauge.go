package types

func NewMsgVoteGauge(creator string) *MsgVoteGauge {
	return &MsgVoteGauge{
		Creator: creator,
	}
}
