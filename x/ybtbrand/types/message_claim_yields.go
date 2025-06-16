package types

func NewMsgClaimYields(sender string, tokenCreator string) *MsgClaimYields {
	return &MsgClaimYields{
		Sender:       sender,
		TokenCreator: tokenCreator,
	}
}
