package types

func NewMsgChallengeForFraud(sender string) *MsgChallengeForFraud {
	return &MsgChallengeForFraud{
		Sender: sender,
	}
}
