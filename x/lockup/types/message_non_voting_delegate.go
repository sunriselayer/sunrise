package types

func NewMsgNonVotingDelegate(owner string) *MsgNonVotingDelegate {
	return &MsgNonVotingDelegate{
		Owner: owner,
	}
}
