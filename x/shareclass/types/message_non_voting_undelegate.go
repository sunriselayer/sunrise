package types

func NewMsgNonVotingUndelegate(creator string) *MsgNonVotingUndelegate {
	return &MsgNonVotingUndelegate{
		Creator: creator,
	}
}
