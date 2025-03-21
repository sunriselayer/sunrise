package types

func NewMsgNonVotingDelegate(sender string) *MsgNonVotingDelegate {
	return &MsgNonVotingDelegate{
		Sender: sender,
	}
}
