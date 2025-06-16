package types

func NewMsgClaimYield(sender string, tokenCreator string) *MsgClaimYield {
	return &MsgClaimYield{
		Sender:       sender,
		TokenCreator: tokenCreator,
	}
}
