package types

func NewMsgNonVotingUndelegate(owner string) *MsgNonVotingUndelegate {
	return &MsgNonVotingUndelegate{
		Owner: owner,
	}
}
