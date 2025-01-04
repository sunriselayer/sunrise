package types

func NewMsgVoteGauge(creator string) *MsgVoteGauge {
	return &MsgVoteGauge{
		Sender: creator,
	}
}
