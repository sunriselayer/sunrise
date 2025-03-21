package types

func NewMsgNonVotingDelegate(creator string) *MsgNonVotingDelegate {
	return &MsgNonVotingDelegate{
		Creator: creator,
	}
}
