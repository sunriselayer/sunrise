package types

func NewMsgMint(admin string, tokenCreator string, amount int64) *MsgMint {
	return &MsgMint{
		Admin:        admin,
		TokenCreator: tokenCreator,
		Amount:       amount,
	}
}
