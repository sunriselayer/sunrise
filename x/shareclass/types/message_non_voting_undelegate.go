package types

func NewMsgNonVotingUndelegate(sender string) *MsgNonVotingUndelegate {
	return &MsgNonVotingUndelegate{
		Sender: sender,
	}
}
