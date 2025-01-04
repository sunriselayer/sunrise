package types

func NewMsgChallengeForFraud(creator string) *MsgChallengeForFraud {
	return &MsgChallengeForFraud{
		Creator: creator,
	}
}
